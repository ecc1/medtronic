package cc1101

import (
	"log"
)

const (
	verbose            = false
	writeUsingTransfer = false
)

func init() {
	if verbose {
		log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	}
}

func (r *Radio) ReadRegister(addr byte) byte {
	if r.Error() != nil {
		return 0
	}
	buf := []byte{READ_MODE | addr, 0xFF}
	r.err = r.device.Transfer(buf)
	return buf[1]
}

func (r *Radio) ReadBurst(addr byte, n int) []byte {
	reg := addr & 0x3F
	if 0x30 <= reg && reg <= 0x3D {
		log.Panicf("burst access for status register %X is not available", reg)
	}
	if r.Error() != nil {
		return nil
	}
	buf := make([]byte, n+1)
	buf[0] = READ_MODE | BURST_MODE | addr
	r.err = r.device.Transfer(buf)
	return buf[1:]
}

func (r *Radio) writeData(data []byte) {
	if r.Error() != nil {
		return
	}
	if writeUsingTransfer {
		r.err = r.device.Transfer(data)
	} else {
		r.err = r.device.Write(data)
	}
}

func (r *Radio) WriteRegister(addr byte, value byte) {
	r.writeData([]byte{addr, value})
}

func (r *Radio) WriteBurst(addr byte, data []byte) {
	r.writeData(append([]byte{BURST_MODE | addr}, data...))
}

func (r *Radio) WriteEach(data []byte) {
	n := len(data)
	if n%2 != 0 {
		log.Panicf("odd data length (%d)", n)
	}
	for i := 0; i < n; i += 2 {
		r.WriteRegister(data[i], data[i+1])
	}
}

func (r *Radio) Strobe(cmd byte) byte {
	if verbose && cmd != SNOP {
		log.Printf("issuing %s command", strobeName(cmd))
	}
	buf := []byte{cmd}
	r.err = r.device.Transfer(buf)
	return buf[0]
}

func (r *Radio) Reset() {
	r.Strobe(SRES)
}

func (r *Radio) Version() uint16 {
	p := r.ReadRegister(PARTNUM)
	v := r.ReadRegister(VERSION)
	return uint16(p)<<8 | uint16(v)
}
