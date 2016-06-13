package cc1101

import (
	"bytes"
	"fmt"

	"github.com/ecc1/gpio"
	"github.com/ecc1/medtronic/radio"
	"github.com/ecc1/spi"
)

const (
	spiSpeed  = 6000000 // Hz
	interruptPin   = 14      // Intel Edison GPIO connected to GDO0
	hwVersion = 0x0014
)

type Radio struct {
	device        *spi.Device
	interruptPin  gpio.InputPin
	receiveBuffer bytes.Buffer
	stats         radio.Statistics
	err           error
}

func Open() *Radio {
	r := Radio{}
	r.device, r.err = spi.Open(spiSpeed)
	if r.Error() != nil {
		return &r
	}
	r.err = r.device.SetMaxSpeed(spiSpeed)
	if r.Error() != nil {
		return &r
	}
	r.interruptPin, r.err = gpio.Input(interruptPin, "both", false)
	if r.Error() != nil {
		return &r
	}
	v := r.Version()
	if v != hwVersion {
		r.err = fmt.Errorf("unexpected hardware version (%04X instead of %04X)", v, hwVersion)
	}
	return &r
}

func (r *Radio) Init(frequency uint32) {
	r.Reset()
	r.InitRF(frequency)
}

func (r *Radio) Statistics() radio.Statistics {
	return r.stats
}

func (r *Radio) Error() error {
	return r.err
}

func (r *Radio) SetError(err error) {
	r.err = err
}
