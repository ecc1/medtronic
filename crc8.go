package cc1100

//go:generate ../gen_crc_table/gen_crc_table

func Crc8(msg []byte) byte {
	res := byte(0)
	for _, b := range msg {
		res = crc8Table[res^b]
	}
	return res
}
