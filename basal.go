package medtronic

import (
	"time"
)

const (
	basalRates    Command = 0x92
	basalPatternA Command = 0x93
	basalPatternB Command = 0x94
)

// BasalRate represents an entry in a basal rate schedule.
type BasalRate struct {
	Start TimeOfDay
	Rate  Insulin
}

// BasalRateSchedule represents a basal rate schedule.
type BasalRateSchedule []BasalRate

func decodeBasalRateSchedule(data []byte) BasalRateSchedule {
	var sched []BasalRate
	for i := 1; i < len(data); i += 3 {
		r := data[i]
		t := data[i+2]
		// Don't stop if the 00:00 rate happens to be zero.
		if i > 1 && r == 0 && t == 0 {
			break
		}
		start := halfHoursToTimeOfDay(t)
		rate := byteToInsulin(r, true)
		sched = append(sched, BasalRate{Start: start, Rate: rate})
	}
	return sched
}

func (pump *Pump) basalSchedule(cmd Command) BasalRateSchedule {
	data := pump.Execute(cmd)
	if pump.Error() != nil {
		return BasalRateSchedule{}
	}
	return decodeBasalRateSchedule(data)
}

// BasalRates returns the pump's basal rate schedule.
func (pump *Pump) BasalRates() BasalRateSchedule {
	return pump.basalSchedule(basalRates)
}

// BasalPatternA returns the pump's basal pattern A.
func (pump *Pump) BasalPatternA() BasalRateSchedule {
	return pump.basalSchedule(basalPatternA)
}

// BasalPatternB returns the pump's basal pattern B.
func (pump *Pump) BasalPatternB() BasalRateSchedule {
	return pump.basalSchedule(basalPatternB)
}

// BasalRateAt returns the basal rate in effect at the given time.
func (s BasalRateSchedule) BasalRateAt(t time.Time) BasalRate {
	d := sinceMidnight(t)
	last := BasalRate{}
	for _, v := range s {
		if v.Start > d {
			break
		}
		last = v
	}
	return last
}
