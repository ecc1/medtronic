package cc1100

import (
	"fmt"
	"log"
	"time"
)

const (
	usePolling      = false
	pollInterval    = time.Millisecond
	maxPacketLength = 100

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
	err := r.changeState(SRX, STATE_RX)
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
		if usePolling {
			n, _ := r.ReadNumRxBytes()
			if n == 0 {
				time.Sleep(pollInterval)
				continue
			}
		} else {
			r.interruptPin.Wait()
		}
		r.interrupt <- struct{}{}
	}
}

func (r *Radio) transmit(data []byte) error {
	if len(data) > maxPacketLength {
		return fmt.Errorf("packet too long (%d bytes)", len(data))
	}
	// Terminate packet with zero byte,
	// and pad with another to ensure final bytes
	// are transmitted before leaving TX state.
	err := r.WriteFifo(append(data, []byte{0, 0}...))
	if err != nil {
		return err
	}
	err = r.changeState(STX, STATE_TX)
	if err != nil {
		return err
	}
	return r.drainTxFifo(len(data) + 1)
}

func (r *Radio) drainTxFifo(numBytes int) error {
	time.Sleep(time.Duration(numBytes) * byteDuration)
	for {
		n, err := r.ReadNumTxBytes()
		if err != nil && err != TxFifoUnderflow {
			return err
		}
		if n == 0 || err == TxFifoUnderflow {
			r.PacketsSent++
			break
		}
		s, err := r.ReadState()
		if err != nil {
			return err
		}
		if s != STATE_TX && s != STATE_TXFIFO_UNDERFLOW {
			return fmt.Errorf("unexpected %s state during TXFIFO drain", StateName(s))
		}
		if verbose {
			log.Printf("waiting to transmit %d bytes\n", n)
		}
	}
	return r.changeState(SIDLE, STATE_IDLE)
}

func (r *Radio) receive() error {
	err := r.changeState(SRX, STATE_RX)
	if err != nil {
		return err
	}
	waiting := false
	for {
		numBytes, err := r.ReadNumRxBytes()
		if err == RxFifoOverflow {
			r.changeState(SRX, STATE_RX)
			continue
		} else if err != nil {
			return err
		}
		// Don't read last byte of FIFO if packet is still
		// being received. See Section 20 of data sheet.
		if numBytes < 2 {
			if waiting {
				return nil
			}
			waiting = true
			time.Sleep(byteDuration)
			continue
		}
		waiting = false
		c, err := r.ReadRegister(RXFIFO)
		if err != nil {
			return err
		}
		if c != 0 {
			err = r.receiveBuffer.WriteByte(c)
			if err != nil {
				return err
			}
			continue
		}
		// End of packet.
		rssi, err := r.ReadRSSI()
		if err != nil {
			return err
		}
		size := r.receiveBuffer.Len()
		if size != 0 {
			r.PacketsReceived++
			p := make([]byte, size)
			_, err := r.receiveBuffer.Read(p)
			if err != nil {
				return err
			}
			r.receiveBuffer.Reset()
			r.receivedPackets <- Packet{Rssi: rssi, Data: p}
		}
		return nil
	}
}

func (r *Radio) changeState(strobe byte, desired byte) error {
	s, err := r.ReadState()
	if err != nil {
		return err
	}
	if verbose && s != desired {
		log.Printf("change from %s to %s\n", StateName(s), StateName(desired))
	}
	for {
		switch s {
		case desired:
			return nil
		case STATE_RXFIFO_OVERFLOW:
			s, err = r.Strobe(SFRX)
		case STATE_TXFIFO_UNDERFLOW:
			s, err = r.Strobe(SFTX)
		default:
			s, err = r.Strobe(strobe)
		}
		if err != nil {
			return err
		}
		s = (s >> STATE_SHIFT) & STATE_MASK
		if verbose {
			log.Printf("  %s\n", StateName(s))
		}
	}
}
