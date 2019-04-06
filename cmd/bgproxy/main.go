package main

import (
	"flag"
	"log"
	"sort"
	"time"

	"github.com/ecc1/medtronic"
	"github.com/ecc1/nightscout"
	"github.com/ecc1/papertrail"
)

const (
	sensorCalibration = 10

	sensorPeriod = 40 * time.Minute
)

var (
	sensorID = flag.String("s", "000000", "sensor `ID`")
	//sensorID      = flag.String("s", "257EB6", "sensor `ID`")
	sensorAddress []byte

	jsonFile = flag.String("f", "glucose.json", "read BGs from JSON `file`")

	nsEntries nightscout.Entries

	intervalSec = []int{
		285, 15, 285, 15, 285, 15, 285, 15, 285, 15, 285, 15, 285, 15, 285, 15,
		//279, 17, 262, 15, 347, 18, 253, 13, 283, 15, 306, 15, 227, 18, 316, 16,
	}

	transmitTime = make([]time.Duration, len(intervalSec))
)

func init() {
	var cur time.Duration
	for i, s := range intervalSec {
		cur += time.Duration(s) * time.Second
		transmitTime[i] = cur
	}
	if cur != sensorPeriod {
		log.Panicf("intervals sum to %v", cur)
	}
}

func main() {
	flag.Parse()
	papertrail.StartLogging()
	var err error
	sensorAddress, err = medtronic.DeviceAddress(*sensorID)
	if err != nil {
		log.Fatal(err)
	}
	readBGs()
	if len(nsEntries) < 2 {
		log.Fatal("not enough BG entries")
	}
	sendBGs()
}

func readBGs() {
	var err error
	nsEntries, err = nightscout.ReadEntries(*jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	nsEntries.Sort()
	if len(nsEntries) > ISIGsPerPacket {
		nsEntries = nsEntries[:ISIGsPerPacket]
	}
}

func sendBGs() {
	// Synchronize to sensor transmission schedule.
	e := nsEntries[0]
	log.Printf("latest BG: %d at %s", e.SGV, e.Time().Format(medtronic.UserTimeLayout))
	d := time.Duration(medtronic.SinceMidnight(time.Now()))
	m := d % sensorPeriod
	index := sort.Search(len(transmitTime), func(i int) bool {
		return transmitTime[i] >= m
	})
	var delta time.Duration
	if index < len(transmitTime) {
		delta = transmitTime[index] - m
	} else {
		index = 0
		delta = sensorPeriod - m
	}
	log.Printf("sleeping for %v", delta)
	time.Sleep(delta)

	for i := range isigBuffer {
		isigBuffer[i] = getISIG(i)
	}
	seq := getSequence(index)
	p := sensorPacket(seq)
	log.Printf("sending packet %02X", seq)
	pump := medtronic.Open()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	pump.Radio.Send(p)
	pump.Close()

	if index%2 == 1 {
		log.Printf("warning: missed initial transmission of this packet")
		return
	}
	delta = transmitTime[index+1] - transmitTime[index]
	log.Printf("sleeping for %v", delta)
	time.Sleep(delta)

	seq = getSequence(index + 1)
	p = sensorPacket(seq)
	log.Printf("sending packet %02X", seq)
	pump = medtronic.Open()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	pump.Radio.Send(p)
	pump.Close()
}

func getSequence(i int) byte {
	return byte((i/2)<<4 | i%2)
}

func getISIG(pos int) int {
	if pos >= len(nsEntries) {
		return 0
	}
	return nsEntries[pos].SGV * sensorCalibration
}
