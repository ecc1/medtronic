package main

import (
	"log"
	"time"

	"github.com/ecc1/medtronic"
	"github.com/ecc1/medtronic/packet"
)

const (
	verbose = true
)

func main() {
	if verbose {
		log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	}
	pump := medtronic.Open()
	defer pump.Close()
	for pump.Error() == nil {
		p, rssi := pump.Radio.Receive(time.Hour)
		if pump.Error() != nil {
			log.Print(pump.Error())
			pump.SetError(nil)
			continue
		}
		if verbose {
			log.Printf("raw data: % X (%d bytes, RSSI = %d)", p, len(p), rssi)
		}
		data, err := packet.Decode(p)
		if err != nil {
			log.Print(err)
			continue
		}
		if verbose {
			log.Printf("decoded:  % X", data)
		} else {
			log.Printf("% X (%d bytes, RSSI = %d)", data, len(data), rssi)
		}

	}
	log.Fatal(pump.Error())
}
