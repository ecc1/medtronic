package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ecc1/medtronic"
	"github.com/ecc1/nightscout"
)

var (
	verbose = flag.Bool("v", false, "print record details")
	model   = flag.Int("m", 523, "pump model")
	nsFlag  = flag.Bool("t", false, "format as Nightscout treatments")

	timeBlank = strings.Repeat(" ", len(medtronic.UserTimeLayout))
)

func main() {
	flag.Parse()
	family := medtronic.Family(*model % 100)
	for _, file := range flag.Args() {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		data, err := readBytes(f)
		_ = f.Close()
		if err != nil {
			log.Fatal(err)
		}
		readHistory(data, family)
	}
}

func readBytes(r io.Reader) ([]byte, error) {
	var data []byte
	for {
		var b byte
		n, err := fmt.Fscanf(r, "%02x", &b)
		if n == 0 {
			break
		}
		if err != nil {
			return data, err
		}
		data = append(data, b)
	}
	return data, nil
}

func readHistory(data []byte, family medtronic.Family) {
	records, err := medtronic.DecodeHistory(data, family)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		*verbose = true
	}
	if *verbose {
		fmt.Println(nightscout.JSON(records))
	} else if *nsFlag {
		medtronic.ReverseHistory(records)
		fmt.Println(nightscout.JSON(medtronic.Treatments(records)))
	} else {
		for _, r := range records {
			printRecord(r)
		}
	}
}

func printRecord(r medtronic.HistoryRecord) {
	t := r.Time
	tStr := timeBlank
	if !t.IsZero() {
		tStr = t.Format(medtronic.UserTimeLayout)
	}
	fmt.Printf("%s %v\n", tStr, r.Type())
}
