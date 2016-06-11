package main

import (
	"fmt"
	"log"

	"github.com/ecc1/medtronic"
)

func main() {
	pump := medtronic.Open()
	pump.Model()
	err := pump.Error()
	if err != nil {
		_, noResponse := err.(medtronic.NoResponseError)
		if noResponse {
			fmt.Println(-128)
			return
		}
		log.Fatal(err)
	}
	fmt.Println(pump.Rssi())
}
