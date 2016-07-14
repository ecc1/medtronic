package cc1101

import (
	"bytes"
	"log"
	"time"
)

const (
	verbose            = false
	maxPacketSize      = 110
	fifoSize           = 64
	readFifoUsingBurst = true

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
	// Terminate packet with zero byte,
	// and pad with another to ensure final bytes
	// are transmitted before leaving TX state.
	packet := make([]byte, len(data)+2)
	copy(packet, data)
	defer r.changeState(SIDLE, STATE_IDLE)
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
		r.hw.WriteBurst(TXFIFO, data[:avail])
		r.changeState(STX, STATE_TX)
		data = data[avail:]
		if len(data) == 0 {
			break
		}
		// Transmitting a packet that is larger than the TXFIFO size.
		// See TI Design Note DN500 (swra109c).
		// Err on the short side here to avoid TXFIFO underflow.
		time.Sleep(fifoSize / 4 * byteDuration)
		for r.Error() == nil {
			n := r.ReadNumTxBytes()
			if n < fifoSize {
				avail = fifoSize - int(n)
				if avail > len(data) {
					avail = len(data)
				}
				break
			}
		}
	}
	r.finishTx(avail)
}

func (r *Radio) finishTx(numBytes int) {
	time.Sleep(time.Duration(numBytes) * byteDuration)
	for r.Error() == nil {
		n := r.ReadNumTxBytes()
		if n == 0 || r.Error() == TxFifoUnderflow {
			break
		}
		s := r.ReadState()
		if s != STATE_TX && s != STATE_TXFIFO_UNDERFLOW {
			log.Panicf("unexpected %s state while finishing TX", StateName(s))
		}
		if verbose {
			log.Printf("waiting to transmit %d bytes in %s state", n, StateName(s))
		}
		time.Sleep(byteDuration)
	}
	if verbose {
		log.Printf("TX finished in %s state", r.State())
	}
}

func (r *Radio) Receive(timeout time.Duration) ([]byte, int) {
	if r.Error() != nil {
		return nil, 0
	}
	r.changeState(SRX, STATE_RX)
	defer r.changeState(SIDLE, STATE_IDLE)
	if verbose {
		log.Printf("waiting for interrupt in %s state", r.State())
	}
	r.hw.AwaitInterrupt(timeout)
	rssi := r.ReadRSSI()
	startedWaiting := time.Time{}
	for r.Error() == nil {
		numBytes := r.ReadNumRxBytes()
		if r.Error() == RxFifoOverflow {
			// Flush RX FIFO and change back to RX.
			r.changeState(SRX, STATE_RX)
			continue
		}
		// Don't read last byte of FIFO if packet is still
		// being received. See Section 20 of data sheet.
		if numBytes < 2 {
			if startedWaiting.IsZero() {
				startedWaiting = time.Now()
			} else if time.Since(startedWaiting) >= timeout {
				break
			}
			time.Sleep(byteDuration)
			continue
		}
		if readFifoUsingBurst {
			data := r.hw.ReadBurst(RXFIFO, int(numBytes))
			if r.Error() != nil {
				break
			}
			i := bytes.IndexByte(data, 0)
			if i == -1 {
				// No zero byte found; packet is still incoming.
				// Append all the data and continue to receive.
				_, r.err = r.receiveBuffer.Write(data)
				continue
			}
			// End of packet.
			_, r.err = r.receiveBuffer.Write(data[:i])
		} else {
			c := r.hw.ReadRegister(RXFIFO)
			if r.Error() != nil {
				break
			}
			if c != 0 {
				r.err = r.receiveBuffer.WriteByte(c)
				continue
			}
		}
		// End of packet.
		r.changeState(SIDLE, STATE_IDLE)
		r.Strobe(SFRX)
		size := r.receiveBuffer.Len()
		if size == 0 {
			break
		}
		r.stats.Packets.Received++
		r.stats.Bytes.Received += size
		p := make([]byte, size)
		_, err := r.receiveBuffer.Read(p)
		r.SetError(err)
		if r.Error() != nil {
			break
		}
		r.receiveBuffer.Reset()
		if verbose {
			log.Printf("received %d-byte packet in %s state; %d bytes remaining", size, r.State(), r.ReadNumRxBytes())
		}
		return p, rssi
	}
	return nil, rssi
}
