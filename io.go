package cc1100

import (
	"errors"
	"fmt"
	"log"
	"time"
)

const (
	Verbose = false

	writeUsingTransfer = false
	verifyWrite        = false
	retryWrite         = false

	readFifoUsingBurst  = true
	writeFifoUsingBurst = true
)

func (dev *Device) ReadRegister(addr byte) (byte, error) {
	buf := []byte{READ_MODE | addr, 0xFF}
	err := dev.spiDev.Transfer(buf)
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
func (dev *Device) ReadNumRxBytes() (byte, error) {
	last := byte(0)
	read := false
	for {
		n, err := dev.ReadRegister(RXBYTES)
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

func (dev *Device) ReadNumTxBytes() (byte, error) {
	n, err := dev.ReadRegister(TXBYTES)
	if err != nil {
		return 0, err
	}
	if n&TXFIFO_UNDERFLOW != 0 {
		return 0, TxFifoUnderflow
	}
	return n & NUM_TXBYTES_MASK, nil
}

func (dev *Device) ReadFifo(n uint8) ([]byte, error) {
	if readFifoUsingBurst {
		buf := make([]byte, n+1)
		buf[0] = READ_MODE | BURST_MODE | RXFIFO
		err := dev.spiDev.Transfer(buf)
		if err != nil {
			return nil, err
		}
		return buf[1:], nil
	} else {
		buf := make([]byte, n)
		var err error
		for i := uint8(0); i < n; i++ {
			buf[i], err = dev.ReadRegister(RXFIFO)
			if err != nil {
				return nil, err
			}
		}
		return buf, nil
	}
}

func (dev *Device) writeData(data []byte) error {
	if writeUsingTransfer {
		return dev.spiDev.Transfer(data)
	} else {
		return dev.spiDev.Write(data)
	}
}

func (dev *Device) WriteRegister(addr byte, value byte) error {
	for {
		err := dev.writeData([]byte{addr, value})
		if err != nil || !verifyWrite || addr == TXFIFO {
			return err
		}
		v, err := dev.ReadRegister(addr)
		if err != nil || v == value {
			return err
		}
		msg := fmt.Sprintf("read(%X) returned %X instead of %X", addr, v, value)
		if !retryWrite {
			return fmt.Errorf("%s", msg)
		}
		if Verbose {
			log.Printf("%s; sleeping\n", msg)
		}
		time.Sleep(time.Millisecond)
	}
}

func (dev *Device) WriteFifo(data []byte) error {
	if writeFifoUsingBurst {
		return dev.writeData(append([]byte{BURST_MODE | TXFIFO}, data...))
	} else {
		for _, b := range data {
			err := dev.WriteRegister(TXFIFO, b)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (dev *Device) WriteEach(data []byte) error {
	n := len(data)
	if n%2 != 0 {
		panic("odd data length")
	}
	for i := 0; i < n; i += 2 {
		err := dev.WriteRegister(data[i], data[i+1])
		if err != nil {
			return err
		}
	}
	return nil
}

func (dev *Device) Strobe(cmd byte) (byte, error) {
	if Verbose && cmd != SNOP {
		log.Printf("issuing %s command\n", strobeName(cmd))
	}
	buf := []byte{cmd}
	err := dev.spiDev.Transfer(buf)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (dev *Device) Reset() error {
	return dev.changeState(SRES, STATE_IDLE)
}

func (dev *Device) ReadState() (byte, error) {
	status, err := dev.Strobe(SNOP)
	if err != nil {
		return 0, err
	}
	return (status >> STATE_SHIFT) & STATE_MASK, nil
}

func (dev *Device) ReadMarcState() (byte, error) {
	state, err := dev.ReadRegister(MARCSTATE)
	if err != nil {
		return 0, err
	}
	return state & MARCSTATE_MASK, nil
}
