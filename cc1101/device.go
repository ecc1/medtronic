package cc1101

import (
	"bytes"
	"log"

	"github.com/ecc1/medtronic/radio"
)

const (
	spiSpeed  = 6000000 // Hz
	hwVersion = 0x0014
)

type flavor struct{}

func (hw flavor) Name() string {
	return "CC1101"
}

func (hw flavor) Speed() int {
	return spiSpeed
}

func (hw flavor) ReadSingleAddress(addr byte) byte {
	return READ_MODE | addr
}

func (hw flavor) ReadBurstAddress(addr byte) byte {
	reg := addr & 0x3F
	if 0x30 <= reg && reg <= 0x3D {
		log.Panicf("no burst access for CC1101 status register %02X", reg)
	}
	return READ_MODE | BURST_MODE | addr
}

func (hw flavor) WriteSingleAddress(addr byte) byte {
	return addr
}

func (hw flavor) WriteBurstAddress(addr byte) byte {
	return BURST_MODE | addr
}

type Radio struct {
	hw            *radio.Hardware
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
	return r
}

func (r *Radio) Close() {
	r.changeState(SIDLE, STATE_IDLE)
	r.hw.Close()
}

func (r *Radio) Version() uint16 {
	p := r.hw.ReadRegister(PARTNUM)
	v := r.hw.ReadRegister(VERSION)
	return uint16(p)<<8 | uint16(v)
}

func (r *Radio) Strobe(cmd byte) byte {
	if verbose && cmd != SNOP {
		log.Printf("issuing %s command", strobeName(cmd))
	}
	buf := []byte{cmd}
	r.err = r.hw.SpiDevice().Transfer(buf)
	return buf[0]
}

func (r *Radio) Reset() {
	r.Strobe(SRES)
}

func (r *Radio) Init(frequency uint32) {
	r.Reset()
	r.InitRF(frequency)
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
