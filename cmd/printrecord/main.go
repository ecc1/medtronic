package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/ecc1/medtronic"
)

var (
	model = flag.Int("m", 523, "pump model")
)

func main() {
	flag.Parse()
	family := medtronic.Family(*model % 100)
	for _, arg := range flag.Args() {
		data, err := base64.StdEncoding.DecodeString(arg)
		if err != nil {
			fmt.Printf("base64 decoding error: %v\n", err)
			continue
		}
		fmt.Printf("[ % X ]\n", data)
		r, err := medtronic.DecodeHistoryRecord(data, family)
		if err != nil {
			fmt.Printf("decoding error: %v\n", err)
			continue
		}
		b, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			fmt.Printf("marshaling error: %v\n", err)
			continue
		}
		fmt.Println(string(b))

	}
}
