package packet

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

func TestPacketEncoding(t *testing.T) {
	cases := []struct {
		decoded []byte // excluding CRC byte
		encoded []byte
	}{
		{parseBytes("A7 12 89 86 5D 00"), parseBytes("A9 6C 72 69 96 A6 94 D5 55 2C E5")},
		{parseBytes("A7 12 89 86 06 00"), parseBytes("A9 6C 72 69 96 A6 56 65 55 C6 55")},
		{parseBytes("A7 12 89 86 15 09"), parseBytes("A9 6C 72 69 96 A6 C6 55 59 96 65")},
		{parseBytes("A7 12 89 86 8D 00"), parseBytes("A9 6C 72 69 96 A6 68 D5 55 2D 55")},
		{parseBytes("A7 12 89 86 8D 09 03 37 32 32 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00"), parseBytes("A9 6C 72 69 96 A6 68 D5 59 56 38 D6 8F 28 F2 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 8D 95")},
		{parseBytes("A8 0F 25 C1 23 0D 19 1C 50 00 8F 00 90 00 34 34 99 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00"), parseBytes("A9 A5 5C CA 5B 31 CA 35 4D C5 9C 6C 95 55 55 69 C5 55 65 55 55 8F 48 F4 65 95 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 55 34 E3 96")},
	}
	for _, c := range cases {
		t.Run("encode", func(t *testing.T) {
			result := Encode(c.decoded)
			if !bytes.Equal(result, c.encoded) {
				t.Errorf("Encode(% X) == % X, want % X", c.decoded, result, c.encoded)
			}
		})
		t.Run("decode", func(t *testing.T) {
			result, err := Decode(c.encoded)
			if err != nil {
				t.Errorf("Decode(% X) == %v, want % X", c.encoded, err, c.decoded)
				return
			}
			if !bytes.Equal(result, c.decoded) {
				t.Errorf("Decode(% X) == % X, want % X", c.encoded, result, c.decoded)
			}
		})
	}
}

func parseBytes(hex string) []byte {
	fields := strings.Fields(hex)
	data := make([]byte, len(fields))
	for i, s := range fields {
		b, err := strconv.ParseUint(string(s), 16, 8)
		if err != nil {
			panic(err)
		}
		data[i] = byte(b)
	}
	return data
}
