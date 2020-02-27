package medtronic

func twoByteUint(data []byte) uint16 {
	return uint16(data[0])<<8 | uint16(data[1])
}

func twoByteUintLE(data []byte) uint16 {
	return uint16(data[1])<<8 | uint16(data[0])
}

func twoByteInt(data []byte) int {
	return int(int16(twoByteUint(data)))
}

func twoByteIntLE(data []byte) int {
	return int(int16(twoByteUintLE(data)))
}

func fourByteUint(data []byte) uint32 {
	return uint32(twoByteUint(data[0:2]))<<16 | uint32(twoByteUint(data[2:4]))
}

func fourByteInt(data []byte) int {
	return int(int32(fourByteUint(data)))
}

func marshalUint16(n uint16) []byte {
	return []byte{byte(n >> 8), byte(n & 0xFF)}
}

func marshalUint16LE(n uint16) []byte {
	return []byte{byte(n & 0xFF), byte(n >> 8)}
}

func marshalUint32(n uint32) []byte {
	return append(marshalUint16(uint16(n>>16)), marshalUint16(uint16(n))...)
}
