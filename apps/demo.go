package main

import (
	"fmt"
	"log"

	"github.com/ecc1/cc1100"
)

func main() {
	dev, err := cc1100.Open()
	if err != nil {
		log.Fatal(err)
	}

	dev.Reset()
	fmt.Printf("\nDefault RF settings:\n")
	dev.DumpRF()

	dev.InitRF()
	fmt.Printf("\nRF settings after initialization:\n")
	dev.DumpRF()
}
