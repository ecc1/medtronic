package medtronic

import (
	"math"
	"testing"
)

func TestTwoByteUint(t *testing.T) {
	cases := []struct {
		val uint16
		rep []byte
	}{
		{0x1234, []byte{0x12, 0x34}},
		{0, []byte{0x00, 0x00}},
		{math.MaxUint16, []byte{0xFF, 0xFF}},
	}
	for _, c := range cases {
		val := twoByteUint(c.rep)
		if val != c.val {
			t.Errorf("twoByteUint(% X) == %04X, want %04X", c.rep, val, c.val)
		}
	}
}

func TestTwoByteInt(t *testing.T) {
	cases := []struct {
		val int
		rep []byte
	}{
		{0x1234, []byte{0x12, 0x34}},
		{0, []byte{0x00, 0x00}},
		{256, []byte{0x01, 0x00}},
		{-1, []byte{0xFF, 0xFF}},
		{-256, []byte{0xFF, 0x00}},
		{math.MaxInt16, []byte{0x7F, 0xFF}},
		{math.MinInt16, []byte{0x80, 0x00}},
	}
	for _, c := range cases {
		val := twoByteInt(c.rep)
		if val != c.val {
			t.Errorf("twoByteInt(% X) == %d, want %d", c.rep, val, c.val)
		}
	}
}

func TestFourByteUint(t *testing.T) {
	cases := []struct {
		val uint32
		rep []byte
	}{
		{0x12345678, []byte{0x12, 0x34, 0x56, 0x78}},
		{0, []byte{0x00, 0x00, 0x00, 0x00}},
		{math.MaxUint32, []byte{0xFF, 0xFF, 0xFF, 0xFF}},
	}
	for _, c := range cases {
		val := fourByteUint(c.rep)
		if val != c.val {
			t.Errorf("fourByteUint(% X) == %08X, want %08X", c.rep, val, c.val)
		}
	}
}

func TestFourByteInt(t *testing.T) {
	cases := []struct {
		val int
		rep []byte
	}{
		{0x12345678, []byte{0x12, 0x34, 0x56, 0x78}},
		{0, []byte{0x00, 0x00, 0x00, 0x00}},
		{-1, []byte{0xFF, 0xFF, 0xFF, 0xFF}},
		{0x0000FFFF, []byte{0x00, 0x00, 0xFF, 0xFF}},
		{-0x10000, []byte{0xFF, 0xFF, 0x00, 0x00}},
		{math.MaxInt32, []byte{0x7F, 0xFF, 0xFF, 0xFF}},
		{math.MinInt32, []byte{0x80, 0x00, 0x00, 0x00}},
	}
	for _, c := range cases {
		val := fourByteInt(c.rep)
		if val != c.val {
			t.Errorf("fourByteInt(% X) == %d, want %d", c.rep, val, c.val)
		}
	}
}
