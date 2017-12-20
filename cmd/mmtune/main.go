package main

import (
	"flag"
	"fmt"
	"log"
	"sort"

	"github.com/ecc1/medtronic"
	"github.com/ecc1/radio"
)

var (
	start     = flag.String("f", "916.300", "scan from this `freq`uency")
	end       = flag.String("t", "916.900", "scan to this `freq`uency")
	delta     = flag.Int("k", 50, "`step` size in kHz")
	worldWide = flag.Bool("ww", false, "scan worldwide frequencies (868 MHz band)")
	showGraph = flag.Bool("g", false, "print graph instead of JSON")

	startFreq uint32
	endFreq   uint32
)

func main() {
	flag.Parse()
	if flag.NArg() != 0 {
		flag.Usage()
		return
	}
	if *worldWide {
		*start = "868.150"
		*end = "868.750"
	}
	var err error
	startFreq, err = medtronic.ParseFrequency(*start)
	if err != nil {
		flag.Usage()
		log.Fatal(err)
	}
	endFreq, err = medtronic.ParseFrequency(*end)
	if err != nil {
		flag.Usage()
		log.Fatal(err)
	}
	pump := medtronic.Open()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	defer pump.Close()
	pump.Wakeup()
	f := searchFrequencies(pump)
	sort.Sort(results)
	if *showGraph {
		showResults(f)
	} else {
		showJSON(f)
	}
}

// Find frequency with maximum RSSI.
func searchFrequencies(pump *medtronic.Pump) uint32 {
	pump.SetRetries(1)
	maxRSSI := -128
	bestFreq := startFreq
	deltaHz := uint32(*delta) * 1000
	for f := startFreq; f <= endFreq; f += deltaHz {
		rssi := tryFrequency(pump, f)
		if rssi > maxRSSI {
			maxRSSI = rssi
			bestFreq = f
		}
	}
	return bestFreq
}

// Result represents the RSSI at a given frequency.
type Result struct {
	frequency uint32
	rssi      int
	count     int
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
	results = append(results, Result{frequency: freq, rssi: rssi, count: count})
	return rssi
}

func showResults(winner uint32) {
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
	fmt.Println(radio.MegaHertz(winner))
}
