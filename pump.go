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

var radios [](func() radio.Interface)

func Open() *Pump {
	pump := &Pump{
		timeout: defaultTimeout,
		retries: defaultRetries,
	}
	found := false
	for _, openRadio := range radios {
		pump.Radio = openRadio()
		if pump.Error() == nil {
			found = true
			break
		}
		_, wrongVersion := pump.Error().(radio.HardwareVersionError)
		if !wrongVersion {
			log.Print(pump.Error())
			break
		}
		pump.SetError(nil)
	}
	if !found {
		pump.SetError(fmt.Errorf("no radio hardware detected"))
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

func (pump *Pump) Timeout() time.Duration {
	return pump.timeout
}

func (pump *Pump) SetTimeout(t time.Duration) {
	pump.timeout = t
}

func (pump *Pump) Retries() int {
	return pump.retries
}

func (pump *Pump) SetRetries(n int) {
	pump.retries = n
}

func (pump *Pump) Rssi() int {
	return pump.rssi
}

func (pump *Pump) Error() error {
	err := pump.Radio.Error()
	if err != nil {
		return err
	}
	return pump.err
}

func (pump *Pump) SetError(err error) {
	pump.Radio.SetError(err)
	pump.err = err
}

func (pump *Pump) PrintStats() {
	stats := pump.Radio.Statistics()
	fmt.Printf("\nTX: %6d    RX: %6d\n", stats.Packets.Sent, stats.Packets.Received)
}
