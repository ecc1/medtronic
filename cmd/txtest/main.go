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
	pkts := 0
	data := make([]byte, *maxPacketSize)
	for pump.Error() == nil {
		if *count != 0 && pkts == *count {
			return
		}
		for i := 0; i < n; i++ {
			data[i] = byte(i + 1)
		}
		packet := packet.Encode(data[:n])
		log.Printf("data:   % X", data[:n])
		log.Printf("packet: % X", packet)
		pump.Radio.Send(packet)
		pkts++
		n++
		if n > *maxPacketSize {
			n = *minPacketSize
		}
		time.Sleep(*interPacketDelay)
	}
	log.Fatal(pump.Error())
}
