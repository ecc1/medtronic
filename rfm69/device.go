package rfm69

import (
	"bytes"
	"time"

	"github.com/ecc1/gpio"
	"github.com/ecc1/medtronic/radio"
)

const (
	spiSpeed  = 6000000 // Hz
	resetPin  = 12      // Intel Edison GPIO for hardware reset
	hwVersion = 0x0204
)

type flavor struct{}

func (hw flavor) Name() string {
	return "RFM69HCW"
}

func (hw flavor) Speed() int {
	return spiSpeed
}

func (hw flavor) ReadSingleAddress(addr byte) byte {
	return addr
}

func (hw flavor) ReadBurstAddress(addr byte) byte {
	return addr
}

func (hw flavor) WriteSingleAddress(addr byte) byte {
	return SpiWriteMode | addr
}

func (hw flavor) WriteBurstAddress(addr byte) byte {
	return SpiWriteMode | addr
}

type Radio struct {
	hw            *radio.Hardware
	resetPin      gpio.OutputPin
	receiveBuffer bytes.Buffer
	stats         radio.Statistics
	err           error
}

func Open() radio.Interface {
	r := &Radio{hw: radio.Open(flavor{})}
	v := r.Version()
	if r.Error() != nil {
		return r
	}
	if v != hwVersion {
		r.hw.Close()
		r.SetError(radio.HardwareVersionError{Actual: v, Expected: hwVersion})
		return r
	}
	r.resetPin, r.err = gpio.Output(resetPin, false)
	if r.Error() != nil {
		r.hw.Close()
		return r
	}
	return r
}

func (r *Radio) Version() uint16 {
	v := r.hw.ReadRegister(RegVersion)
	return uint16(v>>4)<<8 | uint16(v&0xF)
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
	r.setMode(SleepMode)
}

func (r *Radio) Statistics() radio.Statistics {
	return r.stats
}

func (r *Radio) Error() error {
	err := r.hw.Error()
	if err != nil {
		return err
	}
	return r.err
}

func (r *Radio) SetError(err error) {
	r.hw.SetError(err)
	r.err = err
}

func (r *Radio) Hardware() *radio.Hardware {
	return r.hw
}
