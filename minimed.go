package cc1100

import (
	"errors"
	"fmt"
)

const (
	printBinary = false
)

var (
	CrcMismatch = errors.New("CRC mismatch")
)

func (dev *Device) DecodePacket(packet Packet) ([]byte, error) {
	data, err := Decode6b4b(packet.Data)
	if err != nil {
		dev.decodingErrors++
		if Verbose {
			fmt.Printf("%v\n", err)
		}
		return nil, err
	}
	if Verbose {
		fmt.Printf("Decoded:  ")
		PrintBytes(data)
	}
	crc := Crc8(data[:len(data)-1])
	if data[len(data)-1] != crc {
		dev.crcErrors++
		if Verbose {
			fmt.Printf("CRC should be %X, not %X\n", crc, data[len(data)-1])
		}
		return data, CrcMismatch
	}
	return data, nil
}

func EncodePacket(packet []byte) Packet {
	return Packet{Data: Encode4b6b(append(packet, Crc8(packet)))}
}

func PrintBytes(data []byte) {
	fmt.Printf("% X\n", data)
	if !printBinary {
		return
	}
	for i, v := range data {
		fmt.Printf("%08b", v)
		if (i+1)%10 == 0 {
			fmt.Print("\n")
		}
	}
	if len(data)%10 != 0 {
		fmt.Print("\n")
	}
}

func (dev *Device) PrintStats() {
	good := dev.packetsReceived - dev.decodingErrors - dev.crcErrors
	fmt.Printf("\nTX: %6d    RX: %6d    decode errs: %6d    CRC errs: %6d\n", dev.packetsSent, good, dev.decodingErrors, dev.crcErrors)
	s, _ := dev.ReadState()
	m, _ := dev.ReadMarcState()
	fmt.Printf("State: %s / %s\n", StateName(s), MarcStateName(m))
}
