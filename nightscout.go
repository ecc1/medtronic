package medtronic

import (
	"log"
	"time"

	"github.com/ecc1/nightscout"
)

func (r HistoryRecord) NightscoutTreatment(r2 *HistoryRecord) *nightscout.Treatment {
	info := nightscout.Treatment{
		EventTime: r.Time,
		EnteredBy: nightscout.Username(),
	}
	switch r.Type() {
	case BgCapture:
		info.EventType = "BG Check"
		g := r.Glucose.NightscoutGlucose()
		info.Glucose = &g
	case TempBasalRate:
		if !nextEvent(r, r2, TempBasalDuration) {
			return nil
		}
		if *r2.Duration == 0 {
			info.EventType = "Temp Basal End"
		} else {
			info.EventType = "Temp Basal Start"
			ins := r.Insulin.NightscoutInsulin()
			info.Absolute = &ins
			min := int(*r2.Duration / time.Minute)
			info.Duration = &min
		}
	case Bolus:
		info.EventType = "Meal Bolus"
		ins := r.Bolus.Amount.NightscoutInsulin()
		info.Insulin = &ins
		min := int(r.Bolus.Duration / time.Minute)
		info.Duration = &min
	case Rewind:
		if !nextEvent(r, r2, Prime) {
			return nil
		}
		info.EventType = "Site Change"
	case ResumePump:
		info.EventType = "Temp Basal End"
	case SuspendPump:
		info.EventType = "Temp Basal Start"
		zero := Insulin(0).NightscoutInsulin()
		info.Absolute = &zero
		min := 24 * 60
		info.Duration = &min
	default:
		return nil
	}
	return &info
}

func nextEvent(r HistoryRecord, r2 *HistoryRecord, t HistoryRecordType) bool {
	if r2 == nil || r2.Type() != t {
		next := "nothing"
		if r2 != nil {
			next = r2.Type().String()
		}
		log.Printf("expected %v to be followed by %v at %v but found %s", r.Type(), t, r.Time, next)
		return false
	}
	return true
}

func (r Glucose) NightscoutGlucose() nightscout.Glucose {
	return nightscout.Glucose(r)
}

func (r Insulin) NightscoutInsulin() nightscout.Insulin {
	return nightscout.Insulin(float64(r) / 1000)
}

func (r Voltage) NightscoutVoltage() nightscout.Voltage {
	return nightscout.Voltage(float64(r) / 1000)
}

func (sched BasalRateSchedule) NightscoutSchedule() nightscout.Schedule {
	n := len(sched)
	tv := make(nightscout.Schedule, n)
	for i, r := range sched {
		tv[i] = nightscout.TimeValue{
			Time:  r.Start.String(),
			Value: r.Rate,
		}
	}
	return tv
}

func (sched CarbRatioSchedule) NightscoutSchedule() nightscout.Schedule {
	n := len(sched)
	tv := make(nightscout.Schedule, n)
	for i, r := range sched {
		tv[i] = nightscout.TimeValue{
			Time:  r.Start.String(),
			Value: r.CarbRatio,
		}
	}
	return tv
}

func (sched InsulinSensitivitySchedule) NightscoutSchedule() nightscout.Schedule {
	n := len(sched)
	tv := make(nightscout.Schedule, n)
	for i, r := range sched {
		tv[i] = nightscout.TimeValue{
			Time:  r.Start.String(),
			Value: r.Sensitivity,
		}
	}
	return tv
}

func (sched GlucoseTargetSchedule) NightscoutSchedule() (nightscout.Schedule, nightscout.Schedule) {
	n := len(sched)
	low := make(nightscout.Schedule, n)
	high := make(nightscout.Schedule, n)
	for i, r := range sched {
		t := r.Start.String()
		low[i] = nightscout.TimeValue{
			Time:  t,
			Value: r.Low,
		}
		high[i] = nightscout.TimeValue{
			Time:  t,
			Value: r.High,
		}
	}
	return low, high
}
