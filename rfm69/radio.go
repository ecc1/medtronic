package rfm69

import (
	"log"
	"time"
)

const (
	verbose       = false
	maxPacketSize = 110
	fifoSize      = 66

	// The fifoThreshold value allows a maximum-sized packet to be
	// written in two bursts, but large enough to avoid fifo underflow.
	fifoThreshold = 20

	// Approximate time for one byte to be transmitted, based on
	// the data rate.
	byteDuration = time.Millisecond
)

func init() {
	if verbose {
		log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	}
}

func (r *Radio) Send(data []byte) {
	if len(data) > maxPacketSize {
		log.Panicf("attempting to send %d-byte packet", len(data))
	}
	if r.Error() != nil {
		return
	}
	if verbose {
		log.Printf("sending %d-byte packet in %s state", len(data), r.State())
	}
	// Terminate packet with zero byte.
	packet := make([]byte, len(data), len(data)+1)
	copy(packet, data)
	packet = packet[:cap(packet)]
	r.setMode(TransmitterMode)
	defer r.setMode(StandbyMode)
	r.transmit(packet)
	if r.Error() == nil {
		r.stats.Packets.Sent++
		r.stats.Bytes.Sent += len(data)
	}
}

func (r *Radio) transmit(data []byte) {
	avail := fifoSize
	for r.Error() == nil {
		if avail > len(data) {
			avail = len(data)
		}
		if verbose {
			log.Printf("writing %d bytes to TX FIFO\n", avail)
		}
		r.WriteBurst(RegFifo, data[:avail])
		data = data[avail:]
		if len(data) == 0 {
			break
		}
		// Wait until there is room for at least
		// fifoSize - fifoThreshold bytes in the FIFO.
		for r.Error() == nil {
			if !r.fifoThresholdExceeded() {
				avail = fifoSize - fifoThreshold
				break
			}
			if verbose {
				log.Printf("waiting for TX FIFO space in %s state\n", r.State())
			}
			// Err on the short side here to avoid FIFO underflow.
			time.Sleep(byteDuration)
		}
	}
	r.finishTx()
}

func (r *Radio) finishTx() {
	for r.Error() == nil {
		if r.fifoEmpty() {
			break
		}
		if r.mode() != TransmitterMode {
			panic(r.State() + " instead of Transmitter state")
		}
		if verbose {
			log.Printf("waiting for TX to finish in %s state", r.State())
		}
		time.Sleep(byteDuration)
	}
	if verbose {
		log.Printf("TX finished in %s state", r.State())
	}
}

func (r *Radio) fifoEmpty() bool {
	return r.ReadRegister(RegIrqFlags2)&FifoNotEmpty == 0
}

func (r *Radio) fifoFull() bool {
	return r.ReadRegister(RegIrqFlags2)&FifoFull != 0
}

func (r *Radio) fifoThresholdExceeded() bool {
	return r.ReadRegister(RegIrqFlags2)&FifoLevel != 0
}

func (r *Radio) Receive(timeout time.Duration) ([]byte, int) {
	if r.Error() != nil {
		return nil, 0
	}
	r.setMode(ReceiverMode)
	defer r.setMode(StandbyMode)
	if verbose {
		log.Printf("waiting for interrupt in %s state", r.State())
	}
	r.err = r.interruptPin.Wait(timeout)
	startedWaiting := time.Time{}
	for r.Error() == nil {
		if r.fifoEmpty() {
			if startedWaiting.IsZero() {
				startedWaiting = time.Now()
			} else if time.Since(startedWaiting) >= timeout {
				break
			}
			time.Sleep(byteDuration)
			continue
		}
		c := r.ReadRegister(RegFifo)
		if r.Error() != nil {
			break
		}
		if c != 0 {
			r.err = r.receiveBuffer.WriteByte(c)
			continue
		}
		// End of packet.
		rssi := r.ReadRSSI()
		size := r.receiveBuffer.Len()
		if size == 0 {
			if verbose {
				log.Printf("ignoring empty packet in %s state", r.State())
			}
			continue
		}
		r.stats.Packets.Received++
		r.stats.Bytes.Received += size
		p := make([]byte, size)
		_, r.err = r.receiveBuffer.Read(p)
		if r.Error() != nil {
			break
		}
		r.receiveBuffer.Reset()
		if verbose {
			log.Printf("received %d-byte packet in %s state", size, r.State())
		}
		return p, rssi
	}
	return nil, 0
}
