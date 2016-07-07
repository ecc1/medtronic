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
	model = flag.Int("m", 523, "pump model")
)

func main() {
	flag.Parse()
	newer := *model%100 > 22
	d := json.NewDecoder(os.Stdin)
	maps := []interface{}{}
	err := d.Decode(&maps)
	if err != nil {
		log.Fatal(err)
	}
	records := []medtronic.HistoryRecord{}
	for _, j := range maps {
		m := j.(map[string]interface{})
		base64data := m["Data"].(string)
		data, err := base64.StdEncoding.DecodeString(base64data)
		if err != nil {
			log.Fatal(err)
		}
		r, err := medtronic.DecodeHistoryRecord(data, newer)
		if err != nil {
			log.Fatal(err)
		}
		records = append(records, r)
	}
	b, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
