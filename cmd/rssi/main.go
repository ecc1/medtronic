package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ecc1/medtronic"
)

var (
	interval = flag.Int("i", 1, "ping interval in seconds")
	graph    = flag.Bool("g", false, "show bar graph of RSSI value")
)

func main() {
	flag.Parse()
	pump := medtronic.Open()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	pump.Wakeup()
	for {
		pump.Model()
		var rssi int
		if pump.Error() == nil {
			rssi = pump.RSSI()
		} else {
			fmt.Println(pump.Error())
			pump.SetError(nil)
			rssi = -128
		}
		showRSSI(rssi)
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}

func showRSSI(rssi int) {
	log.Printf("%4d\n", rssi)
	if !*graph {
		return
	}
	n := rssi + 128
	for i := 0; i < n; i++ {
		fmt.Print("━")
	}
	fmt.Printf("\n")
}
