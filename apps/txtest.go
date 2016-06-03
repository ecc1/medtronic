package main

import (
	"log"
	"time"

	"github.com/ecc1/medtronic"
)

const (
	maxPacketSize    = 70
	interPacketDelay = time.Second
)

func main() {
	pump, err := medtronic.Open()
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, maxPacketSize)
	for i, _ := range data {
		data[i] = byte(i + 1)
	}
	n := 1
	for {
		log.Printf("data:   % X\n", data[:n])
		packet := medtronic.EncodePacket(data[:n])
		log.Printf("packet: % X\n", packet.Data)
		pump.Radio.Outgoing() <- packet
		n++
		if n > maxPacketSize {
			n = 1
		}
		time.Sleep(interPacketDelay)
	}
}
