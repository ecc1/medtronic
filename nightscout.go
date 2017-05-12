package medtronic

import (
	"log"
	"time"

	"github.com/ecc1/nightscout"
)

var (
	eventType = map[HistoryRecordType]string{
		Bolus:         "Meal Bolus",
		BGCapture:     "BG Check",
		SuspendPump:   "Temp Basal",
		ResumePump:    "Temp Basal",
		Rewind:        "Site Change",
		TempBasalRate: "Temp Basal",
	}
)

// Treatments converts certain pump history records
// into records that can be uploaded as Nightscout treatments.
// History records must be in chronological order.
func Treatments(records []HistoryRecord) []nightscout.Treatment {
	var treatments []nightscout.Treatment
	user := nightscout.Username()
	for i, r := range records {
		var r2 *HistoryRecord
		if i+1 < len(records) {
			r2 = &records[i+1]
		}
		info := nightscout.Treatment{
			CreatedAt: time.Time(r.Time),
			EnteredBy: user,
		}
		if getRecordInfo(r, r2, &info) {
			treatments = append(treatments, info)
		}
	}
	return treatments
}

func getRecordInfo(r HistoryRecord, r2 *HistoryRecord, info *nightscout.Treatment) bool {
	t := r.Type()
	info.EventType = eventType[t]
	switch t {
	case BGCapture:
		g := r.Glucose.NightscoutGlucose()
		info.Glucose = &g
		info.Units = "mg/dl"
	case TempBasalRate:
		return tempBasalInfo(r, r2, info)
	case Bolus:
		ins := r.Bolus.Amount.NightscoutInsulin()
		info.Insulin = &ins
		min := int(r.Bolus.Duration / Duration(time.Minute))
		info.Duration = &min
	case Rewind:
		if !nextEvent(r, r2, Prime) {
			return false
		}
	case ResumePump:
		insulin0 := Insulin(0).NightscoutInsulin()
		info.Absolute = &insulin0
		duration0 := 0
		info.Duration = &duration0
	case SuspendPump:
		insulin0 := Insulin(0).NightscoutInsulin()
		info.Absolute = &insulin0
		min := 24 * 60
		info.Duration = &min
	default:
		return false
	}
	return true
}

func tempBasalInfo(r HistoryRecord, r2 *HistoryRecord, info *nightscout.Treatment) bool {
	if !nextEvent(r, r2, TempBasalDuration) {
		return false
	}
	if *r2.Duration == 0 {
		insulin0 := Insulin(0).NightscoutInsulin()
		info.Absolute = &insulin0
		duration0 := 0
		info.Duration = &duration0
	} else {
		ins := r.Insulin.NightscoutInsulin()
		info.Absolute = &ins
		min := int(*r2.Duration / Duration(time.Minute))
		info.Duration = &min
	}
	return true
}

func nextEvent(r HistoryRecord, r2 *HistoryRecord, t HistoryRecordType) bool {
	if r2 == nil {
		ts := time.Time(r.Time).Format(UserTimeLayout)
		log.Printf("expected %v to be followed by %v at %s", r.Type(), t, ts)
		return false
	}
	if r2.Type() != t {
		ts := time.Time(r.Time).Format(UserTimeLayout)
		log.Printf("expected %v to be followed by %v at %s but found %v", r.Type(), t, ts, r2.Type())
		return false
	}
	return true
}

// NightscoutGlucose converts a Glucose value to a nightscout.Glucose value.
func (r Glucose) NightscoutGlucose() nightscout.Glucose {
	return nightscout.Glucose(r)
}

// NightscoutInsulin converts an Insulin value to a nightscout.Insulin value.
func (r Insulin) NightscoutInsulin() nightscout.Insulin {
	return nightscout.Insulin(float64(r) / 1000)
}

// NightscoutVoltage converts a Voltage value to a nightscout.Voltage value.
func (r Voltage) NightscoutVoltage() nightscout.Voltage {
	return nightscout.Voltage(float64(r) / 1000)
}

// NightscoutSchedule converts a BasalRateSchedule to a nightscout.Schedule.
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

// NightscoutSchedule converts a CarbRatioSchedule to a nightscout.Schedule.
func (sched CarbRatioSchedule) NightscoutSchedule() nightscout.Schedule {
	n := len(sched)
	if n != 0 && sched[0].Units != Grams {
		panic("carb units must be grams")
	}
	tv := make(nightscout.Schedule, n)
	for i, r := range sched {
		tv[i] = nightscout.TimeValue{
			Time:  r.Start.String(),
			Value: float64(r.Ratio) / 10, // Grams
		}
	}
	return tv
}

// NightscoutSchedule converts an InsulinSensitivitySchedule to a nightscout.Schedule.
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

// NightscoutSchedule converts a GlucoseTargetSchedule to a nightscout.Schedule.
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
