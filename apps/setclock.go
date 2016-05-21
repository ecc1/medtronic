package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ecc1/medtronic"
)

const layout = "2006-01-02 15:04:05"

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s YYYY-MM-DD HH:MM:SS\n", os.Args[0])
	}
	date := os.Args[1] + " " + os.Args[2]
	t, err := time.ParseInLocation(layout, date, time.Local)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Setting pump clock to %s\n", t.Format(time.UnixDate))
	pump, err := medtronic.Open()
	if err != nil {
		log.Fatal(err)
	}
	err = pump.SetClock(t, 3)
	if err != nil {
		log.Fatal(err)
	}
}
