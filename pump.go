package medtronic

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ecc1/radio"
)

const (
	freqEnvVar       = "MEDTRONIC_FREQUENCY"
	defaultFrequency = 916600000
	defaultTimeout   = 500 * time.Millisecond
	defaultRetries   = 3
)

// Pump represents a Medtronic insulin pump.
type Pump struct {
	Radio radio.Interface

	// 22 for 522/722, 23 for 523/723, etc.
	family int

	// Implicit parameters for command execution.
	timeout time.Duration
	retries int
	rssi    int
	err     error
}

// Open opens radio communication with a pump.
func Open() *Pump {
	pump := &Pump{
		Radio:   radioInterface(),
		timeout: defaultTimeout,
		retries: defaultRetries,
	}
	if pump.Error() != nil {
		return pump
	}
	hw := pump.Radio.Hardware()
	log.Printf("connected to %s radio on %s", hw.Name(), hw.Device())
	freq := getFrequency()
	log.Printf("setting frequency to %s", radio.MegaHertz(freq))
	pump.Radio.Init(freq)
	go pump.closeWhenSignaled()
	return pump
}

// Close closes communication with the pump.
func (pump *Pump) Close() {
	hw := pump.Radio.Hardware()
	log.Printf("disconnecting %s radio on %s", hw.Name(), hw.Device())
	pump.Radio.Close()
}

func getFrequency() uint32 {
	s := os.Getenv(freqEnvVar)
	if len(s) == 0 {
		return uint32(defaultFrequency)
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatalf("%s: %v", freqEnvVar, err)
	}
	if 860.0 <= f && f <= 920.0 {
		return uint32(f * 1000000.0)
	}
	if 860000000.0 <= f && f <= 920000000.0 {
		return uint32(f)
	}
	log.Fatalf("%s (%s): invalid pump frequency", freqEnvVar, s)
	panic("unreachable")
}

// Timeout returns the timeout used for pump communications.
func (pump *Pump) Timeout() time.Duration {
	return pump.timeout
}

// SetTimeout sets the timeout used for pump communications.
func (pump *Pump) SetTimeout(t time.Duration) {
	pump.timeout = t
}

// Retries returns the number of retries used for pump communications.
func (pump *Pump) Retries() int {
	return pump.retries
}

// SetRetries sets the number of retries used for pump communications.
func (pump *Pump) SetRetries(n int) {
	pump.retries = n
}

// RSSI returns the RSSI received with the most recent packet from the pump.
func (pump *Pump) RSSI() int {
	return pump.rssi
}

// Error returns the error state of the pump.
func (pump *Pump) Error() error {
	err := pump.Radio.Error()
	if err != nil {
		return err
	}
	return pump.err
}

// SetError sets the error state of the pump.
func (pump *Pump) SetError(err error) {
	pump.Radio.SetError(err)
	pump.err = err
}

// PrintStats prints pump communication statistics.
func (pump *Pump) PrintStats() {
	stats := pump.Radio.Statistics()
	fmt.Printf("\nTX: %6d    RX: %6d\n", stats.Packets.Sent, stats.Packets.Received)
}
