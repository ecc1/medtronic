package medtronic

import (
	"log"
	"os"
	"strconv"
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
	freq := getFrequency()
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

func getFrequency() uint32 {
	s := os.Getenv(freqEnvVar)
	if len(s) == 0 {
		return uint32(defaultFrequency)
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatalf("%s: %v\n", freqEnvVar, err)
	}
	if 860.0 <= f && f <= 920.0 {
		return uint32(f * 1000000.0)
	}
	if 860000000.0 <= f && f <= 920000000.0 {
		return uint32(f)
	}
	log.Fatalf("%s (%s): invalid pump frequency\n", freqEnvVar, s)
	panic("unreachable")
}
