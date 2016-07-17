package medtronic

import (
	"time"
)

const (
	BasalRates    Command = 0x92
	BasalPatternA Command = 0x93
	BasalPatternB Command = 0x94
)

type BasalRate struct {
	Start TimeOfDay
	Rate  Insulin
}

type BasalRateSchedule []BasalRate

func (pump *Pump) basalSchedule(cmd Command) BasalRateSchedule {
	data := pump.Execute(cmd)
	if pump.Error() != nil {
		return BasalRateSchedule{}
	}
	sched := []BasalRate{}
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

func (pump *Pump) BasalRates() BasalRateSchedule {
	return pump.basalSchedule(BasalRates)
}

func (pump *Pump) BasalPatternA() BasalRateSchedule {
	return pump.basalSchedule(BasalPatternA)
}

func (pump *Pump) BasalPatternB() BasalRateSchedule {
	return pump.basalSchedule(BasalPatternB)
}

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
