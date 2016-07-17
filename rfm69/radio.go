package rfm69

import (
	"log"
	"time"
)

const (
	verbose       = false
	maxPacketSize = 110
	fifoSize      = 66

	// The fifoThreshold value should allow a maximum-sized packet to be
	// written in two bursts, but be large enough to avoid fifo underflow.
	fifoThreshold = 20

	// Approximate time for one byte to be transmitted, based on the data rate.
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
	packet := make([]byte, len(data)+1)
	copy(packet, data)
	// Prepare for auto-transmit.
	// (Automode from/to sleep mode is not reliable.)
	r.clearFifo()
	r.setMode(StandbyMode)
	r.hw.WriteRegister(RegAutoModes, EnterConditionFifoNotEmpty|ExitConditionFifoEmpty|IntermediateModeTx)
	r.transmit(packet)
	r.setMode(SleepMode)
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
		r.hw.WriteBurst(RegFifo, data[:avail])
		data = data[avail:]
		if len(data) == 0 {
			break
		}
		// Wait until there is room for at least
		// fifoSize - fifoThreshold bytes in the FIFO.
		// Err on the short side here to avoid TXFIFO underflow.
		time.Sleep(fifoSize / 4 * byteDuration)
		for r.Error() == nil {
			if !r.fifoThresholdExceeded() {
				avail = fifoSize - fifoThreshold
				break
			}
		}
	}
	r.finishTx(avail)
}

func (r *Radio) finishTx(numBytes int) {
	time.Sleep(time.Duration(numBytes) * byteDuration)
	// Wait for automatic return to standby mode when FIFO is empty.
	for r.Error() == nil {
		s := r.mode()
		if s == StandbyMode {
			break
		}
		if verbose {
			log.Printf("waiting for TX to finish in %s state", stateName(s))
		}
	}
}

func (r *Radio) fifoEmpty() bool {
	return r.hw.ReadRegister(RegIrqFlags2)&FifoNotEmpty == 0
}

func (r *Radio) fifoFull() bool {
	return r.hw.ReadRegister(RegIrqFlags2)&FifoFull != 0
}

func (r *Radio) fifoThresholdExceeded() bool {
	return r.hw.ReadRegister(RegIrqFlags2)&FifoLevel != 0
}

func (r *Radio) clearFifo() {
	r.hw.WriteRegister(RegIrqFlags2, FifoOverrun)
}

func (r *Radio) Receive(timeout time.Duration) ([]byte, int) {
	if r.Error() != nil {
		return nil, 0
	}
	r.hw.WriteRegister(RegAutoModes, 0)
	r.setMode(ReceiverMode)
	defer r.setMode(SleepMode)
	if verbose {
		log.Printf("waiting for interrupt in %s state", r.State())
	}
	r.hw.AwaitInterrupt(timeout)
	rssi := r.ReadRSSI()
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
		c := r.hw.ReadRegister(RegFifo)
		if r.Error() != nil {
			break
		}
		if c != 0 {
			r.err = r.receiveBuffer.WriteByte(c)
			continue
		}
		// End of packet.
		r.setMode(StandbyMode)
		size := r.receiveBuffer.Len()
		if size == 0 {
			break
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
	return nil, rssi
}
