package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ecc1/medtronic"
)

func main() {
	pump, err := medtronic.Open()
	if err != nil {
		log.Fatal(err)
	}
	t, err := pump.Clock()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(t.Format(time.UnixDate))
}
