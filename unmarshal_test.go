package medtronic

import (
	"bytes"
	"fmt"
	"math"
	"testing"
)

func TestTwoByteUint(t *testing.T) {
	cases := []struct {
		val uint16
		rep []byte
	}{
		{0x1234, parseBytes("12 34")},
		{0, parseBytes("00 00")},
		{math.MaxUint16, parseBytes("FF FF")},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("unmarshal_%d", c.val), func(t *testing.T) {
			val := twoByteUint(c.rep)
			if val != c.val {
				t.Errorf("twoByteUint(% X) == %04X, want %04X", c.rep, val, c.val)
			}
		})
		t.Run(fmt.Sprintf("marshal_%d", c.val), func(t *testing.T) {
			rep := marshalUint16(c.val)
			if !bytes.Equal(rep, c.rep) {
				t.Errorf("marshalUint16(%04X) == % X, want % X", c.val, rep, c.rep)
			}
		})
	}
}

func TestTwoByteUintLE(t *testing.T) {
	cases := []struct {
		val uint16
		rep []byte
	}{
		{0x1234, parseBytes("34 12")},
		{0, parseBytes("00 00")},
		{math.MaxUint16, parseBytes("FF FF")},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("unmarshal_%d", c.val), func(t *testing.T) {
			val := twoByteUintLE(c.rep)
			if val != c.val {
				t.Errorf("twoByteUintLE(% X) == %04X, want %04X", c.rep, val, c.val)
			}
		})
		t.Run(fmt.Sprintf("marshal_%d", c.val), func(t *testing.T) {
			rep := marshalUint16LE(c.val)
			if !bytes.Equal(rep, c.rep) {
				t.Errorf("marshalUint16LE(%04X) == % X, want % X", c.val, rep, c.rep)
			}
		})
	}
}

func TestTwoByteInt(t *testing.T) {
	cases := []struct {
		val int
		rep []byte
	}{
		{0x1234, parseBytes("12 34")},
		{0, parseBytes("00 00")},
		{256, parseBytes("01 00")},
		{-1, parseBytes("FF FF")},
		{-256, parseBytes("FF 00")},
		{math.MaxInt16, parseBytes("7F FF")},
		{math.MinInt16, parseBytes("80 00")},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%d", c.val), func(t *testing.T) {
			val := twoByteInt(c.rep)
			if val != c.val {
				t.Errorf("twoByteInt(% X) == %d, want %d", c.rep, val, c.val)
			}
		})
	}
}

func TestTwoByteIntLE(t *testing.T) {
	cases := []struct {
		val int
		rep []byte
	}{
		{0x1234, parseBytes("34 12")},
		{0, parseBytes("00 00")},
		{256, parseBytes("00 01")},
		{-1, parseBytes("FF FF")},
		{-256, parseBytes("00 FF")},
		{math.MaxInt16, parseBytes("FF 7F")},
		{math.MinInt16, parseBytes("00 80")},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%d", c.val), func(t *testing.T) {
			val := twoByteIntLE(c.rep)
			if val != c.val {
				t.Errorf("twoByteIntLE(% X) == %d, want %d", c.rep, val, c.val)
			}
		})
	}
}

func TestFourByteUint(t *testing.T) {
	cases := []struct {
		val uint32
		rep []byte
	}{
		{0, parseBytes("00 00 00 00")},
		{1, parseBytes("00 00 00 01")},
		{0x12345678, parseBytes("12 34 56 78")},
		{math.MaxUint32, parseBytes("FF FF FF FF")},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("unmarshal_%d", c.val), func(t *testing.T) {
			val := fourByteUint(c.rep)
			if val != c.val {
				t.Errorf("fourByteUint(% X) == %08X, want %08X", c.rep, val, c.val)
			}
		})
		t.Run(fmt.Sprintf("marshal_%d", c.val), func(t *testing.T) {
			rep := marshalUint32(c.val)
			if !bytes.Equal(rep, c.rep) {
				t.Errorf("marshalUint32(%04X) == % X, want % X", c.val, rep, c.rep)
			}
		})
	}
}

func TestFourByteInt(t *testing.T) {
	cases := []struct {
		val int
		rep []byte
	}{
		{0x12345678, parseBytes("12 34 56 78")},
		{0, parseBytes("00 00 00 00")},
		{-1, parseBytes("FF FF FF FF")},
		{0x0000FFFF, parseBytes("00 00 FF FF")},
		{-0x10000, parseBytes("FF FF 00 00")},
		{math.MaxInt32, parseBytes("7F FF FF FF")},
		{math.MinInt32, parseBytes("80 00 00 00")},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%d", c.val), func(t *testing.T) {
			val := fourByteInt(c.rep)
			if val != c.val {
				t.Errorf("fourByteInt(% X) == %d, want %d", c.rep, val, c.val)
			}
		})
	}
}
