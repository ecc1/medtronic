package main

import (
	"fmt"
	"log"

	"github.com/ecc1/cc1101"
)

func main() {
	r, err := cc1101.Open()
	if err != nil {
		log.Fatal(err)
	}
	r.Reset()

	dumpRegs(r)

	fmt.Printf("\nTesting individual writes\n")
	err = r.WriteRegister(cc1101.SYNC1, 0x44)
	if err != nil {
		log.Fatal(err)
	}
	err = r.WriteRegister(cc1101.SYNC0, 0x55)
	if err != nil {
		log.Fatal(err)
	}
	readRegs(r)

	r.Reset()
	fmt.Printf("\nTesting burst writes\n")
	err = r.WriteBurst(cc1101.SYNC1, []byte{0x66, 0x77})
	if err != nil {
		log.Fatal(err)
	}
	readRegs(r)
}

func dumpRegs(r *cc1101.Radio) {
	fmt.Printf("\nConfiguration registers:\n")
	config, err := r.ReadConfiguration()
	if err != nil {
		log.Fatal(err)
	}
	regs := config.Bytes()
	resetValue := cc1101.ResetRfConfiguration.Bytes()
	for i, v := range regs {
		fmt.Printf("%02X  %02X  %08b", cc1101.IOCFG2+i, v, v)
		r := resetValue[i]
		if v == r {
			fmt.Printf("\n")
		} else {
			fmt.Printf("  **** SHOULD BE %02X  %08b\n", r, r)
		}
	}
}

func readRegs(r *cc1101.Radio) {
	x, err := r.ReadRegister(cc1101.SYNC1)
	if err != nil {
		log.Fatal(err)
	}
	y, err := r.ReadRegister(cc1101.SYNC0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("individual: %X %X\n", x, y)
	v, err := r.ReadBurst(cc1101.SYNC1, 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  burst:    %X %X\n", v[0], v[1])
}
