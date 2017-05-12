package packet

//go:generate crcgen -size 8 -poly 0x9B

// CRC8 computes the 8-bit CRC of the given data using the WCDMA polynomial.
func CRC8(msg []byte) byte {
	res := byte(0)
	for _, b := range msg {
		res = crc8Table[res^b]
	}
	return res
}
