package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ecc1/medtronic"
	"github.com/ecc1/medtronic/packet"
)

var (
	meterID = flag.String("m", "123456", "meter `ID`")

	meterAddress []byte

	noResponse = errors.New("no response")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [options] bg", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	var err error
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		usage()
	}
	bg, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		usage()
	}
	meterAddress, err = medtronic.DeviceAddress(*meterID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		usage()
	}

	p := meterPacket(bg)

	pump := medtronic.Open()
	defer pump.Close()

	pump.Wakeup()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}

	pump.SetTimeout(1500 * time.Millisecond)
	for tries := 0; tries < pump.Retries(); tries++ {
		pump.SetError(nil)
		sendPacket(pump, p)
		err := pump.Error()
		if err == nil {
			break
		}
		if err != noResponse {
			log.Print(err)
		}
	}
	err = pump.Error()
	if err != nil {
		if err == noResponse {
			log.Print(err)
		}
		os.Exit(1)
	}
}

func sendPacket(pump *medtronic.Pump, p []byte) {
	response, _ := pump.Radio.SendAndReceive(p, pump.Timeout())
	if pump.Error() != nil {
		return
	}
	if len(response) == 0 {
		pump.SetError(noResponse)
		return
	}
	data, err := packet.Decode(response)
	if err != nil {
		pump.SetError(err)
		return
	}
	if !isAck(data) {
		pump.SetError(fmt.Errorf("unexpected response: % X", data))
	}
}

func meterPacket(bg int) []byte {
	p := make([]byte, 6)
	p[0] = packet.Meter
	copy(p[1:4], meterAddress)
	p[4] = byte(bg>>8) & 0x1
	p[5] = byte(bg)
	return packet.Encode(p)
}

func isAck(data []byte) bool {
	return data[0] == packet.Meter &&
		bytes.Equal(data[1:4], meterAddress) &&
		data[4] == 0x06
}
