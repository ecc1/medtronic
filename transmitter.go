package cc1100

import (
	"log"
	"time"
)

var (
	PacketsSent int
)

func (dev *Device) TransmitMode() error {
	return dev.ChangeState(STX, STATE_TX)
}

func (dev *Device) TransmitPacket(data []byte) error {
	// Terminate packet with zero byte,
	// and pad with another to ensure final bytes
	// are transmitted before going to IDLE state.
	err := dev.WriteFifo(append(data, []byte{0, 0}...))
	if err != nil {
		return err
	}
	err = dev.TransmitMode()
	if err != nil {
		return err
	}
	// Wait for FIFO to drain.
	// The duration is proportional to the length of the packet
	// and depends on the data rate.  It was determined
	// empirically so that there are few if any iterations in the
	// subsequent loop.
	const byteDuration = 900 * time.Microsecond
	time.Sleep(time.Duration(len(data)+1) * byteDuration)
	for {
		n, err := dev.ReadNumTxBytes()
		if err != nil && err != TxFifoUnderflow {
			return err
		}
		if n == 0 {
			break
		}
		s, err := dev.ReadState()
		if err != nil {
			return err
		}
		if s != STATE_TX && s != STATE_TXFIFO_UNDERFLOW {
			panic(StateName(s)) // FIXME
		}
		if Verbose {
			log.Printf("waiting to transmit %d bytes\n", n)
		}
	}
	err = dev.ChangeState(SIDLE, STATE_IDLE)
	if err != nil {
		return err
	}
	PacketsSent++
	return nil
}
