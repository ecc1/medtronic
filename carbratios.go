package medtronic

import (
	"log"
	"time"
)

const (
	carbRatios Command = 0x8A
)

// CarbRatio represents an entry in a carb ratio schedule.
type CarbRatio struct {
	Start TimeOfDay
	Ratio Ratio
	Units CarbUnitsType
}

// Newer pumps store carb ratios as 10x grams/unit or 1000x units/exchange.
// Older pumps store carb ratios as grams/unit or 10x units/exchange.

// Ratio represents a carb ratio using the higher resolution:
// 10x grams/unit or 1000x units/exchange.
type Ratio int

func intToRatio(n int, u CarbUnitsType, family Family) Ratio {
	if family <= 22 {
		// Convert to higher-resolution representation.
		switch u {
		case Grams:
			return Ratio(10 * n)
		case Exchanges:
			return Ratio(100 * n)
		default:
			log.Panicf("unknown carb unit %d", u)
		}
	}
	// Use representation as-is.
	return Ratio(n)
}

// CarbRatioSchedule represents a carb ratio schedule.
type CarbRatioSchedule []CarbRatio

func carbRatioStep(family Family) int {
	if family <= 22 {
		return 2
	}
	return 3
}

func decodeCarbRatioSchedule(data []byte, units CarbUnitsType, family Family) CarbRatioSchedule {
	var sched []CarbRatio
	step := carbRatioStep(family)
	for i := 0; i <= len(data)-step; i += step {
		start := halfHoursToTimeOfDay(data[i])
		if start == 0 && len(sched) != 0 {
			break
		}
		var value int
		if family <= 22 {
			value = int(data[i+1])
		} else {
			value = twoByteInt(data[i+1 : i+3])
		}
		sched = append(sched, CarbRatio{
			Start: start,
			Ratio: intToRatio(value, units, family),
			Units: units,
		})
	}
	return sched
}

// CarbRatios returns the pump's carb ratio schedule..
func (pump *Pump) CarbRatios() CarbRatioSchedule {
	data := pump.Execute(carbRatios)
	if pump.Error() != nil {
		return CarbRatioSchedule{}
	}
	if len(data) < 2 {
		pump.BadResponse(carbRatios, data)
		return CarbRatioSchedule{}
	}
	// Format of response depends on the pump family.
	family := pump.Family()
	n := int(data[0]) - 1
	step := carbRatioStep(family)
	if n%step != 0 {
		pump.BadResponse(carbRatios, data)
		return CarbRatioSchedule{}
	}
	units := CarbUnitsType(data[1])
	return decodeCarbRatioSchedule(data[step:step+n], units, family)
}

// CarbRatioAt returns the carb ratio in effect at the given time.
func (s CarbRatioSchedule) CarbRatioAt(t time.Time) CarbRatio {
	d := sinceMidnight(t)
	last := CarbRatio{}
	for _, v := range s {
		if v.Start > d {
			break
		}
		last = v
	}
	return last
}
