package medtronic

import (
	"fmt"
	"log"
	"time"
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
	for i := 0; i < len(data)-2; i += 3 {
		rate := twoByteInsulinLE(data[i : i+2])
		t := data[i+2]
		// Don't stop if the 00:00 rate happens to be zero.
		if i > 1 && rate == 0 && t == 0 {
			break
		}
		start := halfHoursToTimeOfDay(t)
		sched = append(sched, BasalRate{Start: start, Rate: rate})
	}
	return sched
}

func (pump *Pump) basalSchedule(cmd Command) BasalRateSchedule {
	data := pump.ExtendedResponse(cmd)
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
	d := SinceMidnight(t)
	last := BasalRate{}
	for _, v := range s {
		if v.Start > d {
			break
		}
		last = v
	}
	return last
}

func (pump *Pump) setBasalSchedule(cmd Command, s BasalRateSchedule) {
	if len(s) == 0 {
		pump.SetError(fmt.Errorf("%v: empty schedule", cmd))
		return
	}
	pump.ExtendedRequest(setBasalRates, s.Encode()...)
}

func (s BasalRateSchedule) Encode() []byte {
	data := make([]byte, len(s)*3)
	i := 0
	for _, v := range s {
		m := milliUnitsPerStroke(23)
		strokes := v.Rate / m
		actual := strokes * m
		if actual != v.Rate {
			log.Printf("rounding basal rate from %v to %v", v.Rate, actual)
		}
		copy(data[i:i+2], marshalUint16LE(uint16(strokes)))
		data[i+2] = v.Start.HalfHours()
		i += 3
	}
	return data
}

// SetBasalRates sets the pump's basal rate schedule.
func (pump *Pump) SetBasalRates(s BasalRateSchedule) {
	pump.setBasalSchedule(setBasalRates, s)
}

// SetBasalPatternA sets the pump's basal pattern A.
func (pump *Pump) SetBasalPatternA(s BasalRateSchedule) {
	pump.setBasalSchedule(setBasalPatternA, s)
}

// SetBasalPatternB sets the pump's basal pattern B.
func (pump *Pump) SetBasalPatternB(s BasalRateSchedule) {
	pump.setBasalSchedule(setBasalPatternB, s)
}
