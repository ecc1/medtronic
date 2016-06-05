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

func (pump *Pump) SetTimeout(t time.Duration) time.Duration {
	prev := pump.timeout
	pump.timeout = t
	return prev
}

func (pump *Pump) SetRetries(n int) int {
	prev := pump.retries
	pump.retries = n
	return prev
}

func (pump *Pump) Rssi() int {
	return pump.rssi
}

type ResponseHandler func([]byte) interface{}

func commandPacket(cmd CommandCode, params []byte) radio.Packet {
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
func (pump *Pump) Execute(cmd CommandCode, handler ResponseHandler, params ...byte) (interface{}, error) {
	result, err := pump.perform(cmd, nil, handler)
	if err != nil || len(params) == 0 {
		return result, err
	}
	return pump.perform(cmd, params, handler)
}

func (pump *Pump) perform(cmd CommandCode, params []byte, handler ResponseHandler) (interface{}, error) {
	packet := commandPacket(cmd, params)
	for tries := 0; tries < pump.retries || pump.retries == 0; tries++ {
		pump.Radio.Outgoing() <- packet
		timeout := time.After(pump.timeout)
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
			return nil, BadResponseError{command: cmd, data: data}
		}
		pump.rssi = response.Rssi
		if handler != nil {
			result := handler(data[5:])
			if result == nil {
				return nil, BadResponseError{command: cmd, data: data}
			}
			return result, nil
		}
		return nil, nil
	}
	return nil, NoResponseError(cmd)
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
