package medtronic

import (
	"bytes"
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
		rep := marshalUint16(c.val)
		if !bytes.Equal(rep, c.rep) {
			t.Errorf("marshalUint16(%04X) == % X, want % X", c.val, rep, c.rep)
		}
	}
}

func TestTwoByteUintLE(t *testing.T) {
	cases := []struct {
		val uint16
		rep []byte
	}{
		{0x1234, []byte{0x34, 0x12}},
		{0, []byte{0x00, 0x00}},
		{math.MaxUint16, []byte{0xFF, 0xFF}},
	}
	for _, c := range cases {
		val := twoByteUintLE(c.rep)
		if val != c.val {
			t.Errorf("twoByteUintLE(% X) == %04X, want %04X", c.rep, val, c.val)
		}
		rep := marshalUint16LE(c.val)
		if !bytes.Equal(rep, c.rep) {
			t.Errorf("marshalUint16LE(%04X) == % X, want % X", c.val, rep, c.rep)
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

func TestTwoByteIntLE(t *testing.T) {
	cases := []struct {
		val int
		rep []byte
	}{
		{0x1234, []byte{0x34, 0x12}},
		{0, []byte{0x00, 0x00}},
		{256, []byte{0x00, 0x01}},
		{-1, []byte{0xFF, 0xFF}},
		{-256, []byte{0x00, 0xFF}},
		{math.MaxInt16, []byte{0xFF, 0x7F}},
		{math.MinInt16, []byte{0x00, 0x80}},
	}
	for _, c := range cases {
		val := twoByteIntLE(c.rep)
		if val != c.val {
			t.Errorf("twoByteIntLE(% X) == %d, want %d", c.rep, val, c.val)
		}
	}
}

func TestFourByteUint(t *testing.T) {
	cases := []struct {
		val uint32
		rep []byte
	}{
		{0, []byte{0x00, 0x00, 0x00, 0x00}},
		{1, []byte{0x00, 0x00, 0x00, 0x01}},
		{0x12345678, []byte{0x12, 0x34, 0x56, 0x78}},
		{math.MaxUint32, []byte{0xFF, 0xFF, 0xFF, 0xFF}},
	}
	for _, c := range cases {
		val := fourByteUint(c.rep)
		if val != c.val {
			t.Errorf("fourByteUint(% X) == %08X, want %08X", c.rep, val, c.val)
		}
		rep := marshalUint32(c.val)
		if !bytes.Equal(rep, c.rep) {
			t.Errorf("marshalUint32(%04X) == % X, want % X", c.val, rep, c.rep)
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
