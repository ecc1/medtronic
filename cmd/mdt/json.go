package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ecc1/medtronic"
)

func showInternal(v interface{}) {
	fmt.Printf("%+v\n", v)
}

func showJSON(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println(err)
		fmt.Println(v)
		return
	}
	fmt.Println(string(b))
}

func showOpenAPS(v interface{}) {
	showJSON(OpenAPSJSON(v))
}

// OpenAPSJSON converts v into a value that the json package
// will marshal in a format compatible with openaps.
func OpenAPSJSON(v interface{}) interface{} {
	switch r := v.(type) {
	case medtronic.BasalRateSchedule:
		return convertBasalRateSchedule(r)
	case medtronic.BatteryInfo:
		return convertBatteryInfo(r)
	case medtronic.CarbRatioSchedule:
		return convertCarbRatioSchedule(r)
	case medtronic.CarbUnitsType:
		return carbU(r)
	case medtronic.GlucoseTargetSchedule:
		return convertGlucoseTargetSchedule(r)
	case medtronic.GlucoseUnitsType:
		return glucoseU(r)
	case medtronic.InsulinSensitivitySchedule:
		return convertInsulinSensitivitySchedule(r)
	case medtronic.SettingsInfo:
		return convertSettingsInfo(r)
	case medtronic.StatusInfo:
		return convertStatusInfo(r)
	case medtronic.TempBasalInfo:
		return convertTempBasalInfo(r)
	default:
		return v
	}
}

func hours(d time.Duration) int {
	return int((d + 30*time.Minute) / time.Hour)
}

func minutes(d time.Duration) int {
	return int((d + 30*time.Second) / time.Minute)
}

func carbU(u medtronic.CarbUnitsType) string {
	switch u {
	case medtronic.Grams:
		return "grams"
	case medtronic.Exchanges:
		return "exchanges"
	default:
		log.Panicf("unknown carb unit %d", u)
	}
	panic("unreachable")
}

func glucoseU(u medtronic.GlucoseUnitsType) string {
	switch u {
	case medtronic.MgPerDeciLiter:
		return "mg/dL"
	case medtronic.MMolPerLiter:
		return "mmol/L"
	default:
		log.Panicf("unknown glucose unit %d", u)
	}
	panic("unreachable")
}

// AltBasalRate is an openaps-compatible form of BasalRate.
type AltBasalRate struct {
	Index   int               `json:"i"`
	Start   string            `json:"start"`
	Minutes int               `json:"minutes"`
	Rate    medtronic.Insulin `json:"rate"`
}

func convertBasalRateSchedule(r medtronic.BasalRateSchedule) interface{} {
	sched := make([]AltBasalRate, len(r))
	for i, s := range r {
		sched[i] = AltBasalRate{
			Index:   i,
			Start:   s.Start.String() + ":00",
			Minutes: minutes(time.Duration(s.Start)),
			Rate:    s.Rate,
		}
	}
	return sched
}

// AltCarbRatio is an openaps-compatible form of CarbRatio.
type AltCarbRatio struct {
	X      int     `json:"x"`
	Index  int     `json:"i"`
	Start  string  `json:"start"`
	Offset int     `json:"offset"`
	Ratio  float64 `json:"ratio"`
}

func convertCarbRatioSchedule(r medtronic.CarbRatioSchedule) interface{} {
	sched := make([]AltCarbRatio, len(r))
	u := r[0].Units
	for i, s := range r {
		if s.Units != u {
			log.Fatalf("unexpected carb units type %v", s.Units)
		}
		m := minutes(time.Duration(s.Start))
		ratio := float64(s.Ratio) / 10
		if u == medtronic.Exchanges {
			ratio /= 100
		}
		sched[i] = AltCarbRatio{
			X:      i,
			Index:  m / 30,
			Start:  s.Start.String() + ":00",
			Offset: m,
			Ratio:  ratio,
		}
	}
	return struct {
		First    int            `json:"first"`
		Units    string         `json:"units"`
		Schedule []AltCarbRatio `json:"schedule"`
	}{
		First:    int(u),
		Units:    carbU(u),
		Schedule: sched,
	}
}

// AltGlucoseTarget is an openaps-compatible form of GlucoseTarget.
type AltGlucoseTarget struct {
	X      int     `json:"x"`
	Index  int     `json:"i"`
	Start  string  `json:"start"`
	Offset int     `json:"offset"`
	Low    float64 `json:"low"`
	High   float64 `json:"high"`
}

