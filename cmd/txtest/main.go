package main

import (
	"flag"
	"log"
	"time"

	"github.com/ecc1/medtronic"
	"github.com/ecc1/medtronic/packet"
)

var (
	count            = flag.Int("n", 0, "send only `count` packets")
	minPacketSize    = flag.Int("min", 1, "minimum packet `size` in bytes")
	maxPacketSize    = flag.Int("max", 70, "maximum packet `size` in bytes")
	interPacketDelay = flag.Duration("delay", time.Second, "inter-packet delay")
)

func main() {
	flag.Parse()
	pump := medtronic.Open()
	defer pump.Close()
	n := *minPacketSize
	i := 0
	for pump.Error() == nil {
		if *count != 0 && i == *count {
			return
		}
		data := make([]byte, n+1) // leave space for CRC
		for i := range data {
			data[i] = byte(i + 1)
		}
		packet := packet.Encode(data)
		log.Printf("data:   % X", data)
		log.Printf("packet: % X", packet)
		pump.Radio.Send(packet)
		i++
		n++
		if n > *maxPacketSize {
			n = *minPacketSize
		}
		time.Sleep(*interPacketDelay)
	}
	log.Fatal(pump.Error())
}
