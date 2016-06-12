package rfm69

import (
	"bytes"
	"fmt"
	"time"

	"github.com/ecc1/gpio"
	"github.com/ecc1/medtronic/radio"
	"github.com/ecc1/spi"
)

const (
	spiSpeed     = 10000000 // Hz
	interruptPin = 14       // Intel Edison GPIO connected to DIO0
	resetPin     = 12       // Intel Edison GPIO connected to RESET
	hwVersion    = 0x0204
)

type Radio struct {
	device        *spi.Device
	interruptPin  gpio.InputPin
	resetPin      gpio.OutputPin
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
	r.resetPin, r.err = gpio.Output(resetPin, false)
	if r.Error() != nil {
		return &r
	}
	v := r.Version()
	if v != hwVersion {
		r.err = fmt.Errorf("unexpected hardware version (%04X instead of %04X)", v, hwVersion)
	}
	return &r
}

// Reset module.  See section 7.2.2 of data sheet.
func (r *Radio) Reset() {
	if r.Error() != nil {
		return
	}
	r.err = r.resetPin.Write(true)
	if r.Error() != nil {
		r.resetPin.Write(false)
		return
	}
	time.Sleep(100 * time.Microsecond)
	r.err = r.resetPin.Write(false)
	if r.Error() != nil {
		return
	}
	time.Sleep(5 * time.Millisecond)
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
