package main

import (
	"fmt"
	"log"

	"github.com/ecc1/cc1100"
)

func main() {
	r, err := cc1100.Open()
	if err != nil {
		log.Fatal(err)
	}
	r.Reset()

	dumpRegs(r)

	fmt.Printf("\nTesting individual writes\n")
	err = r.WriteRegister(cc1100.SYNC1, 0x44)
	if err != nil {
		log.Fatal(err)
	}
	err = r.WriteRegister(cc1100.SYNC0, 0x55)
	if err != nil {
		log.Fatal(err)
	}
	readRegs(r)

	r.Reset()
	fmt.Printf("\nTesting burst writes\n")
	err = r.WriteBurst(cc1100.SYNC1, []byte{0x66, 0x77})
	if err != nil {
		log.Fatal(err)
	}
	readRegs(r)
}

func dumpRegs(r *cc1100.Radio) {
	start := cc1100.IOCFG2
	finish := cc1100.TEST0
	numRegs := finish - start + 1
	fmt.Printf("Register dump:\n")
	regs, err := r.ReadBurst(byte(start), numRegs)
	if err != nil {
		log.Fatal(err)
	}
	for i, v := range regs {
		fmt.Printf("%02X  %02X  %08b\n", start+i, v, v)
	}
}

func readRegs(r *cc1100.Radio) {
	x, err := r.ReadRegister(cc1100.SYNC1)
	if err != nil {
		log.Fatal(err)
	}
	y, err := r.ReadRegister(cc1100.SYNC0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("individual: %X %X\n", x, y)
	v, err := r.ReadBurst(cc1100.SYNC1, 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  burst:    %X %X\n", v[0], v[1])
}
