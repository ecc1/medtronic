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
	spiDev *spi.Device
	rxGPIO gpio.InputPin
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
	g, err := gpio.Input(gpioPin, "both", false)
	if err != nil {
		return nil, err
	}
	return &Device{spiDev: spiDev, rxGPIO: g}, nil
}
