package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ecc1/medtronic"
	"github.com/ecc1/nightscout"
)

var (
	verbose = flag.Bool("v", false, "print record details")
	nsFlag  = flag.Bool("t", false, "format as Nightscout treatments")

	timeBlank = strings.Repeat(" ", len(medtronic.UserTimeLayout))
)

func main() {
	flag.Parse()
	var history medtronic.CGMHistory
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
		records, _, err := medtronic.DecodeCGMHistory(data, time.Time{})
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		history = append(history, records...)
	}
	if *verbose {
		fmt.Println(nightscout.JSON(history))
	} else if *nsFlag {
		fmt.Println(nightscout.JSON(medtronic.NightscoutEntries(history)))
	} else {
		for _, r := range history {
			printRecord(r)
		}
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

func printRecord(r medtronic.CGMRecord) {
	t := r.Time
	tStr := timeBlank
	if !t.IsZero() {
		tStr = t.Format(medtronic.UserTimeLayout)
	}
	fmt.Printf("%s %v", tStr, r.Type)
	if r.Glucose != 0 {
		fmt.Printf(" %3d", r.Glucose)
	}
	if r.Value != "" {
		fmt.Printf(" %s", r.Value)
	}
	fmt.Println()
}
