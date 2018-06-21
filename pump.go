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
	pumpEnvVar       = "MEDTRONIC_PUMP_ID"
	freqEnvVar       = "MEDTRONIC_FREQUENCY"
	defaultFrequency = 916600000
	defaultTimeout   = 500 * time.Millisecond
	defaultRetries   = 3
)

// PumpAddress returns the encoded form of the pump ID.
func PumpAddress() []byte {
	id := os.Getenv(pumpEnvVar)
	if len(id) == 0 {
		log.Fatalf("%s environment variable is not set", pumpEnvVar)
	}
	addr, err := DeviceAddress(id)
	if err != nil {
		log.Fatalf("%s: %v", pumpEnvVar, err)
	}
	return addr
}

// DeviceAddress returns the encoded form of a device ID.
func DeviceAddress(id string) ([]byte, error) {
	if len(id) != 6 {
		return nil, fmt.Errorf("device ID %q must be 6 digits", id)
	}
	h, err := strconv.ParseUint(id, 16, 24)
	if err != nil {
		return nil, fmt.Errorf("device ID %q: %v", id, err)
	}
	return []byte{byte(h >> 16), byte(h >> 8), byte(h >> 0)}, nil
}

// Pump represents a Medtronic insulin pump.
type Pump struct {
	Radio radio.Interface

	// 22 for 522/722, 23 for 523/723, etc.
	family Family

	// Implicit parameters for command execution.
	timeout time.Duration
	retries int
	rssi    int
	err     error
}

// Open opens radio communication with a pump.
func Open() *Pump {
	r := radioInterface()
	pump := &Pump{
		Radio:   r,
		timeout: defaultTimeout,
		retries: defaultRetries,
	}
	if pump.Error() != nil {
		log.Printf("cannot connect to %s radio on %s", r.Name(), r.Device())
		return pump
	}
	log.Printf("connected to %s radio on %s", r.Name(), r.Device())
	precomputePackets()
	freq := getFrequency()
	log.Printf("setting frequency to %s", radio.MegaHertz(freq))
	r.Init(freq)
	go pump.closeWhenSignaled()
	return pump
}

// Close closes communication with the pump.
func (pump *Pump) Close() {
	r := pump.Radio
	log.Printf("disconnecting %s radio on %s", r.Name(), r.Device())
	r.Close()
}

// ParseFrequency interprets the given string as a frequency
// and returns its value in Hertz.
func ParseFrequency(s string) (uint32, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	if 860.0 <= f && f <= 920.0 {
		return uint32(f * 1000000.0), nil
	}
	if 860000000.0 <= f && f <= 920000000.0 {
		return uint32(f), nil
	}
	return 0, fmt.Errorf("invalid frequency %s", s)
}

func getFrequency() uint32 {
	s := os.Getenv(freqEnvVar)
	if len(s) == 0 {
		return uint32(defaultFrequency)
	}
	f, err := ParseFrequency(s)
	if err != nil {
		log.Fatalf("%s: %v", freqEnvVar, err)
	}
	return f
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
// For safety, state-changing commands are only tried once.
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
