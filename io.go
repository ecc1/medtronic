package cc1100

import (
	"fmt"
	"log"
	"time"

	"github.com/ecc1/spi"
)

const (
	writeUsingTransfer = false
	verifyWrite        = true
	retryWrite         = false
)

func ReadRegister(dev *spi.Device, addr byte) (byte, error) {
	buf := []byte{READ_MODE | addr, 0xFF}
	err := dev.Transfer(buf)
	if err != nil {
		return 0, err
	}
	return buf[1], nil
}

func ReadFifo(dev *spi.Device, n uint8) ([]byte, error) {
	buf := make([]byte, n+1)
	buf[0] = READ_MODE | BURST_MODE | RXFIFO
	err := dev.Transfer(buf)
	if err != nil {
		return nil, err
	}
	return buf[1:], nil
}

func WriteRegister(dev *spi.Device, addr byte, value byte) error {
	for {
		var err error
		if writeUsingTransfer {
			err = dev.Transfer([]byte{addr, value})
		} else {
			err = dev.Write([]byte{addr, value})
		}
		if err != nil || !verifyWrite {
			return err
		}
		v, err := ReadRegister(dev, addr)
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

func WriteFifo(dev *spi.Device, data []byte) error {
	buf := append([]byte{BURST_MODE | TXFIFO}, data...)
	if writeUsingTransfer {
		return dev.Transfer(buf)
	} else {
		return dev.Write(buf)
	}
}

func WriteEach(dev *spi.Device, data []byte) error {
	n := len(data)
	if n%2 != 0 {
		panic("odd data length")
	}
	for i := 0; i < n; i += 2 {
		err := WriteRegister(dev, data[i], data[i+1])
		if err != nil {
			return err
		}
	}
	return nil
}

func Strobe(dev *spi.Device, cmd byte) (byte, error) {
	buf := []byte{cmd}
	err := dev.Transfer(buf)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func Reset(dev *spi.Device) error {
	err := ChangeState(dev, SRES, STATE_IDLE)
	if err != nil {
		return err
	}
	if verifyWrite {
		err = WriteRegister(dev, SYNC0, 0x55)
	}
	return err
}

func ReadState(dev *spi.Device) (byte, error) {
	status, err := Strobe(dev, SNOP)
	if err != nil {
		return 0, err
	}
	return (status >> STATE_SHIFT) & STATE_MASK, nil
}

func ChangeState(dev *spi.Device, strobe byte, desired byte) error {
	log.Printf("issuing %s command, waiting for %s\n", strobeName(strobe), stateName[desired])
	status, err := Strobe(dev, strobe)
	if err != nil {
		return err
	}
	for {
		s := (status >> STATE_SHIFT) & STATE_MASK
		log.Printf("state = %s\n", stateName[s])
		if s == desired {
			return nil
		}
		status, err = Strobe(dev, SNOP)
		if err != nil {
			return err
		}
	}
}
