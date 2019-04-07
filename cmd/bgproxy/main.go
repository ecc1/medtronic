package main

import (
	"flag"
	"log"
	"time"

	"github.com/ecc1/medtronic"
	"github.com/ecc1/nightscout"
	"github.com/ecc1/papertrail"
)

const (
	sensorPeriod    = 40 * time.Minute
	sensorInterval0 = 285 * time.Second
	sensorInterval1 = 15 * time.Second
	sensorCycle     = sensorInterval0 + sensorInterval1

	sensorCalibrationFactor = 10
)

var (
	jsonFile = flag.String("f", "glucose.json", "read BGs from JSON `file`")
)

func main() {
	flag.Parse()
	papertrail.StartLogging()
	entries := readBGs()
	getISIGs(entries)
	if sendBGs() {
		sendBGs()
	} else {
		log.Printf("warning: missed initial transmission of this packet")
	}
}

func readBGs() nightscout.Entries {
	var err error
	entries, err := nightscout.ReadEntries(*jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	if len(entries) < 2 {
		log.Fatal("not enough BG entries")
	}
	entries.Sort()
	e := entries[0]
	log.Printf("latest BG: %d at %s", e.SGV, e.Time().Format(medtronic.UserTimeLayout))
	return entries
}

func getISIGs(entries nightscout.Entries) {
	for i := range isigBuffer {
		if i >= len(entries) {
			isigBuffer[i] = 0
			continue
		}
		isigBuffer[i] = entries[i].SGV * sensorCalibrationFactor
	}
}

func sendBGs() bool {
	// Synchronize to sensor transmission schedule.
	d := time.Duration(medtronic.SinceMidnight(time.Now()))
	m := d % sensorPeriod
	cycle := int(m / sensorCycle)
	m %= sensorCycle
	var delta time.Duration
	var retransmit int
	if m < sensorInterval0 {
		delta = sensorInterval0 - m
		retransmit = 0
	} else {
		delta = sensorCycle - m
		retransmit = 1
	}
	log.Printf("sleeping for %v", delta)
	time.Sleep(delta)

	seq := byte(cycle<<4 | retransmit)
	p := sensorPacket(seq)
	log.Printf("sending packet %02X", seq)
	pump := medtronic.Open()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	pump.Radio.Send(p)
	pump.Close()

	return retransmit == 0
}
