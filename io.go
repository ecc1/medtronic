package rfm69

import (
	"fmt"
)

const (
	writeUsingTransfer = false
	verifyWrite        = true

	readFifoUsingBurst  = true
	writeFifoUsingBurst = true
)

func (r *Radio) ReadRegister(addr byte) (byte, error) {
	buf := []byte{addr, 0}
	err := r.device.Transfer(buf)
	if err != nil {
		return 0, err
	}
	return buf[1], nil
}

func (r *Radio) ReadBurst(addr byte, n int) ([]byte, error) {
	buf := make([]byte, n+1)
	buf[0] = addr
	err := r.device.Transfer(buf)
	return buf[1:], err
}

func (r *Radio) writeData(data []byte) error {
	if writeUsingTransfer {
		return r.device.Transfer(data)
	} else {
		return r.device.Write(data)
	}
}

func (r *Radio) WriteRegister(addr byte, value byte) error {
	err := r.writeData([]byte{SpiWriteMode | addr, value})
	if err != nil || !verifyWrite || addr == RegFifo {
		return err
	}
	v, err := r.ReadRegister(addr)
	if err != nil {
		return err
	}
	if v != value {
		return fmt.Errorf("read(%X) returned %X instead of %X", addr, v, value)
	}
	return nil
}

func (r *Radio) WriteBurst(addr byte, data []byte) error {
	buf := make([]byte, len(data)+1)
	buf[0] = SpiWriteMode | addr
	copy(buf[1:], data)
	return r.writeData(buf)
}

func (r *Radio) WriteEach(data []byte) error {
	n := len(data)
	if n%2 != 0 {
		panic(fmt.Sprintf("odd data length (%d)", n))
	}
	for i := 0; i < n; i += 2 {
		err := r.WriteRegister(data[i], data[i+1])
		if err != nil {
			return err
		}
	}
	return nil
}
