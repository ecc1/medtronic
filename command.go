package medtronic

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ecc1/medtronic/packet"
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
	Ack Command = 0x06
	Nak Command = 0x15
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
// with the specified command code and parameters.
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
	return packet.Encode(data)
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
// The caller must send an Ack to receive the next fragment
// or a Nak to have the current one retransmitted.
// The 0x80 bit is set in the sequence number of the final fragment.
// The page consists of the concatenated payloads.
// The final 2 bytes are the CRC-16 of the preceding data.

const (
	numFragments    = 16
	fragmentLength  = 65
	maxNaks         = 10
	downloadTimeout = 150 * time.Millisecond
)

func (pump *Pump) Download(cmd Command, page int) []byte {
	data := pump.Execute(cmd, byte(page))
	if pump.Error() != nil {
		return nil
	}
	results := []byte{}
	retries := pump.Retries()
	pump.SetRetries(1)
	defer pump.SetRetries(retries)
	timeout := pump.Timeout()
	pump.SetTimeout(downloadTimeout)
	defer pump.SetTimeout(timeout)
	expected := byte(1)
	for {
		if len(data) != fragmentLength {
			pump.SetError(fmt.Errorf("history page %d: unexpected fragment length (%d)", page, len(data)))
			return nil
		}
		done := data[0]&0x80 != 0
		seqNum := data[0] &^ 0x80
		payload := data[1:]
		if seqNum == expected {
			// Got the next fragment as expected.
			results = append(results, payload...)
			if done {
				if seqNum != numFragments {
					pump.SetError(fmt.Errorf("history page %d: unexpected final sequence number (%d)", page, seqNum))
					return nil
				}
				break
			}
			expected = seqNum + 1
		} else if seqNum < expected {
			// Skip duplicate responses.
		} else {
			// Missed fragment.
			pump.SetError(fmt.Errorf("history page %d: received fragment %d instead of %d", page, seqNum, expected))
			return nil
		}
		next := []byte{}
		// Acknowledge the current fragment.
		next = pump.perform(Ack, cmd, nil)
		if pump.Error() == nil {
			data = next
			continue
		}
		_, noResponse := pump.Error().(NoResponseError)
		if !noResponse {
			return nil
		}
		// No response to ACK. Send NAK to request retransmission.
		pump.SetError(nil)
		for count := 0; count < maxNaks; count++ {
			next = pump.perform(Nak, cmd, nil)
			if pump.Error() == nil {
				format := "history page %d: received fragment %d after %d NAK"
				if count != 0 {
					format += "s"
				}
				log.Printf(format, page, next[0]&^0x80, count+1)
				break
			}
			_, noResponse := pump.Error().(NoResponseError)
			if !noResponse {
				return nil
			}
			pump.SetError(nil)
		}
		if next == nil {
			pump.SetError(fmt.Errorf("history page %d: lost fragment %d", page, expected))
			return nil
		}
		data = next
	}
	if len(results) != historyPageSize {
		pump.SetError(fmt.Errorf("history page %d: unexpected size (%d)", page, len(results)))
		return nil
	}
	dataCrc := twoByteUint(results[historyPageSize-2:])
	results = results[:historyPageSize-2]
	calcCrc := packet.Crc16(results)
	if dataCrc != calcCrc {
		pump.SetError(fmt.Errorf("history page %d: CRC should be %02X, not %02X", page, calcCrc, dataCrc))
		return nil
	}
	return results
}

func (pump *Pump) perform(cmd Command, resp Command, params []byte) []byte {
	if pump.Error() != nil {
		return nil
	}
	p := commandPacket(cmd, params)
	for tries := 0; tries < pump.retries || pump.retries == 0; tries++ {
		pump.Radio.Send(p)
		response, rssi := pump.Radio.Receive(pump.Timeout())
		if len(response) == 0 {
			pump.SetError(nil)
			continue
		}
		data, err := packet.Decode(response)
		if err != nil {
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
	case Nak:
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
