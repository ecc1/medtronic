package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ecc1/medtronic"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s command seq#\n", os.Args[0])
	}
	cmd, err := strconv.ParseUint(os.Args[1], 16, 8)
	if err != nil {
		log.Fatal(err)
	}
	seq, err := strconv.ParseUint(os.Args[2], 16, 8)
	if err != nil {
		log.Fatal(err)
	}
	pump := medtronic.Open()
	defer pump.Close()
	log.Printf("sending RF remote command %02X seq# %02X", cmd, seq)
	pump.SetRetries(50)
	pump.SetTimeout(15 * time.Millisecond)
	pump.RFRemote(medtronic.Command(cmd), uint8(seq))
}
