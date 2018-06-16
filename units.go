package medtronic

import (
	"fmt"
	"log"
)

// Carbs represents a carb value as either grams or 10x exchanges.
type Carbs int

// CarbUnitsType represents the pump's carb unit type (grams or exchanges).
type CarbUnitsType byte

//go:generate stringer -type CarbUnitsType

const (
	// Grams represents the pump's use of grams for carb units.
	Grams CarbUnitsType = 1
	// Exchanges represents the pump's use of exchanges for carb units.
	Exchanges CarbUnitsType = 2
)

// Glucose represents a glucose value as either mg/dL or μmol/L,
// so all conversions must include a GlucoseUnitsType parameter.
type Glucose int

// GlucoseUnitsType represents the pump's glucose unit type (mg/dL or mmol/L).
type GlucoseUnitsType byte

const (
	// MgPerDeciLiter represents the pump's use of mg/dL for glucose units.
	MgPerDeciLiter GlucoseUnitsType = 1
	// MMolPerLiter represents the pump's use of mmol/L for glucose units.
	MMolPerLiter GlucoseUnitsType = 2
)

func (u GlucoseUnitsType) String() string {
	switch u {
	case MgPerDeciLiter:
		return "mg/dL"
	case MMolPerLiter:
		return "μmol/L"
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
	switch t {
	case MgPerDeciLiter:
		return Glucose(n)
	case MMolPerLiter:
		// Convert 10x mmol/L to μmol/L
		return Glucose(n) * 100
	default:
		log.Panicf("unknown glucose unit %d", t)
	}
	panic("unreachable")
}

func byteToGlucose(n byte, t GlucoseUnitsType) Glucose {
	return intToGlucose(int(n), t)
}

// CarbUnits returns the pump's carb units.
func (pump *Pump) CarbUnits() CarbUnitsType {
	return CarbUnitsType(pump.whichUnits(carbUnits))
}

// GlucoseUnits returns the pump's glucose units.
func (pump *Pump) GlucoseUnits() GlucoseUnitsType {
	return GlucoseUnitsType(pump.whichUnits(glucoseUnits))
}

// Insulin represents quantities and rates of insulin delivery, in milliunits.
type Insulin int

func (r Insulin) String() string {
	return fmt.Sprintf("%g", float64(r)/1000)
}

func milliUnitsPerStroke(family Family) Insulin {
	if family <= 22 {
		return 100
	}
	return 25
}

func intToInsulin(strokes int, family Family) Insulin {
	return Insulin(strokes) * milliUnitsPerStroke(family)
}

func byteToInsulin(strokes uint8, family Family) Insulin {
	return intToInsulin(int(strokes), family)
}

func twoByteInsulin(data []byte, family Family) Insulin {
	return intToInsulin(twoByteInt(data), family)
}
