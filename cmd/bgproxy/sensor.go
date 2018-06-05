package main

import (
	"github.com/ecc1/medtronic/packet"
)

const (
	// ISIGsPerPacket is the number of previous ISIG readings transmitted in each sensor packet.
	ISIGsPerPacket = 9
)

// SensorPacket represents a Medtronic CGM sensor packet.
// See github.com/jberian/mmcommander/src/MMsimulator/MMsimulator.c
type SensorPacket struct {
	Address []byte
	Version int
	Adjust  int
	Seq     int
	Repeat  int
	ISIG    []int
	Battery int
}

func marshalSensorPacket(p SensorPacket, modeBits byte) []byte {
	v := make([]byte, 32)
	v[0] = packet.Sensor | modeBits
	v[1] = 0x0F
	copy(v[2:5], p.Address)
	v[5] = byte(p.Version)
	v[6] = 0x17
	v[7] = byte(p.Adjust)
	v[8] = byte(p.Seq<<4) | byte(p.Repeat&0xF)
	isig := p.ISIG[0]
	v[9] = byte(isig >> 8)
	v[10] = byte(isig)
	isig = p.ISIG[1]
	v[11] = byte(isig >> 8)
	v[12] = byte(isig)
	v[13] = 0
	v[14] = 0x67
	v[15] = 0x67
	v[16] = byte(p.Battery)
	for i, j := 2, 17; i < ISIGsPerPacket; i, j = i+1, j+2 {
		isig := p.ISIG[i]
		v[j] = byte(isig >> 8)
		v[j+1] = byte(isig)
	}
	v[31] = 0
	return packet.Encode(v)
}
