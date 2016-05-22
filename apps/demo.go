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

	log.Printf("Initializing radio\n")
	err = r.InitRF()
	if err != nil {
		log.Fatal(err)
	}
	dumpRF(r)

	freq := uint32(915650000)
	log.Printf("Changing frequency to %d\n", freq)
	err = r.WriteFrequency(freq)
	if err != nil {
		log.Fatal(err)
	}
	dumpRF(r)
}

func dumpRF(r *rfm69.Radio) {
	mode, err := r.Mode()
	if err != nil {
		log.Fatal(err)
	}
	s := ""
	switch mode {
	case rfm69.SleepMode:
		s = "Sleep"
	case rfm69.StandbyMode:
		s = "Standby"
	case rfm69.FreqSynthMode:
		s = "Frequency Synthesizer"
	case rfm69.TransmitterMode:
		s = "Transmitter"
	case rfm69.ReceiverMode:
		s = "Receiver"
	default:
		log.Panicf("Unknown operating mode (%X)\n", mode)
	}
	log.Printf("Mode: %s\n", s)

	freq, err := r.ReadFrequency()
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
