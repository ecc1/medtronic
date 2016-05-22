package main

import (
	"fmt"
	"log"

	"github.com/ecc1/rfm69"
)

func main() {
	r, err := rfm69.Open()
	if err != nil {
		log.Fatal(err)
	}
	r.Reset()

	dumpRegs(r)

	fmt.Printf("\nTesting individual writes\n")
	err = r.WriteRegister(rfm69.RegSyncValue1, 0x44)
	if err != nil {
		log.Fatal(err)
	}
	err = r.WriteRegister(rfm69.RegSyncValue2, 0x55)
	if err != nil {
		log.Fatal(err)
	}
	err = r.WriteRegister(rfm69.RegSyncValue3, 0x66)
	if err != nil {
		log.Fatal(err)
	}
	readRegs(r)

	r.Reset()
	fmt.Printf("\nTesting burst writes\n")
	err = r.WriteBurst(rfm69.RegSyncValue1, []byte{0x44,0x55, 0x66})
	if err != nil {
		log.Fatal(err)
	}
	readRegs(r)

}

func dumpRegs(r *rfm69.Radio) {
	start := rfm69.RegOpMode
	finish := rfm69.RegTemp2
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

func readRegs(r *rfm69.Radio) {
	x, err := r.ReadRegister(rfm69.RegSyncValue1)
	if err != nil {
		log.Fatal(err)
	}
	y, err := r.ReadRegister(rfm69.RegSyncValue2)
	if err != nil {
		log.Fatal(err)
	}
	z, err := r.ReadRegister(rfm69.RegSyncValue3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("individual: %X %X %X\n", x, y, z)
	v, err := r.ReadBurst(rfm69.RegSyncValue1, 3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  burst:    %X %X %X\n", v[0], v[1], v[2])
}
