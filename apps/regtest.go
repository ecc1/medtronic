package main

import (
	"fmt"
	"log"

	"github.com/ecc1/cc1100"
)

func main() {
	dev, err := cc1100.Open()
	if err != nil {
		log.Fatal(err)
	}

	dev.Reset()

	x, err := dev.ReadRegister(cc1100.SYNC0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Before write: %#X\n", x)

	err = dev.WriteRegister(cc1100.SYNC0, 0x44)
	if err != nil {
		log.Fatal(err)
	}
	x, err = dev.ReadRegister(cc1100.SYNC0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("After write: %#X\n", x)

	// 24-bit base frequency
	// 0x2340FC * 26MHz / 2^16 == 916.6MHz (916599975 Hz)
	err = dev.WriteEach([]byte{
		cc1100.FREQ2, 0x23,
		cc1100.FREQ1, 0x40,
		cc1100.FREQ0, 0xFC,
	})
	if err != nil {
		log.Fatal(err)
	}

	freq, err := dev.ReadFrequency()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Frequency: %d\n", freq)
}
