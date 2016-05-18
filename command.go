package cc1100

import (
	"fmt"
	"time"
)

const (
	PumpDevice   = 0xA7
	Wakeup       = 0x5D
	GetPumpModel = 0x8D
	Ack          = 0x06
)

var (
	NoResponse = fmt.Errorf("no response from pump %v", pumpID)
)

func commandPacket(cmd byte) Packet {
	command := []byte{
		PumpDevice,
		pumpID[0]<<4 | pumpID[1],
		pumpID[2]<<4 | pumpID[3],
		pumpID[4]<<4 | pumpID[5],
		cmd,
		0,
	}
	return EncodePacket(command)
}

func (dev *Device) Wakeup() error {
	const (
		numWakeups      = 200
		xmitDelay       = 100 * time.Millisecond
		responseTimeout = 10 * time.Second
	)
	packet := commandPacket(Wakeup)
	for i := 0; i < numWakeups; i++ {
		dev.OutgoingPackets() <- packet
		time.Sleep(xmitDelay)
	}
	timeout := time.After(responseTimeout)
	var response Packet
	select {
	case response = <-dev.IncomingPackets():
		break
	case <-timeout:
		return NoResponse
	}
	data, err := dev.DecodePacket(response)
	if err != nil {
		return err
	}
	if len(data) == 7 && data[4] == Ack {
		return nil
	}
	return fmt.Errorf("unexpected wakeup response: % X", data)
}
