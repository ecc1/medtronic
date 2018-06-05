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
	sensorVersion     = 13
	sensorCalibration = 10
)

var (
	sensorID      = flag.String("s", "000000", "sensor `ID`")
	sensorAddress []byte

	jsonFile = flag.String("f", "glucose.json", "read BGs from JSON `file`")

	nsEntries nightscout.Entries
)

func main() {
	flag.Parse()
	papertrail.StartLogging()
	var err error
	sensorAddress, err = medtronic.DeviceAddress(*sensorID)
	if err != nil {
		log.Fatal(err)
	}
	readBGs()
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

var retransmitInterval = []int8{16, 17, 17, 11, 13, 12, 14, 11}

func sendBGs() {
	pump := medtronic.Open()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	defer pump.Close()

	pump.Wakeup()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	p := sensorPacket()
	for n := 0; n < 2; n++ {
		if n != 0 {
			r := time.Duration(retransmitInterval[p.Seq&0x7])
			time.Sleep(r * time.Second)
		}
		p.Repeat = n
		log.Printf("sending %+v", p)
		pump.Radio.Send(marshalSensorPacket(p, 0x3))
	}
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
}

func sensorPacket() SensorPacket {
	p := SensorPacket{Address: sensorAddress, Version: sensorVersion}
	p.Adjust = 0x1D
	p.Battery = 160
	p.Seq = getSequence()
	p.ISIG = make([]int, ISIGsPerPacket)
	for i := range p.ISIG {
		p.ISIG[i] = getISIG(i)
	}
	return p
}

func getSequence() int {
	e := nsEntries[0]
	n := e.Date / (5 * 60 * 1000)
	return int(n & 0xF)
}

func getISIG(pos int) int {
	if pos >= len(nsEntries) {
		return 0
	}
	return nsEntries[pos].SGV * sensorCalibration
}
