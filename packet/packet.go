package packet

import (
	"fmt"
)

// Decode performs 6b/4b decoding and CRC verification.
// It returns the decoded data, excluding the CRC byte.
func Decode(p []byte) ([]byte, error) {
	data, err := Decode6b4b(p)
	if len(data) == 0 {
		err = fmt.Errorf("empty packet")
	}
	if err != nil {
		return data, err
	}
	last := len(data) - 1
	pktCRC := data[last]
	data = data[:last] // without CRC
	calcCRC := CRC8(data)
	if calcCRC != pktCRC {
		err = fmt.Errorf("computed CRC %02X but received %02X", calcCRC, pktCRC)

	}
	return data, err
}

// Encode appends the CRC to the data and returns the 4b/6b-encoded result.
// This may modify data's underlying array.
func Encode(data []byte) []byte {
	return Encode4b6b(append(data, CRC8(data)))
}
