package medtronic

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ecc1/medtronic/cc1101"
	"github.com/ecc1/medtronic/radio"
	"github.com/ecc1/medtronic/rfm69"
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

	DecodingErrors int
	CrcErrors      int
}

var radios = [](func() radio.Interface){cc1101.Open, rfm69.Open}

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
	log.Printf("connected to %s radio", pump.Radio.Hardware().Name())
	freq := getFrequency()
	log.Printf("setting frequency to %s", radio.MegaHertz(freq))
	pump.Radio.Init(freq)
	return pump
}

func (pump *Pump) Close() {
	log.Printf("disconnecting %s radio", pump.Radio.Hardware().Name())
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
