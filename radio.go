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
)

func (dev *Device) StartRadio() {
	if !dev.radioStarted {
		dev.radioStarted = true
		go dev.radio()
		go dev.awaitInterrupts()
	}
}

func (dev *Device) IncomingPackets() <-chan Packet {
	return dev.receivedPackets
}

func (dev *Device) OutgoingPackets() chan<- Packet {
	return dev.transmittedPackets
}

func (dev *Device) radio() {
	err := dev.changeState(SRX, STATE_RX)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case packet := <-dev.transmittedPackets:
			err = dev.transmit(packet.Data)
		case <-dev.interrupt:
			err = dev.receive()
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (dev *Device) awaitInterrupts() {
	for {
		if usePolling {
			n, _ := dev.ReadNumRxBytes()
			if n == 0 {
				time.Sleep(pollInterval)
				continue
			}
		} else {
			dev.interruptPin.Wait()
		}
		dev.interrupt <- struct{}{}
	}
}

const (
	// Approximate time for one byte to be transmitted, based on
	// the data rate.  It was determined empirically so that few
	// if any iterations are needed waiting for the TXFIFO to drain.
	byteDuration = time.Millisecond
)

func (dev *Device) transmit(data []byte) error {
	if len(data) > maxPacketLength {
		return fmt.Errorf("packet too long (%d bytes)", len(data))
	}
	// Terminate packet with zero byte,
	// and pad with another to ensure final bytes
	// are transmitted before leaving TX state.
	err := dev.WriteFifo(append(data, []byte{0, 0}...))
	if err != nil {
		return err
	}
	err = dev.changeState(STX, STATE_TX)
	if err != nil {
		return err
	}
	// Wait for TXFIFO to drain.
	time.Sleep(time.Duration(len(data)+1) * byteDuration)
	for {
		n, err := dev.ReadNumTxBytes()
		if err != nil && err != TxFifoUnderflow {
			return err
		}
		if n == 0 {
			dev.packetsSent++
			return nil
		}
		s, err := dev.ReadState()
		if err != nil {
			return err
		}
		if s != STATE_TX {
			return fmt.Errorf("unexpected %s state during TXFIFO drain", StateName(s))
		}
		if Verbose {
			log.Printf("waiting to transmit %d bytes\n", n)
		}
	}
}

func (dev *Device) receive() error {
	err := dev.changeState(SRX, STATE_RX)
	if err != nil {
		return err
	}
	for {
		numBytes, err := dev.ReadNumRxBytes()
		if err == RxFifoOverflow {
			dev.changeState(SRX, STATE_RX)
			continue
		} else if err != nil {
			return err
		}
		if numBytes == 0 {
			return nil
		}
		c, err := dev.ReadRegister(RXFIFO)
		if err != nil {
			return err
		}
		if c != 0 {
			err = dev.receiveBuffer.WriteByte(c)
			if err != nil {
				return err
			}
			continue
		}
		// End of packet.
		rssi, err := dev.ReadRSSI()
		if err != nil {
			return err
		}
		size := dev.receiveBuffer.Len()
		if size != 0 {
			dev.packetsReceived++
			p := make([]byte, size)
			_, err := dev.receiveBuffer.Read(p)
			if err != nil {
				return err
			}
			dev.receiveBuffer.Reset()
			dev.receivedPackets <- Packet{Rssi: rssi, Data: p}
		}
		return nil
	}
}

func (dev *Device) changeState(strobe byte, desired byte) error {
	s, err := dev.ReadState()
	if err != nil {
		return err
	}
	if Verbose && s != desired {
		log.Printf("change from %s to %s\n", StateName(s), StateName(desired))
	}
	for {
		switch s {
		case desired:
			return nil
		case STATE_RXFIFO_OVERFLOW:
			s, err = dev.Strobe(SFRX)
		case STATE_TXFIFO_UNDERFLOW:
			s, err = dev.Strobe(SFTX)
		default:
			s, err = dev.Strobe(strobe)
		}
		if err != nil {
			return err
		}
		s = (s >> STATE_SHIFT) & STATE_MASK
		if Verbose {
			log.Printf("  %s\n", StateName(s))
		}
	}
}
