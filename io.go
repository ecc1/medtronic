package cc1100

import (
	"errors"
	"fmt"
	"log"
	"time"
)

const (
	verbose = false

	writeUsingTransfer = false
	verifyWrite        = false
	retryWrite         = false

	readFifoUsingBurst  = true
	writeFifoUsingBurst = true
)

func (r *Radio) ReadRegister(addr byte) (byte, error) {
	buf := []byte{READ_MODE | addr, 0xFF}
	err := r.device.Transfer(buf)
	if err != nil {
		return 0, err
	}
	return buf[1], nil
}

var (
	RxFifoOverflow  = errors.New("RXFIFO overflow")
	TxFifoUnderflow = errors.New("TXFIFO underflow")
)

// Per section 20 of data sheet, read NUM_RXBYTES
// repeatedly until same value is returned twice.
func (r *Radio) ReadNumRxBytes() (byte, error) {
	last := byte(0)
	read := false
	for {
		n, err := r.ReadRegister(RXBYTES)
		if err != nil {
			return 0, err
		}
		if n&RXFIFO_OVERFLOW != 0 {
			return 0, RxFifoOverflow
		}
		n &= NUM_RXBYTES_MASK
		if read && n == last {
			return n, nil
		}
		last = n
		read = true
	}
}

func (r *Radio) ReadNumTxBytes() (byte, error) {
	n, err := r.ReadRegister(TXBYTES)
	if err != nil {
		return 0, err
	}
	if n&TXFIFO_UNDERFLOW != 0 {
		return 0, TxFifoUnderflow
	}
	return n & NUM_TXBYTES_MASK, nil
}

func (r *Radio) ReadFifo(n uint8) ([]byte, error) {
	if readFifoUsingBurst {
		buf := make([]byte, n+1)
		buf[0] = READ_MODE | BURST_MODE | RXFIFO
		err := r.device.Transfer(buf)
		if err != nil {
			return nil, err
		}
		return buf[1:], nil
	} else {
		buf := make([]byte, n)
		var err error
		for i := uint8(0); i < n; i++ {
			buf[i], err = r.ReadRegister(RXFIFO)
			if err != nil {
				return nil, err
			}
		}
		return buf, nil
	}
}

func (r *Radio) writeData(data []byte) error {
	if writeUsingTransfer {
		return r.device.Transfer(data)
	} else {
		return r.device.Write(data)
	}
}

func (r *Radio) WriteRegister(addr byte, value byte) error {
	for {
		err := r.writeData([]byte{addr, value})
		if err != nil || !verifyWrite || addr == TXFIFO {
			return err
		}
		v, err := r.ReadRegister(addr)
		if err != nil || v == value {
			return err
		}
		msg := fmt.Sprintf("read(%X) returned %X instead of %X", addr, v, value)
		if !retryWrite {
			return fmt.Errorf("%s", msg)
		}
		if verbose {
			log.Printf("%s; sleeping\n", msg)
		}
		time.Sleep(time.Millisecond)
	}
}

func (r *Radio) WriteFifo(data []byte) error {
	if writeFifoUsingBurst {
		return r.writeData(append([]byte{BURST_MODE | TXFIFO}, data...))
	} else {
		for _, b := range data {
			err := r.WriteRegister(TXFIFO, b)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (r *Radio) WriteEach(data []byte) error {
	n := len(data)
	if n%2 != 0 {
		panic("odd data length")
	}
	for i := 0; i < n; i += 2 {
		err := r.WriteRegister(data[i], data[i+1])
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Radio) Strobe(cmd byte) (byte, error) {
	if verbose && cmd != SNOP {
		log.Printf("issuing %s command\n", strobeName(cmd))
	}
	buf := []byte{cmd}
	err := r.device.Transfer(buf)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (r *Radio) Reset() error {
	return r.changeState(SRES, STATE_IDLE)
}

func (r *Radio) ReadState() (byte, error) {
	status, err := r.Strobe(SNOP)
	if err != nil {
		return 0, err
	}
	return (status >> STATE_SHIFT) & STATE_MASK, nil
}

func (r *Radio) ReadMarcState() (byte, error) {
	state, err := r.ReadRegister(MARCSTATE)
	if err != nil {
		return 0, err
	}
	return state & MARCSTATE_MASK, nil
}
