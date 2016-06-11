package main

import (
	"log"
	"os"
	"path"

	"github.com/ecc1/medtronic"
)

func main() {
	prog := path.Base(os.Args[0])
	var turnOff bool
	switch prog {
	case "suspend":
		turnOff = true
	case "resume":
		turnOff = false
	default:
		log.Fatalf("Program name (%s) must be \"suspend\" or \"resume\"", prog)
	}
	pump := medtronic.Open()
	pump.Suspend(turnOff)
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
}
