package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ecc1/medtronic"
	"github.com/ecc1/nightscout"
)

var (
	noTimes   = flag.Bool("notimes", false, "do not add times to glucose records")
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
		data := readBytes(f)
		_ = f.Close()
		records, err := medtronic.DecodeCGMHistory(data)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		history = append(history, records...)
	}
	if !*noTimes {
		medtronic.AddCGMTimes(history)
	}
	fmt.Println(nightscout.JSON(history))
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

func printRecord(r medtronic.CGMRecord) {
	t := time.Time(r.Time)
	tStr := timeBlank
	if !t.IsZero() {
		tStr = t.Format(medtronic.UserTimeLayout)
	}
	fmt.Printf("%s %v", tStr, r.Type)
	if r.Glucose != 0 {
		fmt.Printf(" %3d", r.Glucose)
	}
	fmt.Println()
}
