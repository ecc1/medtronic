package cc1100

import (
	"log"
)

const (
	verbose            = true
	writeUsingTransfer = false
)

func (r *Radio) ReadRegister(addr byte) (byte, error) {
	buf := []byte{READ_MODE | addr, 0xFF}
	err := r.device.Transfer(buf)
	if err != nil {
		return 0, err
	}
	return buf[1], nil
}

func (r *Radio) ReadBurst(addr byte, n int) ([]byte, error) {
	reg := addr & 0x3F
	if 0x30 <= reg && reg <= 0x3D {
		log.Panicf("burst access for status register %X is not available\n", reg)
	}
	buf := make([]byte, n+1)
	buf[0] = READ_MODE | BURST_MODE | addr
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
	return r.writeData([]byte{addr, value})
}

func (r *Radio) WriteBurst(addr byte, data []byte) error {
	return r.writeData(append([]byte{BURST_MODE | addr}, data...))
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
	_, err := r.Strobe(SRES)
	return err
}

func (r *Radio) Version() (uint16, error) {
	p, err := r.ReadRegister(PARTNUM)
	if err != nil {
		return 0, err
	}
	v, err := r.ReadRegister(VERSION)
	if err != nil {
		return 0, err
	}
	return uint16(p)<<8 | uint16(v), nil
}
