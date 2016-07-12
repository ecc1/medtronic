package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ecc1/medtronic"
)

const (
	verbose = true
)

func main() {
	if verbose {
		log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	}
	pump := medtronic.Open()
	go catchInterrupt(pump)
	for pump.Error() == nil {
		packet, rssi := pump.Radio.Receive(time.Hour)
		if verbose {
			log.Printf("raw data: % X (RSSI = %d)", packet, rssi)
		}
		data := pump.DecodePacket(packet)
		if pump.Error() != nil {
			log.Printf("%v", pump.Error())
			pump.SetError(nil)
			continue
		}
		if verbose {
			log.Printf("decoded:  % X", data)
		} else {
			log.Printf("% X (RSSI = %d)", data, rssi)
		}

	}
	log.Fatal(pump.Error())
}

func catchInterrupt(pump *medtronic.Pump) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	pump.PrintStats()
	os.Exit(0)
}
