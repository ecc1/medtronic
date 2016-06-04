package main

import (
	"log"
	"os"
	"path"

	"github.com/ecc1/medtronic"
)

func main() {
	pump, err := medtronic.Open()
	if err != nil {
		log.Fatal(err)
	}
	prog := path.Base(os.Args[0])
	switch prog {
	case "suspend":
		log.Printf("suspending pump\n")
		err = pump.Suspend(true)
	case "resume":
		log.Printf("resuming pump\n")
		err = pump.Suspend(false)
	default:
		log.Fatalf("Program name (%s) must be \"suspend\" or \"resume\"\n", prog)
	}
	if err != nil {
		log.Fatal(err)
	}
}
