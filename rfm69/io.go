package rfm69

import (
	"fmt"
)

const (
	writeUsingTransfer = false
)

func (r *Radio) ReadRegister(addr byte) byte {
	if r.Error() != nil {
		return 0
	}
	buf := []byte{addr, 0}
	r.err = r.device.Transfer(buf)
	return buf[1]
}

func (r *Radio) ReadBurst(addr byte, n int) []byte {
	if r.Error() != nil {
		return nil
	}
	buf := make([]byte, n+1)
	buf[0] = addr
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
	r.writeData([]byte{SpiWriteMode | addr, value})
}

func (r *Radio) WriteBurst(addr byte, data []byte) {
	buf := make([]byte, len(data)+1)
	buf[0] = SpiWriteMode | addr
	copy(buf[1:], data)
	r.writeData(buf)
}

func (r *Radio) WriteEach(data []byte) {
	n := len(data)
	if n%2 != 0 {
		panic(fmt.Sprintf("odd data length (%d)", n))
	}
	for i := 0; i < n; i += 2 {
		r.WriteRegister(data[i], data[i+1])
	}
}

func (r *Radio) Version() uint16 {
	v := r.ReadRegister(RegVersion)
	return uint16(v>>4)<<8 | uint16(v&0xF)
}
