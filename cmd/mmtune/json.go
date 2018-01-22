package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ecc1/radio"
)

// JSONResults is used to produce JSON output compatible with openaps.
type JSONResults struct {
	ScanDetails []interface{} `json:"scanDetails"`
	SetFreq     float64       `json:"setFreq"`
	UsedDefault bool          `json:"usedDefault"`
}

func showJSON(winner uint32, usedDefault bool) {
	j := JSONResults{
		ScanDetails: make([]interface{}, len(results)),
		SetFreq:     float64(winner) / 1000000,
		UsedDefault: usedDefault,
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
