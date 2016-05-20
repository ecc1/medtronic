package medtronic

import (
	"fmt"
	"log"
	"os"

	"github.com/ecc1/cc1100"
)

const (
	freqEnvVar = "MEDTRONIC_FREQUENCY"
)

type Pump struct {
	Radio *cc1100.Radio

	DecodingErrors int
	CrcErrors      int
}

func Open() (*Pump, error) {
	r, err := cc1100.Open()
	if err != nil {
		return nil, err
	}
	err = r.Init()
	if err != nil {
		return nil, err
	}
	err = r.WriteFrequency(defaultFreq())
	if err != nil {
		return nil, err
	}
	return &Pump{Radio: r}, nil
}

func defaultFreq() uint32 {
	freq := cc1100.DefaultFrequency
	f := os.Getenv(freqEnvVar)
	if len(f) == 0 {
		return freq
	}
	n, err := fmt.Sscanf(f, "%d", &freq)
	if err != nil {
		log.Fatalf("%s value (%s): %v\n", freqEnvVar, f, err)
	}
	if n != 1 || freq < 860000000 || freq > 920000000 {
		log.Fatalf("%s value (%s) should be the pump frequency in Hz\n", freqEnvVar, f)
	}
	return freq
}
