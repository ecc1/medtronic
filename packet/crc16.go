package packet

//go:generate crcgen -size 16 -poly 0x1021

// CRC16 computes the 16-bit CRC of the given data using the CCITT polynomial.
func CRC16(msg []byte) uint16 {
	res := uint16(0xFFFF)
	for _, b := range msg {
		res = res<<8 ^ crc16Table[byte(res>>8)^b]
	}
	return res
}
