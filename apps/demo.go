package main

import (
	"fmt"
	"log"

	"github.com/ecc1/cc1100"
	"github.com/ecc1/spi"
)

func main() {
	dev, err := spi.Open(cc1100.SpiSpeed, 0)
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

	cc1100.Reset(dev)
	fmt.Printf("\nDefault RF settings:\n")
	cc1100.DumpRF(dev)

	cc1100.InitRF(dev)
	fmt.Printf("\nRF settings after initialization:\n")
	cc1100.DumpRF(dev)
}