func convertGlucoseTargetSchedule(r medtronic.GlucoseTargetSchedule) interface{} {
	sched := make([]AltGlucoseTarget, len(r))
	u := r[0].Units
	for i, s := range r {
		if s.Units != u {
			log.Fatalf("unexpected glucose units type %v", s.Units)
		}
		m := minutes(time.Duration(s.Start))
		low := float64(s.Low)
		high := float64(s.High)
		if u == medtronic.MMolPerLiter {
			low /= 1000
			high /= 1000
		}
		sched[i] = AltGlucoseTarget{
			X:      i,
			Index:  m / 30,
			Start:  s.Start.String() + ":00",
			Offset: m,
			Low:    low,
			High:   high,
		}
	}
	return struct {
		First    int                `json:"first"`
		Units    string             `json:"units"`
		Schedule []AltGlucoseTarget `json:"targets"`
	}{
		First:    int(u),
		Units:    glucoseU(u),
		Schedule: sched,
	}
}

// AltInsulinSensitivity is an openaps-compatible form of InsulinSensitivity.
type AltInsulinSensitivity struct {
	X           int     `json:"x"`
	Index       int     `json:"i"`
	Start       string  `json:"start"`
	Offset      int     `json:"offset"`
	Sensitivity float64 `json:"sensitivity"`
}

func convertInsulinSensitivitySchedule(r medtronic.InsulinSensitivitySchedule) interface{} {
	sched := make([]AltInsulinSensitivity, len(r))
	u := r[0].Units
	for i, s := range r {
		if s.Units != u {
			log.Fatalf("unexpected carb units type %v", s.Units)
		}
		m := minutes(time.Duration(s.Start))
		sens := float64(s.Sensitivity)
		if u == medtronic.MMolPerLiter {
			sens /= 1000
		}
		sched[i] = AltInsulinSensitivity{
			X:           i,
			Index:       m / 30,
			Start:       s.Start.String() + ":00",
			Offset:      m,
			Sensitivity: sens,
		}
	}
	return struct {
		First    int                     `json:"first"`
		Units    string                  `json:"units"`
		Schedule []AltInsulinSensitivity `json:"sensitivities"`
	}{
		First:    int(u),
		Units:    glucoseU(u),
		Schedule: sched,
	}
}

func convertBatteryInfo(r medtronic.BatteryInfo) interface{} {
	status := "normal"
	if r.LowBattery {
		status = "low"
	}
	return struct {
		Voltage medtronic.Voltage `json:"voltage"`
		Status  string            `json:"string"`
	}{
		Voltage: r.Voltage,
		Status:  status,
	}
}

func convertSettingsInfo(r medtronic.SettingsInfo) interface{} {
	return struct {
		AutoOffHours         int               `json:"auto_off_duration_hrs"`
		InsulinActionHours   int               `json:"insulin_action_curve"`
		InsulinConcentration int               `json:"insulinConcentration"`
		MaxBolus             medtronic.Insulin `json:"maxBolus"`
		MaxBasal             medtronic.Insulin `json:"maxBasal"`
		RFEnabled            bool              `json:"rf_enable"`
		SelectedPattern      int               `json:"selected_pattern"`
	}{
		AutoOffHours:         hours(r.AutoOff),
		InsulinActionHours:   hours(r.InsulinAction),
		InsulinConcentration: r.InsulinConcentration,
		MaxBolus:             r.MaxBolus,
		MaxBasal:             r.MaxBasal,
		RFEnabled:            r.RFEnabled,
		SelectedPattern:      r.SelectedPattern,
	}
}

func convertStatusInfo(r medtronic.StatusInfo) interface{} {
	status := "normal"
	if r.Code != 0x03 {
		status = "error"
	}
	return struct {
		Status    string `json:"status"`
		Bolusing  bool   `json:"bolusing"`
		Suspended bool   `json:"suspended"`
	}{
		Status:    status,
		Bolusing:  r.Bolusing,
		Suspended: r.Suspended,
	}
}

func convertTempBasalInfo(r medtronic.TempBasalInfo) interface{} {
	t := struct {
		Duration int         `json:"duration"`
		Temp     string      `json:"temp"`
		Rate     interface{} `json:"rate"`
	}{
		Duration: minutes(r.Duration),
		Temp:     strings.ToLower(r.Type.String()),
	}
	if r.Rate != nil {
		t.Rate = *r.Rate
	} else if r.Percent != nil {
		t.Rate = *r.Percent
	}
	return t
}
