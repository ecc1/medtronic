package radio

import (
	"fmt"
	"time"
)

type Counters struct {
	Sent     int
	Received int
}

type Statistics struct {
	Bytes   Counters
	Packets Counters
}

type Interface interface {
	Init(frequency uint32)

	Frequency() uint32
	SetFrequency(uint32)

	Send([]byte)
	Receive(time.Duration) ([]byte, int)

	State() string
	Statistics() Statistics

	Error() error
	SetError(error)
}

func MegaHertz(freq uint32) string {
	MHz := freq / 1000000
	kHz := (freq % 1000000) / 1000
	return fmt.Sprintf("%3d.%03d", MHz, kHz)
}
