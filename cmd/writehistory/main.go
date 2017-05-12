package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/ecc1/medtronic"
)

var (
	model = flag.Int("m", 523, "pump model")
	prune = flag.Bool("p", false, "prune results")
)

func main() {
	flag.Parse()
	newer := *model%100 > 22
	var records []medtronic.HistoryRecord
	for _, file := range flag.Args() {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		data := readBytes(f)
		f.Close() // nolint
		v, err := medtronic.DecodeHistoryRecords(data, newer)
		if err != nil {
			log.Fatal(err)
		}
		records = append(records, v...)
	}
	if *prune {
		log.Printf("pruning %d records", len(records))
		records = pruneRecords(records)
	}
	log.Printf("marshaling %d records", len(records))
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	err := e.Encode(records)
	if err != nil {
		log.Fatal(err)
	}
}

func readBytes(r io.Reader) []byte {
	var data []byte
	s := ""
	for {
		n, err := fmt.Fscan(r, &s)
		if n == 0 {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		b, err := strconv.ParseUint(s, 16, 8)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, byte(b))
	}
	return data
}

// Reduce a set of records to one representative of each record type present.
func pruneRecords(records []medtronic.HistoryRecord) []medtronic.HistoryRecord {
	examples := map[medtronic.HistoryRecordType][]medtronic.HistoryRecord{}
	for _, r := range records {
		examples[r.Type()] = append(examples[r.Type()], r)
	}
	var subset []medtronic.HistoryRecord
	for _, x := range examples {
		subset = append(subset, chooseExample(x))
	}
	return subset
}

// Choose the most complex example, as determined by its JSON length.
func chooseExample(records []medtronic.HistoryRecord) medtronic.HistoryRecord {
	best := medtronic.HistoryRecord{}
	bestLen := 0
	for _, r := range records {
		b, err := json.Marshal(r)
		if err != nil {
			log.Fatal(err)
		}
		if len(b) > bestLen {
			best = r
			bestLen = len(b)
		}
	}
	return best
}
