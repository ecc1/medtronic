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
	buf := []byte{READ_SINGLE | addr, 0xFF}
	err := dev.Transfer(buf)
	if err != nil {
		return 0, err
	}
	return buf[1], nil
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
		time.Sleep(10 * time.Millisecond)
	}
}

func Write(dev *spi.Device, data []byte) error {
	n := len(data)
	if n%2 != 0 {
		panic("odd data length in WriteRegisters()")
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
	Strobe(dev, SRES)
	for {
		status, err := Strobe(dev, SNOP)
		if err != nil {
			return err
		}
		if status&CHIP_RDY == 0 {
			s := (status & STATE_MASK) >> 4
			log.Printf("chip ready, state = %d\n", s)
			break
		}
		log.Printf("chip not yet ready; sleeping\n")
		time.Sleep(10 * time.Millisecond)
	}
	return WriteRegister(dev, SYNC0, 0x55)
}
