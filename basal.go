package medtronic

import (
	"time"
)

const (
	BasalRates    CommandCode = 0x92
	BasalPatternA CommandCode = 0x93
	BasalPatternB CommandCode = 0x94
)

type BasalRate struct {
	Start time.Duration // offset from 00:00:00
	Rate  int           // milliUnits per hour
}

type BasalRateSchedule struct {
	Schedule []BasalRate
}

func (pump *Pump) basalSchedule(cmd CommandCode) BasalRateSchedule {
	result := pump.Execute(cmd, func(data []byte) interface{} {
		info := []BasalRate{}
		for i := 1; i < len(data); i += 3 {
			r := data[i]
			t := data[i+2]
			// Don't stop if the 00:00 rate happens to be zero.
			if i > 1 && r == 0 && t == 0 {
				break
			}
			start := scheduleToDuration(t)
			rate := int(r) * 25
			info = append(info, BasalRate{Start: start, Rate: rate})
		}
		return BasalRateSchedule{Schedule: info}
	})
	if pump.Error() != nil {
		return BasalRateSchedule{}
	}
	return result.(BasalRateSchedule)
}

// Convert a multiple of half-hours to a Duration.
func scheduleToDuration(t uint8) time.Duration {
	return time.Duration(t) * 30 * time.Minute
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
	for _, v := range s.Schedule {
		if v.Start > d {
			break
		}
		last = v
	}
	return last
}

// Convert a time to a Duration representing the offset since 00:00:00.
func sinceMidnight(t time.Time) time.Duration {
	yy, mm, dd := t.Date()
	midnight := time.Date(yy, mm, dd, 0, 0, 0, 0, t.Location())
	return t.Sub(midnight)
}
