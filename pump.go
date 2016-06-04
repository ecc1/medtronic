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
	freqEnvVar       = "MEDTRONIC_FREQUENCY"
	defaultFrequency = 916600000
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
	s := os.Getenv(freqEnvVar)
	if len(s) == 0 {
		return uint32(defaultFrequency)
	}
	MHz := 0.0
	n, err := fmt.Sscanf(s, "%f", &MHz)
	if err == nil && n == 1 && 860.0 <= MHz && MHz <= 920.0 {
		return uint32(MHz * 1000000.0)
	}
	Hz := uint32(0)
	n, err = fmt.Sscanf(s, "%d", &Hz)
	if err == nil && n == 1 && 860000000 <= Hz && Hz <= 920000000 {
		return Hz
	}
	if err != nil {
		log.Fatalf("%s (%s): %v\n", freqEnvVar, s, err)
	}
	log.Fatalf("%s (%s): invalid pump frequency\n", freqEnvVar, s)
	panic("unreachable")
}
