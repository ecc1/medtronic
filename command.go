package medtronic

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

const (
	pumpEnvVar    = "MEDTRONIC_PUMP_ID"
	PumpDevice    = 0xA7
	Ack           = 0x06
	maxPacketSize = 70 // excluding CRC byte
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

//go:generate stringer -type=CommandCode

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

type ResponseHandler func([]byte) interface{}

func commandPacket(cmd CommandCode, params []byte) []byte {
	initCommandPrefix()
	n := len(commandPrefix)
	var data []byte
	if len(params) == 0 {
		data = make([]byte, n+2)
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
// followed by an exchange with arguments.
func (pump *Pump) Execute(cmd CommandCode, handler ResponseHandler, params ...byte) interface{} {
	result := pump.perform(cmd, nil, handler)
	if len(params) != 0 {
		result = pump.perform(cmd, params, handler)
	}
	return result
}

func (pump *Pump) perform(cmd CommandCode, params []byte, handler ResponseHandler) interface{} {
	if pump.Error() != nil {
		return nil
	}
	packet := commandPacket(cmd, params)
	for tries := 0; tries < pump.retries || pump.retries == 0; tries++ {
		pump.Radio.Send(packet)
		response, rssi := pump.Radio.Receive(pump.timeout)
		if response == nil {
			pump.SetError(nil)
			continue
		}
		data := pump.DecodePacket(response)
		if pump.Error() != nil {
			pump.SetError(nil)
			continue
		}
		if !expected(cmd, data) {
			pump.SetError(BadResponseError{command: cmd, data: data})
			return nil
		}
		pump.rssi = rssi
		if handler != nil {
			result := handler(data[5:])
			if result == nil {
				pump.SetError(BadResponseError{command: cmd, data: data})
			}
			return result
		}
		return nil
	}
	pump.err = NoResponseError(cmd)
	return nil
}

func expected(cmd CommandCode, data []byte) bool {
	if len(data) < 5 {
		return false
	}
	if !bytes.Equal(data[:len(commandPrefix)], commandPrefix) {
		return false
	}
	return data[4] == byte(cmd) || data[4] == byte(Ack)
}

func twoByteInt(data []byte) int {
	return int(data[0])<<8 | int(data[1])
}

func fourByteInt(data []byte) int {
	return twoByteInt(data[0:2])<<16 | twoByteInt(data[2:4])
}
