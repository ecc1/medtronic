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

func intToRatio(n int, u CarbUnitsType, newerPump bool) Ratio {
	switch newerPump {
	case true:
		// Use representation as-is.
		return Ratio(n)
	case false:
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
	panic("unreachable")
}

// CarbRatioSchedule represents a carb ratio schedule.
type CarbRatioSchedule []CarbRatio

func carbRatioStep(newerPump bool) int {
	switch newerPump {
	case true:
		return 3
	case false:
		return 2
	}
	panic("unreachable")
}

func decodeCarbRatioSchedule(data []byte, units CarbUnitsType, newerPump bool) CarbRatioSchedule {
	var sched []CarbRatio
	step := carbRatioStep(newerPump)
	for i := 0; i < len(data); i += step {
		start := halfHoursToTimeOfDay(data[i])
		if start == 0 && len(sched) != 0 {
			break
		}
		var value int
		switch newerPump {
		case true:
			value = twoByteInt(data[i+1 : i+3])
		case false:
			value = int(data[i+1])
		}
		sched = append(sched, CarbRatio{
			Start: start,
			Ratio: intToRatio(value, units, newerPump),
			Units: units,
		})
	}
	return sched
}

// CarbRatios returns the pump's carb ratio schedule..
func (pump *Pump) CarbRatios() CarbRatioSchedule {
	// Format of response depends on the pump family.
	newer := pump.Family() >= 23
	data := pump.Execute(carbRatios)
	if pump.Error() != nil {
		return CarbRatioSchedule{}
	}
	if len(data) < 2 {
		pump.BadResponse(carbRatios, data)
		return CarbRatioSchedule{}
	}
	n := int(data[0]) - 1
	step := carbRatioStep(newer)
	if n%step != 0 {
		pump.BadResponse(carbRatios, data)
		return CarbRatioSchedule{}
	}
	units := CarbUnitsType(data[1])
	return decodeCarbRatioSchedule(data[step:step+n], units, newer)
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
