package medtronic

import (
	"log"
)

// Carbs values are represented as either grams or 10x exchanges.
type Carbs int

const (
	CarbUnits    Command = 0x88
	GlucoseUnits Command = 0x89
)

type CarbUnitsType byte

//go:generate stringer -type CarbUnitsType

const (
	Grams     CarbUnitsType = 1
	Exchanges CarbUnitsType = 2
)

// Glucose values are represented as either mg/dL or μmol/L,
// so all conversions must include a GlucoseUnitsType parameter.
type Glucose int

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

func intToGlucose(n int, t GlucoseUnitsType) Glucose {
	if t == MmolPerLiter {
		// Convert 10x mmol/L to μmol/L
		return Glucose(n) * 100
	} else {
		return Glucose(n)
	}
}

func byteToGlucose(n byte, t GlucoseUnitsType) Glucose {
	return intToGlucose(int(n), t)
}

func (pump *Pump) CarbUnits() CarbUnitsType {
	return CarbUnitsType(pump.whichUnits(CarbUnits))
}

func (pump *Pump) GlucoseUnits() GlucoseUnitsType {
	return GlucoseUnitsType(pump.whichUnits(GlucoseUnits))
}

// Quantities and rates of insulin delivery are represented in milliunits.
type Insulin int

func milliUnitsPerStroke(newerPump bool) Insulin {
	if newerPump {
		return 25
	} else {
		return 100
	}
}

func intToInsulin(strokes int, newerPump bool) Insulin {
	return Insulin(strokes) * milliUnitsPerStroke(newerPump)
}

func byteToInsulin(strokes uint8, newerPump bool) Insulin {
	return intToInsulin(int(strokes), newerPump)
}

func twoByteInsulin(data []byte, newerPump bool) Insulin {
	return Insulin(twoByteUint(data)) * milliUnitsPerStroke(newerPump)
}
