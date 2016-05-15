package cc1100

import (
	"github.com/ecc1/gpio"
	"github.com/ecc1/spi"
)

const (
	spiSpeed = 6000000 // Hz
	gpioPin  = 14      // Intel Edison GPIO connected to GDO0
)

type Device struct {
	spiDev       *spi.Device
	interruptPin gpio.InputPin

	receiverStarted bool
	receivedPackets chan []byte

	packetsSent     int
	packetsReceived int
	decodingErrors  int
	crcErrors       int
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
	return &Device{spiDev: spiDev, interruptPin: pin}, nil
}
