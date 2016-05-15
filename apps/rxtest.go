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

	err = dev.ReceiveMode()
	if err != nil {
		log.Fatal(err)
	}
	for {
		packet := <-dev.IncomingPackets()
		r, err := dev.ReadRSSI()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Received %d bytes (RSSI = %d)\n", len(packet), r)
		fmt.Printf("Raw data: ")
		cc1100.PrintBytes(packet)
		if len(packet) == 0 {
			continue
		}
		data, err := cc1100.DecodePacket(packet)
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
			cc1100.PrintStats()
			printState(dev)
		case <-signalChan:
			cc1100.PrintStats()
			printState(dev)
			os.Exit(0)
		}
	}
}

func printState(dev *cc1100.Device) {
	s, _ := dev.ReadState()
	m, _ := dev.ReadMarcState()
	fmt.Printf("State: %s / %s\n", cc1100.StateName(s), cc1100.MarcStateName(m))
}
