package medtronic

import (
	"fmt"
)

func (pump *Pump) DecodePacket(packet []byte) []byte {
	data, err := Decode6b4b(packet)
	if err != nil {
		pump.err = err
		pump.DecodingErrors++
		return data
	}
	last := len(data) - 1
	pktCrc := data[last]
	data = data[:last] // without CRC
	calcCrc := Crc8(data)
	if pktCrc != calcCrc {
		pump.err = fmt.Errorf("CRC should be %X, not %X", calcCrc, pktCrc)
		pump.CrcErrors++
	}
	return data
}

func EncodePacket(data []byte) []byte {
	// Don't use append() to add the CRC, because append
	// may write into the array underlying the caller's slice.
	buf := make([]byte, len(data)+1)
	copy(buf, data)
	buf[len(data)] = Crc8(data)
	return Encode4b6b(buf)
}

func (pump *Pump) PrintStats() {
	stats := pump.Radio.Statistics()
	good := stats.Packets.Received - pump.DecodingErrors - pump.CrcErrors
	fmt.Printf("\nTX: %6d    RX: %6d    decode errs: %6d    CRC errs: %6d\n", stats.Packets.Sent, good, pump.DecodingErrors, pump.CrcErrors)
	fmt.Printf("State: %s\n", pump.Radio.State())
}
