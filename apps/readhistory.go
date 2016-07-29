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

	timeBlank = strings.Repeat(" ", len(medtronic.UserTimeLayout))
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
		if *verbose || err != nil {
			b, err := json.MarshalIndent(records, "", "  ")
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(b))
			}
		} else {
			for _, r := range records {
				printRecord(r)
			}
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

func printRecord(r medtronic.HistoryRecord) {
	t := r.Time
	tStr := timeBlank
	if !t.IsZero() {
		tStr = t.Format(medtronic.UserTimeLayout)
	}
	fmt.Printf("%s %v\n", tStr, r.Type())
}
