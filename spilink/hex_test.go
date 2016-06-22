package main

import (
	"bytes"
	"testing"
)

func TestHex(t *testing.T) {
	cases := []struct {
		hex []byte
		str string
	}{
		{[]byte{}, ""},
		{[]byte{0xFF}, "ff"},
		{[]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}, "0123456789abcdef"},
	}
	for _, c := range cases {
		str := hexEncode(c.hex)
		if str != c.str {
			t.Errorf("hexEncode(% X) == %s, want %s", c.hex, str, c.str)
		}
		hex := hexDecode(c.str)
		if !bytes.Equal(hex, c.hex) {
			t.Errorf("hexDecode(%s) == % X, want % X", c.str, hex, c.hex)
		}
	}
}
