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

type Radio struct {
	device       *spi.Device
	interruptPin gpio.InputPin

	radioStarted       bool
	receiveBuffer      bytes.Buffer
	transmittedPackets chan Packet
	receivedPackets    chan Packet
	interrupt          chan struct{}

	PacketsSent     int
	PacketsReceived int
}

func Open() (*Radio, error) {
	dev, err := spi.Open(spiSpeed)
	if err != nil {
		return nil, err
	}
	err = dev.SetMaxSpeed(spiSpeed)
	if err != nil {
		return nil, err
	}
	pin, err := gpio.Input(gpioPin, "both", false)
	if err != nil {
		return nil, err
	}
	return &Radio{
		device:             dev,
		interruptPin:       pin,
		transmittedPackets: make(chan Packet, 100),
		receivedPackets:    make(chan Packet, 10),
		interrupt:          make(chan struct{}, 10),
	}, nil
}

func (r *Radio) Init() error {
	err := r.Reset()
	if err != nil {
		return err
	}
	err = r.InitRF()
	if err != nil {
		return err
	}
	r.startRadio()
	return nil
}
