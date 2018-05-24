package main

// Regenerate the JSON-encoded history records on stdin
// by extracting their Data fields and decoding them again.

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/ecc1/medtronic"
)

var (
	model = flag.Int("m", 523, "pump model")
)

func main() {
	flag.Parse()
	family := medtronic.Family(*model % 100)
	d := json.NewDecoder(os.Stdin)
	var maps []interface{}
	if err := d.Decode(&maps); err != nil {
		log.Fatal(err)
	}
	var records medtronic.History
	for _, v := range maps {
		m := v.(map[string]interface{})
		base64data := m["Data"].(string)
		var data []byte
		data, err := base64.StdEncoding.DecodeString(base64data)
		if err != nil {
			log.Fatal(err)
		}
		r, err := medtronic.DecodeHistoryRecord(data, family)
		if err != nil {
			log.Fatal(err)
		}
		records = append(records, r)
	}
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	if err := e.Encode(records); err != nil {
		log.Fatal(err)
	}
}
