package main

import (
	"fmt"
	"log"

	"github.com/ecc1/cc1100"
	"github.com/ecc1/spi"
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

	fmt.Printf("\nRF settings after initialization:\n")
	cc1100.DumpRF(dev)

	for {
		packet, err := cc1100.ReceivePacket(dev)
		if err != nil {
			log.Fatal(err)
		}
		r, err := cc1100.ReadRSSI(dev)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Received %d bytes (RSSI = %d)\n", len(packet), r)
		printPacket("raw", packet)
		data, err := cc1100.Decode6b4b(packet)
		if err != nil {
			fmt.Printf("%v\n", err)
		} else {
			printPacket("decoded", data)
		}
	}
}

func printPacket(msg string, data []byte) {
	fmt.Printf("%s:\n", msg)
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
