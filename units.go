package medtronic

import (
	"log"
)

const (
	CarbUnits    Command = 0x88
	GlucoseUnits Command = 0x89
)

type CarbUnitsType byte

const (
	Grams     CarbUnitsType = 1
	Exchanges CarbUnitsType = 2
)

func (u CarbUnitsType) String() string {
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

type GlucoseUnitsType byte

const (
	MgPerDeciLiter GlucoseUnitsType = 1
	MmolPerLiter   GlucoseUnitsType = 2
)

func (u GlucoseUnitsType) String() string {
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

func (pump *Pump) whichUnits(cmd Command) byte {
	data := pump.Execute(cmd)
	if pump.Error() != nil {
		return 0
	}
	if len(data) < 2 || data[0] != 1 {
		pump.BadResponse(cmd, data)
		return 0
	}
	return data[1]
}

func (pump *Pump) CarbUnits() CarbUnitsType {
	return CarbUnitsType(pump.whichUnits(CarbUnits))
}

func (pump *Pump) GlucoseUnits() GlucoseUnitsType {
	return GlucoseUnitsType(pump.whichUnits(GlucoseUnits))
}
