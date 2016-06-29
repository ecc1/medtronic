package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ecc1/medtronic"
)

func usage() {
	log.Fatalf("Usage: %s command-code", os.Args[0])
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	cmd, err := strconv.ParseUint(os.Args[1], 0, 8)
	if err != nil {
		log.Fatal(err)
	}
	pump := medtronic.Open()
	log.Printf("issuing command 0x%02X", cmd)
	result := pump.Query(medtronic.Command(cmd))
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	fmt.Printf("% X\n", result)
}
