package main

import (
	"fmt"
	"log"

	"github.com/ecc1/cc1100"
	"github.com/ecc1/spi"
)

const (
	verbose = false
)

func main() {
	dev, err := spi.Open(cc1100.SpiSpeed)
	if err != nil {
		log.Fatal(err)
	}
	err = dev.SetMaxSpeed(cc1100.SpiSpeed)
	if err != nil {
		log.Fatal(err)
	}
	speed, err := dev.MaxSpeed()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Max speed: %d Hz\n", speed)

	err = cc1100.Reset(dev)
	if err != nil {
		log.Fatal(err)
	}

	err = cc1100.InitRF(dev)
	if err != nil {
		log.Fatal(err)
	}

	if verbose {
		fmt.Printf("\nRF settings after initialization:\n")
		cc1100.DumpRF(dev)
	}

	for {
		packet, err := cc1100.ReceivePacket(dev)
		if err != nil {
			log.Fatal(err)
		}
		r, err := cc1100.ReadRSSI(dev)
		if err != nil {
			log.Fatal(err)
		}
		if verbose {
			fmt.Printf("Received %d bytes (RSSI = %d)\n", len(packet), r)
			printPacket(packet)
		}
		data, err := cc1100.Decode6b4b(packet)
		if err != nil {
			if verbose {
				fmt.Printf("%v\n", err)
			}
			continue
		}
		if verbose {
			printPacket(data)
		}
		crc := cc1100.Crc8(data[:len(data)-1])
		if data[len(data)-1] != crc {
			if verbose {
				fmt.Printf("CRC should be %02X, not %02X\n", crc, data[len(data)-1])
			}
			continue
		}
		if !verbose {
			printPacket(data)
		}
	}
}

func printPacket(data []byte) {
	for i, v := range data {
		fmt.Printf("%02X ", v)
		if (i+1)%20 == 0 {
			fmt.Print("\n")
		}
	}
	if len(data)%20 != 0 {
		fmt.Print("\n")
	}
}
