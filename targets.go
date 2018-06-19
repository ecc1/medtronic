package medtronic

import (
	"time"
)

// GlucoseTarget represents an entry in a glucose target schedule.
type GlucoseTarget struct {
	Start TimeOfDay
	Low   Glucose
	High  Glucose
	Units GlucoseUnitsType
}

// GlucoseTargetSchedule represents a glucose target schedule.
type GlucoseTargetSchedule []GlucoseTarget

func glucoseTargetStep(family Family) int {
	if family <= 12 {
		return 2
	}
	return 3
}

func decodeGlucoseTargetSchedule(data []byte, units GlucoseUnitsType, family Family) GlucoseTargetSchedule {
	var sched []GlucoseTarget
	step := glucoseTargetStep(family)
	for i := 0; i <= len(data)-step; i += step {
		start := halfHoursToTimeOfDay(data[i])
		if start == 0 && len(sched) != 0 {
			break
		}
		low := byteToGlucose(data[i+1], units)
		high := low
		if family > 12 {
			high = byteToGlucose(data[i+2], units)
		}
		sched = append(sched, GlucoseTarget{
			Start: start,
			Low:   low,
			High:  high,
			Units: units,
		})
	}
	return sched
}

// GlucoseTargets returns the pump's glucose target schedule.
func (pump *Pump) GlucoseTargets() GlucoseTargetSchedule {
	// Command opcode and format of response depend on the pump family.
	family := pump.Family()
	var cmd Command
	if family <= 12 {
		cmd = glucoseTargets512
	} else {
		cmd = glucoseTargets
	}
	data := pump.Execute(cmd)
	if pump.Error() != nil {
		return GlucoseTargetSchedule{}
	}
	if len(data) < 2 {
		pump.BadResponse(cmd, data)
		return GlucoseTargetSchedule{}
	}
	n := int(data[0]) - 1
	step := glucoseTargetStep(family)
	if n%step != 0 {
		pump.BadResponse(cmd, data)
		return GlucoseTargetSchedule{}
	}
	units := GlucoseUnitsType(data[1])
	return decodeGlucoseTargetSchedule(data[2:2+n], units, family)
}

// GlucoseTargetAt returns the glucose target in effect at the given time.
func (s GlucoseTargetSchedule) GlucoseTargetAt(t time.Time) GlucoseTarget {
	d := SinceMidnight(t)
	last := GlucoseTarget{}
	for _, v := range s {
		if v.Start > d {
			break
		}
		last = v
	}
	return last
}
