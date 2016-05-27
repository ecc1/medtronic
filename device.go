package cc1100

import (
	"bytes"
	"fmt"

	"github.com/ecc1/gpio"
	"github.com/ecc1/radio"
	"github.com/ecc1/spi"
)

const (
	spiSpeed  = 6000000 // Hz
	gpioPin   = 14      // Intel Edison GPIO connected to GDO0
	hwVersion = 0x0014
)

type Radio struct {
	device       *spi.Device
	interruptPin gpio.InputPin

	radioStarted       bool
	receiveBuffer      bytes.Buffer
	transmittedPackets chan radio.Packet
	receivedPackets    chan radio.Packet
	interrupt          chan struct{}
	stats              radio.Statistics
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
	r := &Radio{
		device:             dev,
		interruptPin:       pin,
		transmittedPackets: make(chan radio.Packet, 100),
		receivedPackets:    make(chan radio.Packet, 10),
		interrupt:          make(chan struct{}, 10),
	}
	v, err := r.Version()
	if err == nil && v != hwVersion {
		err = fmt.Errorf("unexpected hardware version (%04X instead of %04X)", v, hwVersion)
	}
	return r, err
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

func (r *Radio) Statistics() radio.Statistics {
	return r.stats
}
