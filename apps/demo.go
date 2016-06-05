package main

import (
	"log"

	"github.com/ecc1/rfm69"
)

func main() {
	r, err := rfm69.Open()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Resetting radio\n")
	err = r.Reset()
	if err != nil {
		log.Fatal(err)
	}
	dumpRF(r)

	freq := uint32(916600000)
	log.Println("")
	log.Printf("Initializing radio to %d MHz\n", freq)
	err = r.InitRF(freq)
	if err != nil {
		log.Fatal(err)
	}
	dumpRF(r)

	log.Println("")
	freq += 500000
	log.Printf("Changing frequency to %d\n", freq)
	err = r.SetFrequency(freq)
	if err != nil {
		log.Fatal(err)
	}
	dumpRF(r)

	log.Println("")
	log.Printf("Sleeping\n")
	err = r.Sleep()
	if err != nil {
		log.Fatal(err)
	}
	dumpRF(r)

}

func dumpRF(r *rfm69.Radio) {
	log.Printf("Mode: %s\n", r.State())

	freq, err := r.Frequency()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Frequency: %d Hz\n", freq)

	mod, err := r.ReadModulationType()
	if err != nil {
		log.Fatal(err)
	}
	switch mod {
	case rfm69.ModulationTypeFSK:
		log.Printf("Modulation type: FSK\n")
	case rfm69.ModulationTypeOOK:
		log.Printf("Modulation type: OOK\n")
	default:
		log.Panicf("Unknown modulation mode %X\n", mod)
	}

	bitrate, err := r.ReadBitrate()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Bitrate: %d baud\n", bitrate)

	bw, err := r.ReadChannelBw()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Channel BW: %d Hz\n", bw)
}
