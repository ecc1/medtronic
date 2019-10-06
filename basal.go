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

func encodeBasalRateSchedule(s BasalRateSchedule, family Family) ([]byte, error) {
	data := make([]byte, 0, len(s)*3)
	for _, v := range s {
		r, err := encodeBasalRate("basal", v.Rate, family)
		if err != nil {
			return nil, err
		}
		b := append(marshalUint16LE(r), v.Start.HalfHours())
		data = append(data, b...)
	}
	return data, nil
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

// BasalRateAt returns the index of the basal rate in effect at the given time.
func (s BasalRateSchedule) BasalRateAt(t time.Time) int {
	d := SinceMidnight(t)
	last := -1
	for i, v := range s {
		if v.Start > d {
			break
		}
		last = i
	}
	if last == -1 {
		// Schedule started after time t?
		log.Printf("cannot find basal rate at %s in profile %+v", t.Format(UserTimeLayout), s)
	}
	return last
}

// NextChange returns the time when the next scheduled rate will take effect (strictly after t).
func (s BasalRateSchedule) NextChange(t time.Time) time.Time {
	i := s.BasalRateAt(t)
	var next time.Duration
	if i+1 < len(s) {
		next = time.Duration(s[i+1].Start)
	} else {
		next = 24 * time.Hour
	}
	return t.Add(next - time.Duration(SinceMidnight(t)))
}

func (pump *Pump) setBasalSchedule(cmd Command, s BasalRateSchedule) {
	if len(s) == 0 {
		pump.SetError(fmt.Errorf("%v: empty schedule", cmd))
		return
	}
	data, err := encodeBasalRateSchedule(s, pump.Family())
	if err != nil {
		pump.SetError(err)
		return
	}
	pump.ExtendedRequest(setBasalRates, data...)
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

func encodeBasalRate(kind string, rate Insulin, family Family) (uint16, error) {
	if rate < 0 {
		return 0, fmt.Errorf("%s rate (%d) is negative", kind, rate)
	}
	if rate > maxBasal {
		return 0, fmt.Errorf("%s rate (%d) is too large", kind, rate)
	}
	// Round the rate to the pump's delivery resolution.
	var res Insulin
	if family <= 22 {
		res = 50
	} else if rate < 1000 {
		res = 25
	} else if rate < 10000 {
		res = 50
	} else {
		res = 100
	}
	actual := (rate / res) * res
	if actual != rate {
		log.Printf("rounding %s rate from %v to %v", kind, rate, actual)
	}
	// Encode the rounded value using 25 milliUnits/stroke.
	m := milliUnitsPerStroke(23)
	return uint16(actual / m), nil
}
