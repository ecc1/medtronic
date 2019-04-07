package main

import (
	"github.com/ecc1/medtronic/packet"
)

// ISIGsPerPacket is the number of previous ISIG readings transmitted in each sensor packet.
const ISIGsPerPacket = 9

var isigBuffer = make([]int, ISIGsPerPacket)

// See github.com/jberian/mmcommander/src/MMsimulator/MMsimulator.c
// and github.com/sarunia/Symulator-nadajnika-Enlite/PROGRAM/tasks.c
func sensorPacket(seq byte) []byte {
	v := make([]byte, 32)
	v[0] = packet.Sensor | 0x3 // 0xAB
	v[1] = 0x0F
	// v[2:5] = sensor ID 000000
	v[5] = 13 // Firmware version 1.3
	v[6] = 0x0E
	v[7] = 0x1E
	v[8] = seq
	s := isigBuffer[0]
	v[9] = byte(s >> 8)
	v[10] = byte(s)
	s = isigBuffer[1]
	v[11] = byte(s >> 8)
	v[12] = byte(s)
	v[13] = 0
	v[14] = 0x5C
	v[15] = 0x5C
	v[16] = 200 // Battery level
	for i, j := 2, 17; i < ISIGsPerPacket; i, j = i+1, j+2 {
		s := isigBuffer[i]
		v[j] = byte(s >> 8)
		v[j+1] = byte(s)
	}
	v[31] = 0
	return packet.Encode(v)
}
