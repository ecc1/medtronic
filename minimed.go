package cc1100

import (
	"errors"
	"fmt"
)

const (
	printBinary = false
)

var (
	DecodingErrors = 0
	CrcErrors      = 0
	CrcMismatch    = errors.New("CRC mismatch")
)

func (dev *Device) DecodePacket(packet []byte) ([]byte, error) {
	data, err := Decode6b4b(packet)
	if err != nil {
		DecodingErrors++
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
		CrcErrors++
		if Verbose {
			fmt.Printf("CRC should be %X, not %X\n", crc, data[len(data)-1])
		}
		return data, CrcMismatch
	}
	return data, nil
}

func (dev *Device) EncodePacket(packet []byte) []byte {
	packet = append(packet, Crc8(packet))
	if Verbose {
		fmt.Printf("Packet:  ")
		PrintBytes(packet)
	}
	data := Encode4b6b(packet)
	if Verbose {
		fmt.Printf("Encoded: ")
		PrintBytes(data)
	}
	return data
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
	good := PacketsReceived - DecodingErrors - CrcErrors
	fmt.Printf("\nTX: %6d    RX: %6d    decode errs: %6d    CRC errs: %6d\n", PacketsSent, good, DecodingErrors, CrcErrors)
	s, _ := dev.ReadState()
	fmt.Printf("State: %s\n", StateName(s))
}
