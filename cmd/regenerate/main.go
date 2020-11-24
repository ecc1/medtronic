package main

// Regenerate the JSON-encoded history records on stdin
// by extracting their Data fields and decoding them again.

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ecc1/medtronic"
)

var (
	model       = flag.Int("m", 523, "pump model")
	dataField   = flag.String("d", "Data", "JSON `field` containing base64-encoded data")
	nsFlag      = flag.Bool("ns", false, "format as Nightscout treatments")
	reverseFlag = flag.Bool("r", false, "reverse the history records")
	binaryFlag  = flag.Bool("b", false, "output base64-encoded data")
)

func main() {
	flag.Parse()
	family := medtronic.Family(*model % 100)
	d := json.NewDecoder(os.Stdin)
	var maps []interface{}
	var err error
	err = d.Decode(&maps)
	if err != nil {
		log.Fatal(err)
	}
	if *reverseFlag {
		Reverse(maps)
	}
	var records medtronic.History
	for _, v := range maps {
		m := v.(map[string]interface{})
		base64data := m[*dataField].(string)
		data, err := base64.StdEncoding.DecodeString(base64data)
		if err != nil {
			log.Fatal(err)
		}
		if *binaryFlag {
			fmt.Printf("% X ", data)
			continue
		}
		r, err := medtronic.DecodeHistoryRecord(data, family)
		if err != nil {
			log.Fatal(err)
		}
		records = append(records, r)
	}
	if *binaryFlag {
		fmt.Println()
		return
	}
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	if *nsFlag {
		err = e.Encode(medtronic.Treatments(records))
	} else {
		err = e.Encode(records)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func Reverse(a []interface{}) {
	for i, j := 0, len(a)-1; i < len(a)/2; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}
