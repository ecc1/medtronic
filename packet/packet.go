package packet

import (
	"fmt"
)

// Decode performs 6b/4b decoding and CRC verification.
// It returns the decoded data, excluding the CRC byte.
func Decode(p []byte) ([]byte, error) {
	data, err := Decode6b4b(p)
	if err != nil {
		return data, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("empty packet")
	}
	if data[0] == Sensor {
		return checkCRC16(data)
	}
	return checkCRC8(data)
}

func checkCRC8(data []byte) ([]byte, error) {
	last := len(data) - 1
	pktCRC := data[last]
	data = data[:last] // without CRC
	calcCRC := CRC8(data)
	if calcCRC != pktCRC {
		return nil, fmt.Errorf("computed CRC %02X but received %02X", calcCRC, pktCRC)
	}
	return data, nil
}

func checkCRC16(data []byte) ([]byte, error) {
	last := len(data) - 2
	pktCRC := uint16(data[last])<<8 | uint16(data[last+1])
	data = data[:last] // without CRC
	calcCRC := CRC16(data)
	if calcCRC != pktCRC {
		return nil, fmt.Errorf("computed CRC %04X but received %04X", calcCRC, pktCRC)
	}
	return data, nil
}

// Encode appends the CRC to the data and returns the 4b/6b-encoded result.
// This may modify data's underlying array.
func Encode(data []byte) []byte {
	var msg []byte
	if data[0] == Sensor {
		crc := CRC16(data)
		msg = append(data, byte(crc>>8), byte(crc&0xFF))
	} else {
		msg = append(data, CRC8(data))
	}
	return Encode4b6b(msg)
}
