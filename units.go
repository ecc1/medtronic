package medtronic

import (
	"log"
)

const (
	CarbUnits    CommandCode = 0x88
	GlucoseUnits CommandCode = 0x89
)

type CarbUnitsInfo byte

const (
	Grams     CarbUnitsInfo = 1
	Exchanges CarbUnitsInfo = 2
)

func (u CarbUnitsInfo) String() string {
	switch u {
	case Grams:
		return "grams"
	case Exchanges:
		return "exchanges"
	default:
		log.Panicf("unknown carb unit %d", u)
	}
	panic("unreachable")
}

type GlucoseUnitsInfo byte

const (
	MgPerDeciLiter GlucoseUnitsInfo = 1
	MmolPerLiter   GlucoseUnitsInfo = 2
)

func (u GlucoseUnitsInfo) String() string {
	switch u {
	case MgPerDeciLiter:
		return "mg/dL"
	case MmolPerLiter:
		return "mmol/L"
	default:
		log.Panicf("unknown glucose unit %d", u)
	}
	panic("unreachable")
}

func (pump *Pump) whichUnits(cmd CommandCode) byte {
	result := pump.Execute(cmd, func(data []byte) interface{} {
		if len(data) < 2 || data[0] != 1 {
			return nil
		}
		switch data[1] {
		case 1:
			return byte(1)
		case 2:
			return byte(2)
		default:
			return nil
		}
	})
	if pump.Error() != nil {
		return 0
	}
	return result.(byte)
}

func (pump *Pump) CarbUnits() CarbUnitsInfo {
	return CarbUnitsInfo(pump.whichUnits(CarbUnits))
}

func (pump *Pump) GlucoseUnits() GlucoseUnitsInfo {
	return GlucoseUnitsInfo(pump.whichUnits(GlucoseUnits))
}
