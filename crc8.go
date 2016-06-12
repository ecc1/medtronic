package medtronic

//go:generate ../crcgen/crcgen

func Crc8(msg []byte) byte {
	res := byte(0)
	for _, b := range msg {
		res = crc8Table[res^b]
	}
	return res
}
