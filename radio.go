package rfm69

import (
	"log"
	"time"

	"github.com/ecc1/radio"
)

const (
	verbose       = false
	fifoSize      = 66
	fifoThreshold = 30

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

func (r *Radio) Incoming() <-chan radio.Packet {
	return r.receivedPackets
}

func (r *Radio) Outgoing() chan<- radio.Packet {
	return r.transmittedPackets
}

func (r *Radio) radio() {
	for {
		//XXX ???
		err := r.WriteEach([]byte{
			RegPacketConfig1, VariableLength,
			RegPayloadLength, 0xFF,
		})
		if err != nil {
			log.Fatal(err)
		}
		err = r.setMode(ReceiverMode)
		if err != nil {
			log.Fatal(err)
		}
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
		r.interrupt <- struct{}{}
	}
}

func (r *Radio) transmit(data []byte) error {
	// Terminate packet with zero byte.
	data = append(data, 0)
	err := r.setMode(StandbyMode)
	if err != nil {
		return err
	}
	if verbose {
		log.Printf("transmitting %d bytes\n", len(data))
	}
	defer r.setMode(StandbyMode)
	//XXX ???
	err = r.WriteEach([]byte{
		RegPacketConfig1, FixedLength,
		RegPayloadLength, uint8(len(data)),
	})
	if err != nil {
		return err
	}
	err = r.send(data)
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
		err := r.WriteBurst(RegFifo, data[:avail])
		if err != nil {
			return err
		}
		err = r.setMode(TransmitterMode)
		if err != nil {
			return err
		}
		data = data[avail:]
		if len(data) == 0 {
			break
		}
		// Once the FifoLevel bit is clear, there will be
		// at least fifoSize - fifoThreshold bytes available.
		for {
			flags, err := r.ReadRegister(RegIrqFlags2)
			if err != nil {
				return err
			}
			if flags&FifoLevel == 0 {
				avail = fifoSize - fifoThreshold
				break
			}
			// Err on the short side here to avoid TXFIFO underflow.
			time.Sleep(fifoThreshold / 2 * byteDuration)
		}

	}
	return r.awaitTxDone(avail)
}

/* XXX ???
func (r *Radio) awaitTxDone(numBytes int) error {
	for {
		empty, err := r.txFifoEmpty()
		if err != nil {
			return err
		}
		if empty {
			return nil
		}
		if verbose {
			log.Printf("waiting for TXFIFO to empty\n")
		}
		time.Sleep(byteDuration)
	}
}
*/

func (r *Radio) awaitTxDone(numBytes int) error {
	for {
		flags, err := r.ReadRegister(RegIrqFlags2)
		if err != nil {
			return err
		}
		if flags&PacketSent != 0 {
			return nil
		}
		if verbose {
			log.Printf("waiting for packet to be sent\n")
		}
		time.Sleep(byteDuration)
	}
}

func (r *Radio) txFifoEmpty() (bool, error) {
	flags, err := r.ReadRegister(RegIrqFlags2)
	if err != nil {
		return false, err
	}
	return flags&FifoNotEmpty == 0, nil
}

func (r *Radio) receive() error {
	if verbose {
		log.Printf("receiving\n")
	}
	// Make sure to enter standby mode when we're done
	// so that we continue to receive SyncAddressMatch interrupts.
	defer r.setMode(StandbyMode)
	for {
		flags, err := r.ReadRegister(RegIrqFlags2)
		if err != nil {
			return err
		}
		if flags&FifoNotEmpty == 0 {
			time.Sleep(byteDuration)
			continue
		}
		c, err := r.ReadRegister(RegFifo)
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
			r.stats.Packets.Received++
			r.stats.Bytes.Received += size
			p := make([]byte, size)
			_, err := r.receiveBuffer.Read(p)
			if err != nil {
				return err
			}
			r.receiveBuffer.Reset()
			r.receivedPackets <- radio.Packet{Rssi: rssi, Data: p}
		}
		return nil
	}
}
