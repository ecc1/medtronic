package packet

import (
	"fmt"
)

// DecodePacket performs 4b/6b decoding and CRC verification.
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
	if pktCRC != calcCRC {
		err = fmt.Errorf("CRC should be %X, not %X", calcCRC, pktCRC)
	}
	return data, err
}

// EncodePacket calculates and stores the final CRC byte
// and returns the 4b/6b-encoded result.
// The caller must provide space for the CRC byte.
func Encode(data []byte) []byte {
	n := len(data) - 1
	data[n] = CRC8(data[:n])
	return Encode4b6b(data)
}
