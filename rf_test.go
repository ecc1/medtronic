package rfm69

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
		{315000000, []byte{0x4E, 0xC0, 0x00}, 0},
		{434000000, []byte{0x6C, 0x80, 0x00}, 0},
		{868000000, []byte{0xD9, 0x00, 0x00}, 0},
		{915000000, []byte{0xE4, 0xC0, 0x00}, 0},
		// some that can't be represented exactly:
		{916300000, []byte{0xE5, 0x13, 0x33}, 916299987},
		{916600000, []byte{0xE5, 0x26, 0x66}, 916599975},
	}
	for _, c := range cases {
		b := frequencyToRegisters(c.f)
		if !bytes.Equal(b, c.b) {
			t.Errorf("frequencyToRegisters(%d) == % X, want % X", c.f, b, c.b)
		}
		f := registersToFrequency(c.b)
		if c.fApprox != 0 {
			if f != c.fApprox {
				t.Errorf("registersToFrequency(% X) == %d, want %d", c.b, f, c.fApprox)
			}
		} else {
			if f != c.f {
				t.Errorf("registersToFrequency(% X) == %d, want %d", c.b, f, c.f)
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
