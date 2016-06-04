package medtronic

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ecc1/cc1101"
	"github.com/ecc1/radio"
)

const (
	DefaultFrequency = 916600000
	freqEnvVar       = "MEDTRONIC_FREQUENCY"
	defaultTimeout   = 500 * time.Millisecond
	defaultRetries   = 3
)

type Pump struct {
	Radio radio.Interface

	// Implicit parameters for command execution.
	timeout time.Duration
	retries int
	rssi    int

	DecodingErrors int
	CrcErrors      int
}

func Open() (*Pump, error) {
	r, err := cc1101.Open()
	if err != nil {
		return nil, err
	}
	freq := defaultFreq()
	log.Printf("setting frequency to %d\n", freq)
	err = r.Init(freq)
	if err != nil {
		return nil, err
	}
	return &Pump{
		Radio:   r,
		timeout: defaultTimeout,
		retries: defaultRetries,
	}, nil
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
