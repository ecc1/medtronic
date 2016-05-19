package medtronic

import (
	"fmt"

	"github.com/ecc1/cc1100"
)

func (pump *Pump) DecodePacket(packet cc1100.Packet) ([]byte, error) {
	data, err := Decode6b4b(packet.Data)
	if err != nil {
		pump.DecodingErrors++
		return nil, err
	}
	crc := Crc8(data[:len(data)-1])
	if data[len(data)-1] != crc {
		pump.CrcErrors++
		return data, fmt.Errorf("CRC should be %X, not %X", crc, data[len(data)-1])
	}
	return data, nil
}

func EncodePacket(packet []byte) cc1100.Packet {
	return cc1100.Packet{Data: Encode4b6b(append(packet, Crc8(packet)))}
}

func (pump *Pump) PrintStats() {
	good := pump.Radio.PacketsReceived - pump.DecodingErrors - pump.CrcErrors
	fmt.Printf("\nTX: %6d    RX: %6d    decode errs: %6d    CRC errs: %6d\n", pump.Radio.PacketsSent, good, pump.DecodingErrors, pump.CrcErrors)
	s, _ := pump.Radio.ReadState()
	m, _ := pump.Radio.ReadMarcState()
	fmt.Printf("State: %s / %s\n", cc1100.StateName(s), cc1100.MarcStateName(m))
}
