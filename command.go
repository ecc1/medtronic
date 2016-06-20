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

type CommandCode byte

//go:generate stringer -type CommandCode

const Ack CommandCode = 0x06

type NoResponseError CommandCode

func (e NoResponseError) Error() string {
	return fmt.Sprintf("no response to %s", CommandCode(e).String())
}

type BadResponseError struct {
	command CommandCode
	data    []byte
}

func (e BadResponseError) Error() string {
	return fmt.Sprintf("unexpected response to %s: % X", e.command.String(), e.data)
}

func (pump *Pump) BadResponse(cmd CommandCode, data []byte) {
	pump.SetError(BadResponseError{command: cmd, data: data})
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
func commandPacket(cmd CommandCode, params []byte) []byte {
	initCommandPrefix()
	n := len(commandPrefix)
	data := []byte{}
	if len(params) == 0 {
		data = make([]byte, n+3)
	} else {
		data = make([]byte, maxPacketSize)
	}
	copy(data, commandPrefix)
	data[n] = byte(cmd)
	data[n+1] = byte(len(params))
	if len(params) != 0 {
		copy(data[n+2:], params)
	}
	return EncodePacket(data)
}

// Commands with parameters require an initial exchange with no parameters,
// followed by an exchange with the actual arguments.
func (pump *Pump) Execute(cmd CommandCode, params ...byte) []byte {
	if len(params) != 0 {
		pump.perform(cmd, Ack, nil)
		return pump.perform(cmd, Ack, params)
	}
	return pump.perform(cmd, cmd, nil)
}

// History pages are returned as a series of 65-byte records:
//   sequence number (1 to 16)
//   64 bytes of payload
// The caller must Ack the record in order to receive the next one.
// The 0x80 bit is set in the sequence number of the final record.
// The page consists of the concatenated payloads.
// The final 2 bytes are the CRC-16 of the preceding data.
func (pump *Pump) Download(cmd CommandCode, page int) []byte {
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
			pump.SetError(fmt.Errorf("unexpected history record length (%d)", len(data)))
			break
		}
		done := data[0]&0x80 != 0
		seqNum := data[0] &^ 0x80
		payload := data[1:]
		// Skip duplicate responses.
		if seqNum != prev {
			results = append(results, payload...)
			prev = seqNum
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

func (pump *Pump) perform(cmd CommandCode, resp CommandCode, params []byte) []byte {
	if pump.Error() != nil {
		return nil
	}
	packet := commandPacket(cmd, params)
	for tries := 0; tries < pump.retries || pump.retries == 0; tries++ {
		pump.Radio.Send(packet)
		response, rssi := pump.Radio.Receive(pump.timeout)
		if len(response) == 0 {
			pump.SetError(nil)
			continue
		}
		data := pump.DecodePacket(response)
		if pump.Error() != nil {
			pump.SetError(nil)
			continue
		}
		if !expected(cmd, resp, data) {
			pump.BadResponse(cmd, data)
			return nil
		}
		pump.rssi = rssi
		return data[5:]
	}
	pump.SetError(NoResponseError(cmd))
	return nil
}

func expected(cmd CommandCode, resp CommandCode, data []byte) bool {
	if len(data) < 5 {
		return false
	}
	if !bytes.Equal(data[:len(commandPrefix)], commandPrefix) {
		return false
	}
	return data[4] == byte(cmd) || data[4] == byte(resp) || (cmd == Wakeup && data[4] == byte(Ack))
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
