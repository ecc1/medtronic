package main

// Fetch recent CGM readings from a Medtronic pump,
// with options to upload to Nightscout and update a local JSON file.

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ecc1/medtronic"
	"github.com/ecc1/nightscout"
	"github.com/ecc1/papertrail"
)

type (
	// Entries is an alias, for conciseness.
	Entries = nightscout.Entries
)

const (
	maxClockDelta = 5 * time.Minute
	gapDuration   = 7 * time.Minute
)

var (
	cgmHistory         = flag.Duration("b", 20*time.Minute, "maximum age of CGM entries to fetch")
	sinceFlag          = flag.String("t", "", "get records since the specified `time` in RFC3339 format")
	uploadFlag         = flag.Bool("u", false, "upload to Nightscout")
	simulateUploadFlag = flag.Bool("s", false, "simulate upload to Nightscout")
	verboseFlag        = flag.Bool("v", false, "verbose mode")
	jsonFile           = flag.String("f", "", "append results to JSON `file`")
	jsonCutoff         = flag.Duration("k", 7*24*time.Hour, "maximum age of CGM entries to keep in JSON file")

	pump       *medtronic.Pump
	cgmTime    time.Time
	cgmEpoch   time.Time
	cgmRecords medtronic.CGMHistory
	oldEntries Entries
	newEntries Entries

	somethingFailed = false
)

func main() {
	flag.Parse()
	if *simulateUploadFlag {
		*uploadFlag = true
	}
	nightscout.SetNoUpload(*simulateUploadFlag)
	nightscout.SetVerbose(*verboseFlag)
	papertrail.StartLogging()
	if *jsonFile != "" {
		oldEntries = readJSON()
	}
	getCGMInfo()
	if *verboseFlag && !*uploadFlag {
		newEntries.Print()
	}
	if *jsonFile != "" {
		updateJSON()
	}
	if *uploadFlag {
		uploadEntries()
	}
	if somethingFailed {
		os.Exit(1)
	}
}

func getCGMInfo() {
	pump = medtronic.Open()
	pump.Wakeup()
	cgmTime = checkCGMClock()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	if *sinceFlag != "" {
		var err error
		cgmEpoch, err = time.Parse(medtronic.JSONTimeLayout, *sinceFlag)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		cgmEpoch = cgmTime.Add(-*cgmHistory)
	}
	// Use time of most recent entry to reduce how far back to go.
	cutoff := cgmEpoch
	if len(oldEntries) != 0 {
		lastTime := oldEntries[0].Time()
		if cutoff.Before(lastTime) {
			cutoff = lastTime
		}
	}
	log.Printf("retrieving records since %s", cutoff.Format(medtronic.UserTimeLayout))
	cgmRecords = pump.CGMHistory(cutoff)
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	log.Printf("%d CGM records", len(cgmRecords))
	newEntries = discardIncomplete(medtronic.NightscoutEntries(cgmRecords))
	describeEntries(newEntries, "Nightscout")
}

func timeStr(e nightscout.Entry) string {
	return e.Time().Format(medtronic.UserTimeLayout)
}

func describeEntries(v Entries, kind string) {
	n := len(v)
	switch n {
	case 0:
		log.Printf("0 %s entries", kind)
	case 1:
		log.Printf("1 %s entry at %s", kind, timeStr(v[0]))
	default:
		log.Printf("%d %s entries from %s to %s", n, kind, timeStr(v[0]), timeStr(v[n-1]))
	}
}

func uploadEntries() {
	gaps, err := nightscout.Gaps(cgmEpoch, gapDuration)
	if err != nil {
		log.Print(err)
		somethingFailed = true
		return
	}
	if *verboseFlag {
		printGaps(gaps)
	}
	if len(gaps) == 0 {
		log.Printf("no Nightscout gaps")
		return
	}
	missing := nightscout.Missing(newEntries, gaps)
	log.Printf("uploading %d entries to Nightscout", len(missing))
	for _, e := range missing {
		err := nightscout.Upload("POST", "entries", e)
		if err != nil {
			log.Print(err)
			somethingFailed = true
			return
		}
	}
}

// If the most recent glucose entry is incomplete, discard it.
// This can happen if the loop runs at the same time the sensor
// transmits a new reading, or if the sensor is warming up. If we
// simply wait until next time, both EGV and sensor records might be
// available.  If the sensor is warming up, there still won't be a
// matching EGV record, but the raw-only entry will be uploaded
// because it's no longer the most recent.
func discardIncomplete(entries Entries) Entries {
	if len(entries) == 0 {
		return entries
	}
	e := entries[0]
	if e.Type == nightscout.SGVType && (e.SGV == 0) {
		return entries[1:]
	}
	return entries
}

func checkCGMClock() time.Time {
	t := pump.Clock()
	if pump.Error() != nil {
		return t
	}
	delta := time.Until(t)
	if delta < 0 {
		delta = -delta
	}
	log.Printf("CGM clock difference = %v", delta)
	if delta > maxClockDelta {
		pump.SetError(fmt.Errorf("CGM clock difference is greater than %v", maxClockDelta))
	}
	return t
}

func printGaps(gaps []nightscout.Gap) {
	for _, g := range gaps {
		t1 := g.Start
		t2 := g.Finish
		gap := t2.Sub(t1)
		s1 := t1.Format(medtronic.UserTimeLayout)
		s2 := t2.Format(medtronic.UserTimeLayout)
		log.Printf("%v gap from %s to %s", gap, s1, s2)
	}
}

func readJSON() Entries {
	entries, err := nightscout.ReadEntries(*jsonFile)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("%s: %v", *jsonFile, err)
		somethingFailed = true
		return nil
	}
	log.Printf("read %d entries from %s", len(entries), *jsonFile)
	entries.Sort()
	return entries
}

func updateJSON() {
	log.Printf("merging %d old and %d new entries", len(oldEntries), len(newEntries))
	merged := nightscout.MergeEntries(oldEntries, newEntries)
	describeEntries(merged, "merged")
	cutoff := cgmTime.Add(-*jsonCutoff)
	trimmed := merged.TrimAfter(cutoff)
	describeEntries(trimmed, "trimmed")
	// Back up JSON file with a "~" suffix.
	err := os.Rename(*jsonFile, *jsonFile+"~")
	if err != nil && !os.IsNotExist(err) {
		log.Print(err)
		somethingFailed = true
	}
	err = trimmed.Save(*jsonFile)
	if err != nil {
		log.Print(err)
		somethingFailed = true
	}
	log.Printf("wrote %d entries to %s", len(trimmed), *jsonFile)
}
