package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ecc1/medtronic"
)

func showJSON(v interface{}) {
	result := OpenapsJSON(v)
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println(err)
		fmt.Println(result)
		return
	}
	fmt.Println(string(b))
}

// OpenapsJSON converts v into a value that the json package
// will marshal in a format compatible with openaps.
func OpenapsJSON(v interface{}) interface{} {
	switch r := v.(type) {
	case medtronic.BasalRateSchedule:
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
	case medtronic.BatteryInfo:
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
	case medtronic.CarbRatioSchedule:
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
				Index:  m / 30,
				Start:  s.Start.String() + ":00",
				Offset: m,
				Ratio:  ratio,
			}
		}
		return struct {
			Units    string         `json:"units"`
			Schedule []AltCarbRatio `json:"schedule"`
		}{
			Units:    carbU(u),
			Schedule: sched,
		}
	case medtronic.CarbUnitsType:
		return carbU(r)
	case medtronic.GlucoseTargetSchedule:
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
				Index:  m / 30,
				Start:  s.Start.String() + ":00",
				Offset: m,
				Low:    low,
				High:   high,
			}
		}
		return struct {
			Units    string             `json:"units"`
			Schedule []AltGlucoseTarget `json:"targets"`
		}{
			Units:    glucoseU(u),
			Schedule: sched,
		}
	case medtronic.GlucoseUnitsType:
		return glucoseU(r)
	case medtronic.InsulinSensitivitySchedule:
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
				Index:       m / 30,
				Start:       s.Start.String() + ":00",
				Offset:      m,
				Sensitivity: sens,
			}
		}
		return struct {
			Units    string                  `json:"units"`
			Schedule []AltInsulinSensitivity `json:"sensitivities"`
		}{
			Units:    glucoseU(u),
			Schedule: sched,
		}
	case medtronic.SettingsInfo:
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
	case medtronic.StatusInfo:
		status := "normal"
		if r.Code != 0x03 {
			status = "error"
		}
		return struct {
			Status    string `json:"string"`
			Bolusing  bool   `json:"bolusing"`
			Suspended bool   `json:"suspended"`
		}{
			Status:    status,
			Bolusing:  r.Bolusing,
			Suspended: r.Suspended,
		}
	case medtronic.TempBasalInfo:
		if r.Type != medtronic.Absolute {
			log.Fatalf("temp basal type is %v", r.Type)
		}
		return struct {
			Duration int               `json:"duration"`
			Rate     medtronic.Insulin `json:"rate"`
			Temp     string            `json:"temp"`
		}{
			Duration: minutes(r.Duration),
			Rate:     *r.Rate,
			Temp:     "absolute",
		}
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

// AltCarbRatio is an openaps-compatible form of CarbRatio.
type AltCarbRatio struct {
	Index  int     `json:"i"`
	Start  string  `json:"start"`
	Offset int     `json:"offset"`
	Ratio  float64 `json:"ratio"`
}

// AltGlucoseTarget is an openaps-compatible form of GlucoseTarget.
type AltGlucoseTarget struct {
	Index  int     `json:"i"`
	Start  string  `json:"start"`
	Offset int     `json:"offset"`
	Low    float64 `json:"low"`
	High   float64 `json:"high"`
}

// AltInsulinSensitivity is an openaps-compatible form of InsulinSensitivity.
type AltInsulinSensitivity struct {
	Index       int     `json:"i"`
	Start       string  `json:"start"`
	Offset      int     `json:"offset"`
	Sensitivity float64 `json:"sensitivity"`
}
