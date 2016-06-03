package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/ecc1/medtronic"
)

const (
	verbose = true
)

var (
	signalChan = make(chan os.Signal, 1)
)

func main() {
	if verbose {
		log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	}
	pump, err := medtronic.Open()
	if err != nil {
		log.Fatal(err)
	}

	signal.Notify(signalChan, os.Interrupt)
	go catchInterrupt(pump)

	for packet := range pump.Radio.Incoming() {
		if verbose {
			log.Printf("raw data: % X (RSSI = %d)\n", packet.Data, packet.Rssi)
		}
		data, err := pump.DecodePacket(packet)
		if err != nil {
			log.Printf("%v\n", err)
			continue
		}
		if verbose {
			log.Printf(" decoded: % X\n", data)
		} else {
			log.Printf("% X (RSSI = %d)\n", data, packet.Rssi)
		}
	}
}

func catchInterrupt(pump *medtronic.Pump) {
	<-signalChan
	pump.PrintStats()
	os.Exit(0)
}
