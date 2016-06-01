package main

// Route the CC1101 clock oscillator (divided by 24) to the GDO0 pin,
// so that its frequency can be measured easily with an oscilloscope
// or frequency counter.

// IMPORTANT: disconnect the Edison GPIO used for interrupts before
// running this program.  Otherwise the Edison will be interrupted
// at 1 MHz and become non-responsive.

import (
	"log"

	"github.com/ecc1/cc1100"
)

func main() {
	r, err := cc1100.Open()
	if err != nil {
		log.Fatal(err)
	}
	r.Reset()
	// Route CLK_XOSC/24 to GDO0 pin.
	// See data sheet, Table 41.
	err = r.WriteRegister(cc1100.IOCFG0, 0x39)
	if err != nil {
		log.Fatal(err)
	}
}
