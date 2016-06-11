package main

import (
	"log"
	"time"

	"github.com/ecc1/medtronic"
)

const (
	minPacketSize    = 1
	maxPacketSize    = 70
	interPacketDelay = time.Second
)

func main() {
	pump := medtronic.Open()
	data := make([]byte, maxPacketSize)
	for i, _ := range data {
		data[i] = byte(i + 1)
	}
	n := minPacketSize
	for pump.Error() == nil {
		log.Printf("data:   % X", data[:n])
		packet := medtronic.EncodePacket(data[:n])
		log.Printf("packet: % X", packet)
		pump.Radio.Send(packet)
		n++
		if n > maxPacketSize {
			n = minPacketSize
		}
		time.Sleep(interPacketDelay)
	}
	log.Fatal(pump.Error())
}
