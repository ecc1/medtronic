package main

import (
	"fmt"
	"log"

	"github.com/ecc1/medtronic"
)

const (
	startFreq = uint32(916500000)
	endFreq   = uint32(916800000)
	precision = uint32(10000)
)

func main() {
	pump, err := medtronic.Open()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(shortFreq(searchFrequencies(pump)))
}

func shortFreq(freq uint32) string {
	MHz := freq / 1000000
	kHz := (freq % 1000000) / 1000
	return fmt.Sprintf("%3d.%03d", MHz, kHz)
}

// Use ternary search to find frequency with maximum RSSI.
func searchFrequencies(pump *medtronic.Pump) uint32 {
	pump.SetRetries(1)
	lower := startFreq
	upper := endFreq
	for {
		delta := upper - lower
		if delta < precision {
			return (lower + upper) / 2
		}
		delta /= 3
		lowerThird := lower + delta
		r1 := tryFrequency(pump, lowerThird)
		upperThird := upper - delta
		r2 := tryFrequency(pump, upperThird)
		if r1 < r2 {
			lower = lowerThird
		} else {
			upper = upperThird
		}
	}
}

func tryFrequency(pump *medtronic.Pump, freq uint32) int {
	err := pump.Radio.SetFrequency(freq)
	if err != nil {
		log.Fatal(err)
	}
	f, err := pump.Radio.Frequency()
	if err != nil {
		log.Fatal(err)
	}
	_, err = pump.Model()
	rssi := 0
	if err == nil {
		rssi = pump.Rssi()
	} else {
		rssi = -99
	}
	log.Printf("%s MHz: %d\n", shortFreq(f), rssi)
	return rssi
}
