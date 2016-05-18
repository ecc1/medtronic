package cc1100

import (
	"bytes"

	"github.com/ecc1/gpio"
	"github.com/ecc1/spi"
)

const (
	spiSpeed = 6000000 // Hz
	gpioPin  = 14      // Intel Edison GPIO connected to GDO0
)

type Packet struct {
	Rssi int
	Data []byte
}

type Device struct {
	spiDev             *spi.Device
	interruptPin       gpio.InputPin
	radioStarted       bool
	receiveBuffer      bytes.Buffer
	transmittedPackets chan Packet
	receivedPackets    chan Packet
	interrupt          chan struct{}
	packetsSent        int
	packetsReceived    int
}

func Open() (*Device, error) {
	spiDev, err := spi.Open(spiSpeed)
	if err != nil {
		return nil, err
	}
	err = spiDev.SetMaxSpeed(spiSpeed)
	if err != nil {
		return nil, err
	}
	pin, err := gpio.Input(gpioPin, "both", false)
	if err != nil {
		return nil, err
	}
	return &Device{
		spiDev:             spiDev,
		interruptPin:       pin,
		transmittedPackets: make(chan Packet, 100),
		receivedPackets:    make(chan Packet, 10),
		interrupt:          make(chan struct{}, 10),
	}, nil
}

func (dev *Device) Init() error {
	err := dev.Reset()
	if err != nil {
		return err
	}
	err = dev.InitRF()
	if err != nil {
		return err
	}
	dev.StartRadio()
	return nil
}
