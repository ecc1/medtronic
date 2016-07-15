package medtronic

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

const (
	pumpEnvVar      = "MEDTRONIC_PUMP_ID"
	PumpDevice      = 0xA7
	maxPacketSize   = 71 // including CRC byte
	historyPageSize = 1024
)

var (
	commandPrefix []byte
)

func initCommandPrefix() {
	if len(commandPrefix) != 0 {
		return
	}
	id := os.Getenv(pumpEnvVar)
	if len(id) == 0 {
		log.Fatalf("%s environment variable is not set", pumpEnvVar)
	}
	if len(id) != 6 {
		log.Fatalf("%s environment variable must be 6 digits", pumpEnvVar)
	}
	commandPrefix = []byte{
		PumpDevice,
		(id[0]-'0')<<4 | (id[1] - '0'),
		(id[2]-'0')<<4 | (id[3] - '0'),
		(id[4]-'0')<<4 | (id[5] - '0'),
	}
}

type Command byte

//go:generate stringer -type Command

const (
	Ack          Command = 0x06
	CommandError Command = 0x15
)

type NoResponseError Command

func (e NoResponseError) Error() string {
	return fmt.Sprintf("no response to %v", Command(e))
}

type InvalidCommandError Command

func (e InvalidCommandError) Error() string {
	return fmt.Sprintf("invalid %v command", Command(e))
}

type BadResponseError struct {
	Command Command
	Data    []byte
}

func (e BadResponseError) Error() string {
	return fmt.Sprintf("unexpected response to %v: % X", e.Command, e.Data)
}

func (pump *Pump) BadResponse(cmd Command, data []byte) {
	pump.SetError(BadResponseError{Command: cmd, Data: data})
}

// commandPacket constructs a packet
// with the specified command code and parameters,
// A command packet with no parameters is 7 bytes long:
//   device type (0xA7)
//   3 bytes of pump ID
//   command code
//   length of parameters (0)
//   CRC-8
// A command packet with parameters is 71 bytes long:
//   device type (0xA7)
//   3 bytes of pump ID
//   command code
//   length of parameters
//   64 bytes of parameters plus padding
//   CRC-8
func commandPacket(cmd Command, params []byte) []byte {
	initCommandPrefix()
	data := []byte{}
	if len(params) == 0 {
		data = make([]byte, 7)
	} else {
		data = make([]byte, maxPacketSize)
	}
	copy(data, commandPrefix)
	data[4] = byte(cmd)
	data[5] = byte(len(params))
	if len(params) != 0 {
		copy(data[6:], params)
	}
	return EncodePacket(data)
}

// Commands with parameters require an initial exchange with no parameters,
// followed by an exchange with the actual arguments.
func (pump *Pump) Execute(cmd Command, params ...byte) []byte {
	if len(params) != 0 {
		pump.perform(cmd, Ack, nil)
		return pump.perform(cmd, Ack, params)
	}
	return pump.perform(cmd, cmd, nil)
}

// History pages are returned as a series of 65-byte fragments:
//   sequence number (1 to 16)
//   64 bytes of payload
// The caller must Ack the fragment in order to receive the next one.
// The 0x80 bit is set in the sequence number of the final fragment.
// The page consists of the concatenated payloads.
// The final 2 bytes are the CRC-16 of the preceding data.
func (pump *Pump) Download(cmd Command, page int) []byte {
	data := pump.Execute(cmd, byte(page))
	if pump.Error() != nil {
		return nil
	}
	results := []byte{}
	retries := pump.Retries()
	pump.SetRetries(1)
	defer pump.SetRetries(retries)
	prev := byte(0)
	for {
		if len(data) != 65 {
			pump.SetError(fmt.Errorf("unexpected history fragment length (%d)", len(data)))
			break
		}
		done := data[0]&0x80 != 0
		seqNum := data[0] &^ 0x80
		payload := data[1:]
		if seqNum == prev {
			// Skip duplicate responses.
		} else if seqNum == prev+1 {
			results = append(results, payload...)
			prev = seqNum
		} else {
			pump.SetError(fmt.Errorf("received fragment %d instead of %d in history page", seqNum, prev+1))
			break
		}
		if done {
			if seqNum != 16 {
				pump.SetError(fmt.Errorf("unexpected final sequence number for history page (%d)", seqNum))
			}
			break
		}
		next := pump.perform(Ack, cmd, nil)
		if pump.Error() != nil {
			_, noResponse := pump.Error().(NoResponseError)
			if noResponse {
				pump.SetError(nil)
				continue
			}
			break
		}
		data = next
	}
	if pump.Error() != nil {
		return nil
	}
	if len(results) != historyPageSize {
		pump.SetError(fmt.Errorf("unexpected history page size (%d)", len(results)))
		return nil
	}
	dataCrc := twoByteUint(results[historyPageSize-2:])
	results = results[:historyPageSize-2]
	calcCrc := Crc16(results)
	if dataCrc != calcCrc {
		pump.SetError(fmt.Errorf("CRC should be %02X, not %02X", calcCrc, dataCrc))
		return nil
	}
	return results
}

func (pump *Pump) perform(cmd Command, resp Command, params []byte) []byte {
	if pump.Error() != nil {
		return nil
	}
	packet := commandPacket(cmd, params)
	for tries := 0; tries < pump.retries || pump.retries == 0; tries++ {
		pump.Radio.Send(packet)
		response, rssi := pump.Radio.Receive(pump.Timeout())
		if len(response) == 0 {
			pump.SetError(nil)
			continue
		}
		data := pump.DecodePacket(response)
		if pump.Error() != nil {
			pump.SetError(nil)
			continue
		}
		if pump.unexpected(cmd, resp, data) {
			return nil
		}
		pump.rssi = rssi
		return data[5:]
	}
	pump.SetError(NoResponseError(cmd))
	return nil
}

func (pump *Pump) unexpected(cmd Command, resp Command, data []byte) bool {
	if len(data) < 5 {
		pump.BadResponse(cmd, data)
		return true
	}
	n := len(commandPrefix)
	if !bytes.Equal(data[:n], commandPrefix) {
		pump.BadResponse(cmd, data)
		return true
	}
	switch Command(data[n]) {
	case cmd:
		return false
	case resp:
		return false
	case Ack:
		if cmd != Wakeup {
			break
		}
		return false
	case CommandError:
		pump.SetError(InvalidCommandError(cmd))
		return true
	}
	pump.BadResponse(cmd, data)
	return true
}

func twoByteInt(data []byte) int {
	return int(data[0])<<8 | int(data[1])
}

func twoByteUint(data []byte) uint16 {
	return uint16(data[0])<<8 | uint16(data[1])
}

func fourByteInt(data []byte) int {
	return twoByteInt(data[0:2])<<16 | twoByteInt(data[2:4])
}
