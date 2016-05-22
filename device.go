package rfm69

import (
	"bytes"
	"time"

	"github.com/ecc1/gpio"
	"github.com/ecc1/spi"
)

const (
	spiSpeed     = 10000000 // Hz
	interruptPin = 14       // Intel Edison GPIO connected to DIO0
	resetPin     = 12       // Intel Edison GPIO connected to RESET
)

type Packet struct {
	Rssi int
	Data []byte
}

type Radio struct {
	device       *spi.Device
	interruptPin gpio.InputPin
	resetPin     gpio.OutputPin

	radioStarted       bool
	receiveBuffer      bytes.Buffer
	transmittedPackets chan Packet
	receivedPackets    chan Packet
	interrupt          chan struct{}

	PacketsSent     int
	PacketsReceived int
}

func Open() (*Radio, error) {
	dev, err := spi.Open(spiSpeed)
	if err != nil {
		return nil, err
	}
	err = dev.SetMaxSpeed(spiSpeed)
	if err != nil {
		return nil, err
	}
	intr, err := gpio.Input(interruptPin, "both", false)
	if err != nil {
		return nil, err
	}
	reset, err := gpio.Output(resetPin, false)
	if err != nil {
		return nil, err
	}
	return &Radio{
		device:             dev,
		interruptPin:       intr,
		resetPin:           reset,
		transmittedPackets: make(chan Packet, 100),
		receivedPackets:    make(chan Packet, 10),
		interrupt:          make(chan struct{}, 10),
	}, nil
}

// Reset module.  See section 7.2.2 of data sheet.
func (r *Radio) Reset() error {
	err := r.resetPin.Write(true)
	if err != nil {
		r.resetPin.Write(false)
		return err
	}
	time.Sleep(100 * time.Microsecond)
	err = r.resetPin.Write(false)
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Millisecond)
	return nil
}

func (r *Radio) Init() error {
	err := r.Reset()
	if err != nil {
		return err
	}
	err = r.InitRF()
	if err != nil {
		return err
	}
	r.startRadio()
	return nil
}
