package medtronic

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ecc1/radio"
)

const (
	pumpEnvVar             = "MEDTRONIC_PUMP_ID"
	PumpDevice             = 0xA7
	Ack                    = 0x06
	maxPacketSize          = 70 // excluding CRC byte
	defaultResponseTimeout = 500 * time.Millisecond
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
		log.Fatalf("%s environment variable is not set\n", pumpEnvVar)
	}
	if len(id) != 6 {
		log.Fatalf("%s environment variable must be 6 digits\n", pumpEnvVar)
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

func noResponse(code CommandCode) error {
	return fmt.Errorf("no response to %s", code.String())
}

func unexpectedResponse(code CommandCode, data []byte) error {
	return fmt.Errorf("unexpected response to %s: % X", code.String(), data)
}

type PumpCommand struct {
	Code            CommandCode
	Params          []byte
	ResponseHandler func([]byte) interface{}
	ResponseTimeout time.Duration
	NumRetries      int
	Rssi            *int
}

func commandPacket(cmd PumpCommand) radio.Packet {
	initCommandPrefix()
	n := len(commandPrefix)
	data := make([]byte, n+2)
	copy(data, commandPrefix)
	data[n] = byte(cmd.Code)
	data[n+1] = byte(0)
	return EncodePacket(data)
}

func paramsPacket(cmd PumpCommand) radio.Packet {
	initCommandPrefix()
	n := len(commandPrefix)
	data := make([]byte, maxPacketSize)
	copy(data, commandPrefix)
	data[n] = byte(cmd.Code)
	data[n+1] = byte(len(cmd.Params))
	copy(data[n+2:], cmd.Params)
	return EncodePacket(data)
}

// Commands with parameters require an initial exchange with no parameters,
// followed by an exchange with arguments.
func (pump *Pump) Execute(cmd PumpCommand) (interface{}, error) {
	result, err := pump.perform(cmd, false)
	if err != nil || len(cmd.Params) == 0 {
		return result, err
	}
	if cmd.ResponseHandler != nil {
		panic("ResponseHandler != nil")
	}
	return pump.perform(cmd, true)
}

func (pump *Pump) perform(cmd PumpCommand, params bool) (interface{}, error) {
	var packet radio.Packet
	if params {
		packet = paramsPacket(cmd)
	} else {
		packet = commandPacket(cmd)
	}
	responseTimeout := defaultResponseTimeout
	if cmd.ResponseTimeout != 0 {
		responseTimeout = cmd.ResponseTimeout
	}
	for tries := 0; tries < cmd.NumRetries || cmd.NumRetries == 0; tries++ {
		pump.Radio.Outgoing() <- packet
		timeout := time.After(responseTimeout)
		var response radio.Packet
		select {
		case response = <-pump.Radio.Incoming():
			break
		case <-timeout:
			continue
		}
		data, err := pump.DecodePacket(response)
		if err != nil {
			continue
		}
		if !expected(cmd, data) {
			return nil, unexpectedResponse(cmd.Code, data)
		}
		if cmd.Rssi != nil {
			*cmd.Rssi = response.Rssi
		}
		if cmd.ResponseHandler != nil {
			result := cmd.ResponseHandler(data[5:])
			if result == nil {
				return nil, unexpectedResponse(cmd.Code, data)
			}
			return result, nil
		}
		return nil, nil
	}
	return nil, noResponse(cmd.Code)
}

func expected(cmd PumpCommand, data []byte) bool {
	if len(data) < 5 {
		return false
	}
	if !bytes.Equal(data[:len(commandPrefix)], commandPrefix) {
		return false
	}
	if len(cmd.Params) != 0 || cmd.Code == PowerControl {
		return data[4] == byte(Ack)
	} else {
		return data[4] == byte(cmd.Code)
	}
}
