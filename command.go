package medtronic

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ecc1/cc1100"
)

const (
	pumpEnvVar             = "MEDTRONIC_PUMP_ID"
	PumpDevice             = 0xA7
	Ack                    = 0x06
	defaultResponseTimeout = 100 * time.Millisecond
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

func commandPacket(cmd PumpCommand) cc1100.Packet {
	initCommandPrefix()
	data := append(commandPrefix, byte(cmd.Code), byte(len(cmd.Params)))
	if len(cmd.Params) != 0 {
		data = append(data, cmd.Params...)
	}
	return EncodePacket(data)
}

func (pump *Pump) Execute(cmd PumpCommand) (interface{}, error) {
	packet := commandPacket(cmd)
	responseTimeout := defaultResponseTimeout
	if cmd.ResponseTimeout != 0 {
		responseTimeout = cmd.ResponseTimeout
	}
	for tries := 0; tries < cmd.NumRetries || cmd.NumRetries == 0; tries++ {
		pump.Radio.Outgoing() <- packet
		timeout := time.After(responseTimeout)
		var response cc1100.Packet
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
		if !expected(cmd.Code, data) {
			return nil, unexpectedResponse(cmd.Code, data)
		}
		if cmd.Rssi != nil {
			*cmd.Rssi = response.Rssi
		}
		result := cmd.ResponseHandler(data[5:])
		if result == nil {
			return nil, unexpectedResponse(cmd.Code, data)
		}
		return result, nil

	}
	return nil, noResponse(cmd.Code)
}

func expected(code CommandCode, data []byte) bool {
	if len(data) < 5 {
		return false
	}
	if !bytes.Equal(data[:len(commandPrefix)], commandPrefix) {
		return false
	}
	if code == PowerControl {
		return data[4] == byte(Ack)
	} else {
		return data[4] == byte(code)
	}
}
