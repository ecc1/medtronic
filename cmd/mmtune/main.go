package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/ecc1/medtronic"
	"github.com/ecc1/radio"
)

const (
	startFreq = uint32(916000000)
	endFreq   = uint32(917000000)
	precision = uint32(10000)
)

func main() {
	pump := medtronic.Open()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	defer pump.Close()
	f := searchFrequencies(pump)
	showResults(f)
	fmt.Println(radio.MegaHertz(f))
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

// Result represents the RSSI at a given frequency.
type Result struct {
	frequency uint32
	rssi      int
}

// Results implements sort.Interface based on frequency.
type Results []Result

func (r Results) Len() int           { return len(r) }
func (r Results) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r Results) Less(i, j int) bool { return r[i].frequency < r[j].frequency }

var results Results

func tryFrequency(pump *medtronic.Pump, freq uint32) int {
	const sampleSize = 2
	pump.Radio.SetFrequency(freq)
	log.Printf("frequency set to %s", radio.MegaHertz(freq))
	rssi := -128
	count := 0
	sum := 0
	for i := 0; i < sampleSize; i++ {
		pump.Model()
		if pump.Error() != nil {
			pump.SetError(nil)
			continue
		}
		sum += pump.RSSI()
		count++
	}
	if count != 0 {
		rssi = (sum + count/2) / count
	}
	results = append(results, Result{frequency: freq, rssi: rssi})
	return rssi
}

func showResults(winner uint32) {
	sort.Sort(results)
	for _, r := range results {
		fmt.Printf("%s  %4d ", radio.MegaHertz(r.frequency), r.rssi)
		n := r.rssi + 128
		for i := 0; i < n; i++ {
			fmt.Print("━")
		}
		if r.frequency == winner {
			fmt.Print(" ⏺")
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}
