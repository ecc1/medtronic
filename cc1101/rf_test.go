package cc1101

import (
	"bytes"
	"testing"
)

func TestFrequency(t *testing.T) {
	cases := []struct {
		f       uint32
		b       []byte
		fApprox uint32 // 0 => equal to f
	}{
		{315000000, []byte{0x0D, 0x20, 0x00}, 0},
		{915000000, []byte{0x26, 0x20, 0x00}, 0},
		// some that can't be represented exactly:
		{434000000, []byte{0x12, 0x15, 0x55}, 433999877},
		{868000000, []byte{0x24, 0x2A, 0xAB}, 868000122},
		{916300000, []byte{0x26, 0x2D, 0xDE}, 916300048},
		{916600000, []byte{0x26, 0x31, 0x11}, 916599975},
	}
	for _, c := range cases {
		b := frequencyToRegisters(c.f)
		if !bytes.Equal(b, c.b) {
			t.Errorf("frequencyToRegisters(%d) == % X, want % X", c.f, b, c.b)
		}
		f := registersToFrequency(c.b)
		if c.fApprox == 0 {
			if f != c.f {
				t.Errorf("registersToFrequency(% X) == %d, want %d", c.b, f, c.f)
			}
		} else {
			if f != c.fApprox {
				t.Errorf("registersToFrequency(% X) == %d, want %d", c.b, f, c.fApprox)
			}
		}
	}
}

/*
func TestBitrate(t *testing.T) {
	cases := []struct {
		br       uint32
		b        []byte
		brApprox uint32 // 0 => equal to br
	}{
		{1200, []byte{0x68, 0x2B}, 0},
		{2400, []byte{0x34, 0x15}, 0},
		{25000, []byte{0x05, 0x00}, 0},
		{50000, []byte{0x02, 0x80}, 0},
		// some that can't be represented exactly:
		{16384, []byte{0x07, 0xA1}, 16385},
		{19200, []byte{0x06, 0x83}, 19196},
		{38400, []byte{0x03, 0x41}, 38415},
		{150000, []byte{0x00, 0xD5}, 150235},
	}
	for _, c := range cases {
		b := bitrateToRegisters(c.br)
		if !bytes.Equal(b, c.b) {
			t.Errorf("bitrateToRegisters(%d) == % X, want % X", c.br, b, c.b)
		}
		f := registersToBitrate(c.b)
		if c.brApprox != 0 {
			if f != c.brApprox {
				t.Errorf("registersToBitrate(% X) == %d, want %d", c.b, f, c.brApprox)
			}
		} else {
			if f != c.br {
				t.Errorf("registersToBitrate(% X) == %d, want %d", c.b, f, c.br)
			}
		}
	}
}

func TestChannelBw(t *testing.T) {
	cases := []struct {
		bw       uint32
		r        byte
		bwApprox uint32 // 0 => equal to bw
	}{
		{12500, RxBwMant20 | 4<<RxBwExpShift, 0},
		{25000, RxBwMant20 | 3<<RxBwExpShift, 0},
		{83333, RxBwMant24 | 1<<RxBwExpShift, 0},
		{166666, RxBwMant24 | 0<<RxBwExpShift, 0},
		{200000, RxBwMant20 | 0<<RxBwExpShift, 0},
		{250000, RxBwMant16 | 0<<RxBwExpShift, 0},
		// some that can't be represented exactly:
		{0, RxBwMant24 | 7<<RxBwExpShift, 1302},
		{1000, RxBwMant24 | 7<<RxBwExpShift, 1302},
		{48000, RxBwMant20 | 2<<RxBwExpShift, 50000},
		{150000, RxBwMant24 | 0<<RxBwExpShift, 166666},
		{112500, RxBwMant20 | 1<<RxBwExpShift, 100000},
		{300000, RxBwMant16 | 0<<RxBwExpShift, 250000},
	}
	for _, c := range cases {
		r := channelBwToRegister(c.bw)
		if r != c.r {
			t.Errorf("channelBwToRegister(%d) == %02X, want %02X", c.bw, r, c.r)
		}
		bw := registerToChannelBw(c.r, ModulationTypeOOK)
		if c.bwApprox != 0 {
			if bw != c.bwApprox {
				t.Errorf("registerToChannelBw(%02X) == %d, want %d", c.r, bw, c.bwApprox)
			}
		} else {
			if bw != c.bw {
				t.Errorf("registerToChannelBw(%02X) == %d, want %d", c.r, bw, c.bw)
			}
		}
	}
}
*/
