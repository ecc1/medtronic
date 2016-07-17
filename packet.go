package medtronic

import (
	"fmt"
)

// DecodePacket performs 4b/6b decoding and CRC verification.
// It returns the decoded data, excluding the CRC byte.
func (pump *Pump) DecodePacket(packet []byte) []byte {
	data, err := Decode6b4b(packet)
	if err != nil || len(data) == 0 {
		if len(data) == 0 {
			err = fmt.Errorf("empty packet")
		}
		pump.SetError(err)
		pump.DecodingErrors++
		return data
	}
	last := len(data) - 1
	pktCrc := data[last]
	data = data[:last] // without CRC
	calcCrc := Crc8(data)
	if pktCrc != calcCrc {
		pump.SetError(fmt.Errorf("CRC should be %X, not %X", calcCrc, pktCrc))
		pump.CrcErrors++
	}
	return data
}

// EncodePacket calculates and stores the final CRC byte
// and returns the 4b/6b-encoded result.
// The caller must provide space for the CRC byte.
func EncodePacket(data []byte) []byte {
	n := len(data) - 1
	data[n] = Crc8(data[:n])
	return Encode4b6b(data)
}

func (pump *Pump) PrintStats() {
	stats := pump.Radio.Statistics()
	good := stats.Packets.Received - pump.DecodingErrors - pump.CrcErrors
	fmt.Printf("\nTX: %6d    RX: %6d    decode errs: %6d    CRC errs: %6d\n", stats.Packets.Sent, good, pump.DecodingErrors, pump.CrcErrors)
}
