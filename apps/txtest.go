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
	n := minPacketSize
	for pump.Error() == nil {
		data := make([]byte, n+1) // leave space for CRC
		for i, _ := range data {
			data[i] = byte(i + 1)
		}
		packet := medtronic.EncodePacket(data)
		log.Printf("data:   % X", data)
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
