package packet

import (
	"bytes"
	"testing"
)

func Test4b6bEncoding(t *testing.T) {
	cases := []struct {
		decoded []byte
		encoded []byte
	}{
		{parseBytes(""), parseBytes("")},
		{parseBytes("00"), parseBytes("55 55")},
		{parseBytes("00 00"), parseBytes("55 55 55")},
		{parseBytes("A7 12 89 86 5D 00 BE"), parseBytes("A9 6C 72 69 96 A6 94 D5 55 2C E5")},
		{parseBytes("A7 12 89 86 06 00 15"), parseBytes("A9 6C 72 69 96 A6 56 65 55 C6 55")},
		{parseBytes("A7 12 89 86 15 09 56"), parseBytes("A9 6C 72 69 96 A6 C6 55 59 96 65")},
		{parseBytes("A7 12 89 86 8D 00 B0"), parseBytes("A9 6C 72 69 96 A6 68 D5 55 2D 55")},
		{parseBytes("A7 12 89 86 8D 09 03 37 32 32 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 39"), parseBytes("A9 6C 72 69 96 A6 68 D5 59 56 38 D6 8F 28 F2 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 8D 95")},
	}
	for _, c := range cases {
		t.Run("encode", func(t *testing.T) {
			result := Encode4b6b(c.decoded)
			if !bytes.Equal(result, c.encoded) {
				t.Errorf("Encode4b6b(% X) == % X, want % X", c.decoded, result, c.encoded)
			}
		})
		t.Run("decode", func(t *testing.T) {
			result, err := Decode6b4b(c.encoded)
			if err != nil {
				t.Errorf("Decode6b4b(% X) == %v, want % X", c.encoded, err, c.decoded)
			} else if !bytes.Equal(result, c.decoded) {
				t.Errorf("Decode6b4b(% X) == % X, want % X", c.encoded, result, c.decoded)
			}
		})
	}
}
