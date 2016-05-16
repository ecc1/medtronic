package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ecc1/cc1100"
)

var (
	signalChan = make(chan os.Signal, 1)
)

func main() {
	dev, err := cc1100.Open()
	if err != nil {
		log.Fatal(err)
	}

	err = dev.Reset()
	if err != nil {
		log.Fatal(err)
	}

	err = dev.InitRF()
	if err != nil {
		log.Fatal(err)
	}

	if cc1100.Verbose {
		fmt.Printf("\nRF settings after initialization:\n")
		dev.DumpRF()
	}

	signal.Notify(signalChan, os.Interrupt)
	go stats(dev)

	dev.StartRadio()
	for packet := range dev.IncomingPackets() {
		fmt.Printf("Received %d bytes (RSSI = %d)\n", len(packet.Data), packet.Rssi)
		fmt.Printf("Raw data: ")
		cc1100.PrintBytes(packet.Data)
		data, err := dev.DecodePacket(packet)
		if !cc1100.Verbose {
			if err == nil {
				cc1100.PrintBytes(data)
			} else {
				fmt.Printf("%v\n", err)
			}
		}
	}
}

func stats(dev *cc1100.Device) {
	tick := time.Tick(10 * time.Second)
	for {
		select {
		case <-tick:
			dev.PrintStats()
		case <-signalChan:
			dev.PrintStats()
			os.Exit(0)
		}
	}
}
