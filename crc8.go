package medtronic

//go:generate crcgen -size 8 -poly 0x9B

// Compute CRC-8 using WCDMA polynomial.
func Crc8(msg []byte) byte {
	res := byte(0)
	for _, b := range msg {
		res = crc8Table[res^b]
	}
	return res
}
