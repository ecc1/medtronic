package medtronic

import (
	"fmt"
	"log"
	"os"

	"github.com/ecc1/radio"
	"github.com/ecc1/rfm69"
)

const (
	DefaultFrequency = 916600000
	freqEnvVar       = "MEDTRONIC_FREQUENCY"
)

type Pump struct {
	Radio radio.Interface

	DecodingErrors int
	CrcErrors      int
}

func Open() (*Pump, error) {
	r, err := rfm69.Open()
	if err != nil {
		return nil, err
	}
	err = r.Init()
	if err != nil {
		return nil, err
	}
	freq := defaultFreq()
	log.Printf("setting frequency to %d\n", freq)
	err = r.SetFrequency(freq)
	if err != nil {
		return nil, err
	}
	return &Pump{Radio: r}, nil
}

func defaultFreq() uint32 {
	freq := uint32(DefaultFrequency)
	f := os.Getenv(freqEnvVar)
	if len(f) != 0 {
		n, err := fmt.Sscanf(f, "%d", &freq)
		if err != nil {
			log.Fatalf("%s value (%s): %v\n", freqEnvVar, f, err)
		}
		if n != 1 || freq < 860000000 || freq > 920000000 {
			log.Fatalf("%s value (%s) should be the pump frequency in Hz\n", freqEnvVar, f)
		}
	}
	return freq
}
