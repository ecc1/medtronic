package cc1100

import (
	"bytes"
	"fmt"
	"time"
)

const (
	PumpDevice = 0xA7
	Ack        = 0x06
)

type PumpCommand byte

//go:generate stringer -type=PumpCommand

const (
	PowerControl PumpCommand = 0x5D
	PumpID       PumpCommand = 0x71
	PumpModel    PumpCommand = 0x8D
)

// TODO: use stringer

func (cmd PumpCommand) String() string {
	switch cmd {
	case PowerControl:
		return "PowerControl"
	case PumpID:
		return "PumpID"
	case PumpModel:
		return "PumpModel"
	default:
		return fmt.Sprintf("PumpCommand %X", cmd)
	}
}

func noResponse(cmd PumpCommand) error {
	return fmt.Errorf("no response to %s", cmd.String())
}

func unexpectedResponse(cmd PumpCommand, data []byte) error {
	return fmt.Errorf("unexpected response to %s: % X", cmd.String(), data)
}

var commandPrefix = []byte{
	PumpDevice,
	pumpID[0]<<4 | pumpID[1],
	pumpID[2]<<4 | pumpID[3],
	pumpID[4]<<4 | pumpID[5],
}

func commandPacket(cmd PumpCommand) Packet {
	return EncodePacket(append(commandPrefix, byte(cmd), 0))
}

type ResponseHandler func([]byte) (interface{}, error)

func (dev *Device) Execute(cmd PumpCommand, responseTimeout time.Duration, retries int, handler ResponseHandler, rssi *int) (interface{}, error) {
	packet := commandPacket(cmd)
	for tries := 0; tries < retries || retries == 0; tries++ {
		dev.OutgoingPackets() <- packet
		timeout := time.After(responseTimeout)
		var response Packet
		select {
		case response = <-dev.IncomingPackets():
			break
		case <-timeout:
			continue
		}
		data, err := dev.DecodePacket(response)
		if err != nil {
			continue
		}
		if !expected(cmd, data) {
			return nil, unexpectedResponse(cmd, data)
		}
		if rssi != nil {
			*rssi = response.Rssi
		}
		return handler(data[5:])
	}
	return nil, noResponse(cmd)
}

func expected(cmd PumpCommand, data []byte) bool {
	if len(data) < 5 {
		return false
	}
	if !bytes.Equal(data[:len(commandPrefix)], commandPrefix) {
		return false
	}
	if cmd == PowerControl {
		return data[4] == byte(Ack)
	} else {
		return data[4] == byte(cmd)
	}
}

func (dev *Device) PumpID(retries int) (string, error) {
	getResult := func(data []byte) (interface{}, error) {
		if len(data) >= 1 {
			n := int(data[0])
			if len(data) >= 1+n {
				return string(data[1 : 1+n]), nil
			}
		}
		return nil, unexpectedResponse(PumpID, data)
	}
	result, err := dev.Execute(PumpID, 200*time.Millisecond, retries, getResult, nil)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

func (dev *Device) PumpModel(retries int, rssi *int) (string, error) {
	getResult := func(data []byte) (interface{}, error) {
		if len(data) >= 2 {
			n := int(data[1])
			if len(data) >= 2+n {
				return string(data[2 : 2+n]), nil
			}
		}
		return nil, unexpectedResponse(PumpModel, data)
	}
	result, err := dev.Execute(PumpModel, 200*time.Millisecond, retries, getResult, rssi)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

func (dev *Device) PowerControl(retries int) error {
	nop := func(_ []byte) (interface{}, error) {
		return nil, nil
	}
	_, err := dev.Execute(PowerControl, 10*time.Second, retries, nop, nil)
	return err
}

func (dev *Device) Wakeup() error {
	const (
		numWakeups = 250
		xmitDelay  = 50 * time.Millisecond
	)
	packet := commandPacket(PowerControl)
	for i := 0; i < numWakeups; i++ {
		dev.OutgoingPackets() <- packet
		time.Sleep(xmitDelay)
	}
	return dev.PowerControl(1)
}
