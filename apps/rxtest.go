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

	err = cc1100.ChangeState(dev, cc1100.SRX, cc1100.STATE_RX)
	if err != nil {
		log.Fatal(err)
	}

	for {
		packet, err := cc1100.ReceivePacket(dev)
		if err != nil {
			log.Fatal(err)
		}
		printPacket(packet)
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
