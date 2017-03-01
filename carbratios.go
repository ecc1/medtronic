package medtronic

import (
	"log"
	"time"
)

const (
	CarbRatios Command = 0x8A
)

type CarbRatio struct {
	Start TimeOfDay
	Ratio Ratio
	Units CarbUnitsType
}

// Newer pumps store carb ratios as 10x grams/unit or 1000x units/exchange.
// Older pumps store carb ratios as grams/unit or 10x units/exchange.

// Higher-resolution representation: 10x grams/unit or 1000x units/exchange.
type Ratio int

func intToRatio(n int, u CarbUnitsType, newerPump bool) Ratio {
	if newerPump {
		// Use representation as-is.
		return Ratio(n)
	}
	// Convert to higher-resolution representation.
	switch u {
	case Grams:
		return Ratio(10 * n)
	case Exchanges:
		return Ratio(100 * n)
	default:
		log.Panicf("unknown carb unit %d", u)
	}
	panic("unreachable")
}

type CarbRatioSchedule []CarbRatio

func carbRatioStep(newerPump bool) int {
	if newerPump {
		return 3
	} else {
		return 2
	}
}

func decodeCarbRatioSchedule(data []byte, units CarbUnitsType, newerPump bool) CarbRatioSchedule {
	var sched []CarbRatio
	step := carbRatioStep(newerPump)
	for i := 0; i < len(data); i += step {
		start := halfHoursToTimeOfDay(data[i])
		if start == 0 && len(sched) != 0 {
			break
		}
		value := 0
		if newerPump {
			value = twoByteInt(data[i+1 : i+3])
		} else {
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

func (pump *Pump) CarbRatios() CarbRatioSchedule {
	// Format of response depends on the pump family.
	newer := pump.Family() >= 23
	data := pump.Execute(CarbRatios)
	if pump.Error() != nil {
		return CarbRatioSchedule{}
	}
	if len(data) < 2 {
		pump.BadResponse(CarbRatios, data)
		return CarbRatioSchedule{}
	}
	n := int(data[0]) - 1
	step := carbRatioStep(newer)
	if n%step != 0 {
		pump.BadResponse(CarbRatios, data)
		return CarbRatioSchedule{}
	}
	units := CarbUnitsType(data[1])
	return decodeCarbRatioSchedule(data[step:step+n], units, newer)
}

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
