package medtronic

import (
	"time"
)

// Radio is a mock implementation of radio.Interface.
type Radio struct {
	freq uint32
	err  error
}

// Init initializes the radio device.
func (r *Radio) Init(freq uint32) {
	r.freq = freq
}

// Reset resets the radio device.
func (r *Radio) Reset() {
}

// Close closes the radio device.
func (r *Radio) Close() {}

// Frequency returns the radio's current frequency, in Hertz.
func (r *Radio) Frequency() uint32 {
	return r.freq
}

// SetFrequency sets the radio to the given frequency, in Hertz.
func (r *Radio) SetFrequency(freq uint32) {
	r.freq = freq
}

// Send transmits the given packet.
func (r *Radio) Send(data []byte) {
}

// Receive listens with the given timeout for an incoming packet.
// It returns the packet and the associated RSSI.
func (r *Radio) Receive(timeout time.Duration) ([]byte, int) {
	return nil, 0
}

// SendAndReceive transmits the given packet,
// then listens with the given timeout for an incoming packet.
// It returns the packet and the associated RSSI.
func (r *Radio) SendAndReceive(data []byte, timeout time.Duration) ([]byte, int) {
	return nil, 0
}

// State returns the radio's current state as a string.
func (r *Radio) State() string {
	return "idle"
}

// Error returns the error state of the radio device.
func (r *Radio) Error() error {
	return r.err
}

// SetError sets the error state of the radio device.
func (r *Radio) SetError(err error) {
	r.err = err
}

// Name returns the radio's name.
func (r *Radio) Name() string {
	return "mock"
}

// Device returns the pathname of the radio's device.
func (r *Radio) Device() string {
	return "/dev/null"
}
