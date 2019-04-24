package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ecc1/medtronic"
)

var (
	glucosePage    = flag.Int("g", -1, "read glucose history `page`")
	isigPage       = flag.Int("i", -1, "read ISIG history `page`")
	vcntrPage      = flag.Int("v", -1, "read vcntr history `page`")
	calFactor      = flag.Bool("c", false, "read calibration factor")
	writeTimestamp = flag.Bool("w", false, "write timestamp")
)

func main() {
	flag.Parse()
	pump := medtronic.Open()
	defer pump.Close()
	pump.Wakeup()
	var data []byte
	if *glucosePage >= 0 {
		data = pump.GlucosePage(*glucosePage)
	} else if *isigPage >= 0 {
		data = pump.ISIGPage(*isigPage)
	} else if *vcntrPage >= 0 {
		data = pump.VcntrPage(*vcntrPage)
	}
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	if len(data) != 0 {
		fmt.Printf("% X\n", data)
		return
	}
	if *calFactor {
		getCal(pump)
		return
	}
	if *writeTimestamp {
		writeTS(pump)
		return
	}
	getCur(pump)
}

func getCal(pump *medtronic.Pump) {
	f := pump.CalibrationFactor()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	fmt.Println(f)
}

func writeTS(pump *medtronic.Pump) {
	pump.CGMWriteTimestamp()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
}

func getCur(pump *medtronic.Pump) {
	n := pump.CGMCurrentGlucosePage()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	fmt.Println(n)
}
