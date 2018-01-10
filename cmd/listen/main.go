package main

import (
	"flag"
	"log"
	"time"

	"github.com/ecc1/medtronic"
	"github.com/ecc1/medtronic/packet"
)

var (
	listenDuration = flag.Duration("t", time.Hour, "max `duration` to listen")
)

func main() {
	flag.Parse()
	pump := medtronic.Open()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	defer pump.Close()
	log.Printf("listening for %v", *listenDuration)
	p, rssi := pump.Radio.Receive(*listenDuration)
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	if len(p) == 0 {
		log.Fatal("timeout")
	}
	log.Printf("raw data: % X", p)
	data, err := packet.Decode(p)
	if err == nil {
		log.Printf("decoded:  % X", data)
	} else {
		log.Print(err)
	}
	log.Printf("RSSI %d", rssi)
}
