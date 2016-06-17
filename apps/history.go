package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ecc1/medtronic"
)

func main() {
	pump := medtronic.Open()
	switch len(os.Args) {
	case 1:
		result := pump.CurrentPage()
		if pump.Error() != nil {
			log.Fatal(pump.Error())
		}
		fmt.Println(result)
	case 2:
		page, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		data := pump.History(page)
		if pump.Error() != nil {
			log.Fatal(pump.Error())
		}
		for i, v := range data {
			fmt.Printf("[%02d] % X\n", i, v)
		}
	default:
		log.Fatalf("Usage: %s [page#]", os.Args[0])
	}
}
