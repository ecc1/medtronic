package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"sort"

	"github.com/ecc1/medtronic"
	"github.com/ecc1/radio"
)

var (
	numSteps  = flag.Int("n", 24, "number of frequencies to scan")
	worldWide = flag.Bool("ww", false, "scan worldwide frequencies (868 MHz band)")
	sFlag     = flag.String("f", "916.300", "scan from this frequency")
	eFlag     = flag.String("t", "916.900", "scan to this frequency")
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
		*sFlag = "868.150"
		*eFlag = "868.750"
	}
	var err error
	startFreq, err = medtronic.ParseFrequency(*sFlag)
	if err != nil {
		flag.Usage()
		log.Fatal(err)
	}
	endFreq, err = medtronic.ParseFrequency(*eFlag)
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
	deltaHz := (endFreq - startFreq) / uint32(*numSteps)
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

// JSONResults is used to produce JSON output compatible with openaps.
type JSONResults struct {
	ScanDetails []interface{} `json:"scanDetails"`
	SetFreq     float64       `json:"setFreq"`
	UsedDefault bool          `json:"usedDefault"`
}

func showJSON(winner uint32) {
	j := JSONResults{
		ScanDetails: make([]interface{}, len(results)),
		SetFreq:     float64(winner) / 1000000,
		UsedDefault: winner == startFreq,
	}
	// Convert each Result struct into a slice of interfaces
	// so it will be marshaled as a JSON array.
	for i, r := range results {
		j.ScanDetails[i] = []interface{}{
			radio.MegaHertz(r.frequency),
			r.count,
			r.rssi,
		}
	}
	b, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

}
