package cc1101

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/ecc1/radio"
)

const (
	readFifoUsingBurst = true
	fifoSize           = 64
	maxPacketSize      = 110

	// Approximate time for one byte to be transmitted, based on
	// the data rate.  It was determined empirically so that few
	// if any iterations are needed in drainTxFifo().
	byteDuration = time.Millisecond
)

func (r *Radio) Start() {
	if !r.radioStarted {
		r.radioStarted = true
		go r.radio()
	}
}

func (r *Radio) Stop() {
	// stop radio goroutines and enter IDLE state
	panic("not implemented")
}

func (r *Radio) Incoming() <-chan radio.Packet {
	return r.receivedPackets
}

func (r *Radio) Outgoing() chan<- radio.Packet {
	return r.transmittedPackets
}

func (r *Radio) radio() {
	go r.awaitInterrupts()
	for {
		err := r.changeState(SRX, STATE_RX)
		if err != nil {
			log.Fatal(err)
		}
		select {
		case packet := <-r.transmittedPackets:
			err := r.changeState(SIDLE, STATE_IDLE)
			if err != nil {
				log.Fatal(err)
			}
			err = r.transmit(packet.Data)
			if err != nil {
				log.Fatal(err)
			}
			err = r.changeState(SIDLE, STATE_IDLE)
			if err != nil {
				log.Fatal(err)
			}
		case <-r.interrupt:
			err = r.receive()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (r *Radio) awaitInterrupts() {
	for {
		if verbose {
			log.Printf("waiting for interrupt in %s state", r.State())
		}
		r.interruptPin.Wait()
		r.interrupt <- struct{}{}
	}
}

func (r *Radio) transmit(data []byte) error {
	if len(data) > maxPacketSize {
		log.Panicf("attempting to send %d-byte packet", len(data))
	}
	if verbose {
		log.Printf("sending %d-byte packet in %s state", len(data), r.State())
	}
	// Terminate packet with zero byte,
	// and pad with another to ensure final bytes
	// are transmitted before leaving TX state.
	packet := make([]byte, len(data), len(data)+2)
	copy(packet, data)
	packet = packet[:cap(packet)]
	err := r.send(packet)
	if err == nil {
		r.stats.Packets.Sent++
		r.stats.Bytes.Sent += len(data)
	}
	return err
}

func (r *Radio) send(data []byte) error {
	avail := fifoSize
	for {
		if avail > len(data) {
			avail = len(data)
		}
		err := r.WriteBurst(TXFIFO, data[:avail])
		if err != nil {
			return err
		}
		err = r.changeState(STX, STATE_TX)
		if err != nil {
			return err
		}
		data = data[avail:]
		if len(data) == 0 {
			break
		}
		// Transmitting a packet that is larger than the TXFIFO size.
		// See TI Design Note DN500 (swra109c).
		// Err on the short side here to avoid TXFIFO underflow.
		time.Sleep(fifoSize / 4 * byteDuration)
		for {
			n, err := r.ReadNumTxBytes()
			if err != nil {
				return err
			}
			if n < fifoSize {
				avail = fifoSize - int(n)
				if avail > len(data) {
					avail = len(data)
				}
				break
			}
		}
	}
	return r.drainTxFifo(avail)
}

func (r *Radio) drainTxFifo(numBytes int) error {
	time.Sleep(time.Duration(numBytes) * byteDuration)
	for {
		n, err := r.ReadNumTxBytes()
		if err != nil && err != TxFifoUnderflow {
			return err
		}
		if n == 0 || err == TxFifoUnderflow {
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
			log.Printf("waiting to transmit %d bytes in %s state", n, r.State())
		}
		time.Sleep(byteDuration)
	}
	if verbose {
		log.Printf("TX FIFO drained in %s state", r.State())
	}
	return nil
}

func (r *Radio) receive() error {
	waiting := false
	for {
		numBytes, err := r.ReadNumRxBytes()
		if err == RxFifoOverflow {
			// Flush RX FIFO and change back to RX.
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
		if readFifoUsingBurst {
			data, err := r.ReadBurst(RXFIFO, int(numBytes))
			if err != nil {
				return err
			}
			i := bytes.IndexByte(data, 0)
			if i == -1 {
				// No zero byte found; packet is still incoming.
				// Append all the data and continue to receive.
				_, err = r.receiveBuffer.Write(data)
				if err != nil {
					return err
				}
				continue
			}
			// End of packet.
			_, err = r.receiveBuffer.Write(data[:i])
			if err != nil {
				return err
			}
		} else {
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
		}
		// End of packet.
		rssi, err := r.ReadRSSI()
		if err != nil {
			return err
		}
		size := r.receiveBuffer.Len()
		if size != 0 {
			r.stats.Packets.Received++
			r.stats.Bytes.Received += size
			p := make([]byte, size)
			_, err := r.receiveBuffer.Read(p)
			if err != nil {
				return err
			}
			r.receiveBuffer.Reset()
			r.receivedPackets <- radio.Packet{Rssi: rssi, Data: p}
			if verbose {
				n, _ := r.ReadNumRxBytes()
				log.Printf("received %d-byte packet in %s state; %d bytes remaining", size, r.State(), n)
			}
		}
		return nil
	}
}

func (r *Radio) drainRxFifo() error {
	n, err := r.ReadNumRxBytes()
	if err == RxFifoOverflow {
		// Flush RX FIFO and change back to RX.
		return r.changeState(SRX, STATE_RX)
	}
	if err != nil || n == 0 {
		return err
	}
	s, err := r.ReadState()
	if err != nil {
		return err
	}
	switch s {
	case STATE_RX:
		log.Printf("draining %d bytes from RXFIFO", n)
		_, err = r.ReadBurst(RXFIFO, int(n))
		return err
	default:
		return fmt.Errorf("unexpected %s state during RXFIFO drain", StateName(s))
	}
}

func (r *Radio) changeState(strobe byte, desired byte) error {
	s, err := r.ReadState()
	if err != nil {
		return err
	}
	if verbose && s != desired {
		log.Printf("change from %s to %s", StateName(s), StateName(desired))
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
			log.Printf("  %s", StateName(s))
		}
	}
}
