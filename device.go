package rfm69

import (
	"bytes"
	"fmt"
	"time"

	"github.com/ecc1/gpio"
	"github.com/ecc1/radio"
	"github.com/ecc1/spi"
)

const (
	spiSpeed     = 10000000 // Hz
	interruptPin = 14       // Intel Edison GPIO connected to DIO0
	resetPin     = 12       // Intel Edison GPIO connected to RESET
	hwVersion    = 0x0204
)

type Radio struct {
	device       *spi.Device
	interruptPin gpio.InputPin
	resetPin     gpio.OutputPin

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
	intr, err := gpio.Input(interruptPin, "both", false)
	if err != nil {
		return nil, err
	}
	reset, err := gpio.Output(resetPin, false)
	if err != nil {
		return nil, err
	}
	r := &Radio{
		device:             dev,
		interruptPin:       intr,
		resetPin:           reset,
		transmittedPackets: make(chan radio.Packet, 100),
		receivedPackets:    make(chan radio.Packet, 100),
		interrupt:          make(chan struct{}),
	}
	v, err := r.Version()
	if err == nil && v != hwVersion {
		err = fmt.Errorf("unexpected hardware version (%04X instead of %04X)", v, hwVersion)
	}
	return r, err
}

// Reset module.  See section 7.2.2 of data sheet.
func (r *Radio) Reset() error {
	err := r.resetPin.Write(true)
	if err != nil {
		r.resetPin.Write(false)
		return err
	}
	time.Sleep(100 * time.Microsecond)
	err = r.resetPin.Write(false)
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Millisecond)
	return nil
}

func (r *Radio) Init(frequency uint32) error {
	err := r.Reset()
	if err != nil {
		return err
	}
	err = r.InitRF(frequency)
	if err != nil {
		return err
	}
	r.Start()
	return nil
}

func (r *Radio) Statistics() radio.Statistics {
	return r.stats
}
