package medtronic

//go:generate crcgen -size 16 -poly 0x1021

// Compute CRC-16 using CCITT polynomial.
func Crc16(msg []byte) uint16 {
	res := uint16(0xFFFF)
	for _, b := range msg {
		res = res<<8 ^ crc16Table[byte(res>>8)^b]
	}
	return res
}
