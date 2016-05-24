package main

import (
	"fmt"
	"log"

	"github.com/ecc1/medtronic"
)

const (
	startFreq  = uint32(916550000)
	endFreq    = uint32(916770000)
	stepSize   = uint32(10000)
	sampleSize = 5
)

func main() {
	pump, err := medtronic.Open()
	if err != nil {
		log.Fatal(err)
	}
	printResults(scanFrequencies(pump))
}

type Result struct {
	frequency uint32
	count     int
	rssi      int
}

func printResults(results []Result) {
	for _, r := range results {
		fmt.Printf("%s  %d  %3d\n", shortFreq(r.frequency), r.count, r.rssi)
	}
}

func shortFreq(freq uint32) string {
	MHz := freq / 1000000
	kHz := (freq % 1000000) / 1000
	return fmt.Sprintf("%3d.%03d", MHz, kHz)
}

func scanFrequencies(pump *medtronic.Pump) []Result {
	var results []Result
	for freq := startFreq; freq <= endFreq; freq += stepSize {
		results = append(results, tryFrequency(pump, freq))
	}
	return results
}

func tryFrequency(pump *medtronic.Pump, freq uint32) Result {
	err := pump.Radio.SetFrequency(freq)
	if err != nil {
		log.Fatal(err)
	}
	f, err := pump.Radio.Frequency()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("frequency set to %s MHz\n", shortFreq(f))
	count := 0
	rssiSum := 0
	rssi := -99
	for i := 0; i < sampleSize; i++ {
		r := -99
		_, err := pump.Model(1, &r)
		if err != nil {
			continue
		}
		count++
		rssiSum += r
	}
	if count != 0 {
		rssi = (rssiSum + count/2) / count
	}
	return Result{frequency: freq, count: count, rssi: rssi}
}
