package main

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
	err = r.WriteRegister(cc1100.IOCFG0, 0x39)
	if err != nil {
		log.Fatal(err)
	}
}
