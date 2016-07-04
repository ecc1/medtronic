package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ecc1/medtronic"
)

var (
	verbose = flag.Bool("v", false, "print record details")
	model   = flag.Int("m", 523, "pump model")

	timeBlank = strings.Repeat(" ", len(medtronic.TimeLayout))
)

func main() {
	flag.Parse()
	newer := *model%100 > 22
	for _, file := range flag.Args() {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		data := readBytes(f)
		f.Close()
		records, err := medtronic.DecodeHistoryRecords(data, newer)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		for _, r := range records {
			printRecord(r, *verbose || err != nil)
		}
	}
}

func readBytes(r io.Reader) []byte {
	data := []byte{}
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

func printRecord(r medtronic.HistoryRecord, verbose bool) {
	if verbose {
		b, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			fmt.Printf("%v %v\n", r.Type(), err)
		} else {
			fmt.Printf("%v %s\n", r.Type(), string(b))
		}
	} else {
		t := r.Time
		tStr := timeBlank
		if !t.IsZero() {
			tStr = t.Format(medtronic.TimeLayout)
		}
		fmt.Printf("%s %v\n", tStr, r.Type())
	}
}
