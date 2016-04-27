package cc1100

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

const (
	writeUsingTransfer = false
	verifyWrite        = true
	retryWrite         = false
	verbose            = false
)

func init() {
	if !verbose {
		log.SetOutput(ioutil.Discard)
	}
}

func (dev *Device) ReadRegister(addr byte) (byte, error) {
	buf := []byte{READ_MODE | addr, 0xFF}
	err := dev.spiDev.Transfer(buf)
	if err != nil {
		return 0, err
	}
	return buf[1], nil
}

func (dev *Device) ReadFifo(n uint8) ([]byte, error) {
	buf := make([]byte, n+1)
	buf[0] = READ_MODE | BURST_MODE | RXFIFO
	err := dev.spiDev.Transfer(buf)
	if err != nil {
		return nil, err
	}
	return buf[1:], nil
}

func (dev *Device) WriteRegister(addr byte, value byte) error {
	for {
		var err error
		if writeUsingTransfer {
			err = dev.spiDev.Transfer([]byte{addr, value})
		} else {
			err = dev.spiDev.Write([]byte{addr, value})
		}
		if err != nil || !verifyWrite {
			return err
		}
		v, err := dev.ReadRegister(addr)
		if err != nil || v == value {
			return err
		}
		msg := fmt.Sprintf("read(%02X) returned %02X instead of %02X", addr, v, value)
		if !retryWrite {
			return fmt.Errorf("%s", msg)
		}
		log.Printf("%s; sleeping\n", msg)
		time.Sleep(time.Millisecond)
	}
}

func (dev *Device) WriteFifo(data []byte) error {
	buf := append([]byte{BURST_MODE | TXFIFO}, data...)
	if writeUsingTransfer {
		return dev.spiDev.Transfer(buf)
	} else {
		return dev.spiDev.Write(buf)
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
	buf := []byte{cmd}
	err := dev.spiDev.Transfer(buf)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (dev *Device) Reset() error {
	err := dev.ChangeState(SRES, STATE_IDLE)
	if err != nil {
		return err
	}
	if verifyWrite {
		err = dev.WriteRegister(SYNC0, 0x55)
	}
	return err
}

func (dev *Device) ReadState() (byte, error) {
	status, err := dev.Strobe(SNOP)
	if err != nil {
		return 0, err
	}
	return (status >> STATE_SHIFT) & STATE_MASK, nil
}

func (dev *Device) ChangeState(strobe byte, desired byte) error {
	cmd := strobe
	for {
		log.Printf("issuing %s command, waiting for %s\n", strobeName(cmd), stateName[desired])
		status, err := dev.Strobe(cmd)
		if err != nil {
			return err
		}
		s := (status >> STATE_SHIFT) & STATE_MASK
		log.Printf("state = %s\n", stateName[s])
		if s == desired {
			return nil
		}
		switch s {
		case STATE_RXFIFO_OVERFLOW:
			dev.Strobe(SIDLE)
			dev.Strobe(SFRX)
			cmd = strobe
		case STATE_TXFIFO_UNDERFLOW:
			dev.Strobe(SIDLE)
			dev.Strobe(SFTX)
			cmd = strobe
		default:
			cmd = SNOP
		}
	}
}
