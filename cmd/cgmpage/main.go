package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ecc1/medtronic"
)

var (
	glucosePage = flag.Int("g", -1, "read glucose history `page`")
	isigPage    = flag.Int("i", -1, "read ISIG history `page`")
	vcntrPage   = flag.Int("v", -1, "read vcntr history `page`")
	calFactor   = flag.Bool("c", false, "read calibration factor")
)

func main() {
	flag.Parse()
	pump := medtronic.Open()
	defer pump.Close()
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
		f := pump.CalibrationFactor()
		if pump.Error() != nil {
			log.Fatal(pump.Error())
		}
		fmt.Println(f)
		return
	}
	m, n := pump.CGMPageRange()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	fmt.Printf("[%d, %d)\n", m, n)
}
