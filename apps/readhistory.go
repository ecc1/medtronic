package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ecc1/medtronic"
)

const (
	timeLayout = "2006-01-02 15:04:05"
)

var (
	verbose = flag.Bool("v", false, "print record details")
	model   = flag.Int("m", 523, "pump model")

	timeBlank = strings.Repeat(" ", len(timeLayout))
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
		for _, r := range medtronic.DecodeHistoryRecords(data, newer) {
			printRecord(r)
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
	t := r.Time()
	tStr := timeBlank
	if !t.IsZero() {
		tStr = t.Format(timeLayout)
	}
	if *verbose {
		fmt.Printf("%s %v %+v\n", tStr, r.Type(), r)
	} else {
		fmt.Printf("%s %v\n", tStr, r.Type())
	}
}
