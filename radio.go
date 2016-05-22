package rfm69

import (
	"log"
	"time"
)

const (
	fifoSize = 66

	// Approximate time for one byte to be transmitted, based on
	// the data rate.  It was determined empirically so that few
	// if any iterations are needed in drainTxFifo().
	byteDuration = time.Millisecond
)

func (r *Radio) startRadio() {
	if !r.radioStarted {
		r.radioStarted = true
		go r.radio()
		go r.awaitInterrupts()
	}
}

func (r *Radio) Incoming() <-chan Packet {
	return r.receivedPackets
}

func (r *Radio) Outgoing() chan<- Packet {
	return r.transmittedPackets
}

func (r *Radio) radio() {
	err := r.SetMode(ReceiverMode)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case packet := <-r.transmittedPackets:
			err = r.transmit(packet.Data)
		case <-r.interrupt:
			err = r.receive()
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (r *Radio) awaitInterrupts() {
	for {
		r.interruptPin.Wait()
		log.Printf("Interrupt!\n") //XXX
		r.interrupt <- struct{}{}
	}
}

func (r *Radio) transmit(data []byte) error {
	// Terminate packet with zero byte,
	// and pad with another to ensure final bytes
	// are transmitted before leaving TX state.
	data = append(data, []byte{0, 0}...)
	if len(data) <= fifoSize {
		return r.transmitSmall(data)
	} else {
		return r.transmitLarge(data)
	}
}

func (r *Radio) transmitSmall(data []byte) error {
	// FIXME
	return nil
}

// Transmit a packet that is larger than the TXFIFO size.
// See TI Design Note DN500 (swra109c).
func (r *Radio) transmitLarge(data []byte) error {
	// FIXME
	return nil
}

func (r *Radio) receive() error {
	// FIXME
	return nil
}
