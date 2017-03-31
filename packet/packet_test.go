package packet

import (
	"bytes"
	"testing"
)

func TestPacketEncoding(t *testing.T) {
	cases := []struct {
		decoded []byte
		encoded []byte
	}{
		{[]byte{0x00, 0x00}, []byte{0x55, 0x55, 0x55}},
		{[]byte{0xA7, 0x12, 0x34, 0x56, 0x8D, 0x00, 0xA6, 0x00}, []byte{0xA9, 0x6C, 0x72, 0x8F, 0x49, 0x66, 0x68, 0xD5, 0x55, 0xAA, 0x63, 0x4E}},
	}
	for _, c := range cases {
		result := Encode(c.decoded)
		if !bytes.Equal(result, c.encoded) {
			t.Errorf("Encode(% X) == % X, want % X", c.decoded, result, c.encoded)
		}
		result, err := Decode(c.encoded)
		if err != nil {
			t.Errorf("Decode(% X) == %v, want % X", c.encoded, err, c.decoded)
			continue
		}
		d := c.decoded[:len(c.decoded)-1]
		if !bytes.Equal(result, d) {
			t.Errorf("Decode(% X) == % X, want % X", c.encoded, result, d)
		}
	}
}
