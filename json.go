package medtronic

import (
	"encoding/json"
	"fmt"
	"time"
)

func (r HistoryRecord) MarshalJSON() ([]byte, error) {
	t := ""
	if !r.Time.IsZero() {
		t = r.Time.Format(TimeLayout)
	}
	type Original HistoryRecord
	rep := struct {
		Type     string
		Time     string `json:",omitempty"`
		Duration string `json:",omitempty"`
		Original
	}{
		Type:     r.Type().String(),
		Time:     t,
		Original: Original(r),
	}
	if r.Duration != nil {
		rep.Duration = r.Duration.String()
	}
	return json.Marshal(rep)
}

func (r *HistoryRecord) UnmarshalJSON(data []byte) error {
	type Original HistoryRecord
	rep := struct {
		Type     string
		Time     string `json:",omitempty"`
		Duration string `json:",omitempty"`
		*Original
	}{
		Original: (*Original)(r),
	}
	err := json.Unmarshal(data, &rep)
	if err != nil {
		return err
	}
	if rep.Time != "" {
		r.Time, err = time.Parse(TimeLayout, rep.Time)
		if err != nil {
			return err
		}
	}
	if rep.Duration != "" {
		d, err := time.ParseDuration(rep.Duration)
		r.Duration = &d
		if err != nil {
			return err
		}
	}
	return nil
}

func (r BasalRate) MarshalJSON() ([]byte, error) {
	type Original BasalRate
	rep := struct {
		Start string
		Original
	}{
		Start:    r.Start.String(),
		Original: Original(r),
	}
	return json.Marshal(rep)
}

func (r *BasalRate) UnmarshalJSON(data []byte) error {
	type Original BasalRate
	rep := struct {
		Start string
		*Original
	}{
		Original: (*Original)(r),
	}
	err := json.Unmarshal(data, &rep)
	if err != nil {
		return err
	}
	r.Start, err = parseTimeOfDay(rep.Start)
	return err
}

func (r BolusRecord) MarshalJSON() ([]byte, error) {
	type Original BolusRecord
	rep := struct {
		Duration string
		Original
	}{
		Duration: r.Duration.String(),
		Original: Original(r),
	}
	return json.Marshal(rep)
}

func (r *BolusRecord) UnmarshalJSON(data []byte) error {
	type Original BolusRecord
	rep := struct {
		Duration string
		*Original
	}{
		Original: (*Original)(r),
	}
	err := json.Unmarshal(data, &rep)
	if err != nil {
		return err
	}
	r.Duration, err = time.ParseDuration(rep.Duration)
	return err
}

func (r BolusWizardConfig) MarshalJSON() ([]byte, error) {
	type Original BolusWizardConfig
	rep := struct {
		InsulinAction string
		Original
	}{
		InsulinAction: r.InsulinAction.String(),
		Original:      Original(r),
	}
	return json.Marshal(rep)
}

func (r *BolusWizardConfig) UnmarshalJSON(data []byte) error {
	type Original BolusWizardConfig
	rep := struct {
		InsulinAction string
		*Original
	}{
		Original: (*Original)(r),
	}
	err := json.Unmarshal(data, &rep)
	if err != nil {
		return err
	}
	r.InsulinAction, err = time.ParseDuration(rep.InsulinAction)
	return err
}

func (r Tenths) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(r) / 10)
}

func (r *Tenths) UnmarshalJSON(data []byte) error {
	v := 0.0
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*r = Tenths(10*v + 0.5)
	return nil
}

func (r CarbRatio) MarshalJSON() ([]byte, error) {
	type Original CarbRatio
	rep := struct {
		Start string
		Original
	}{
		Start:    r.Start.String(),
		Original: Original(r),
	}
	return json.Marshal(rep)
}

func (r *CarbRatio) UnmarshalJSON(data []byte) error {
	type Original CarbRatio
	rep := struct {
		Start string
		*Original
	}{
		Original: (*Original)(r),
	}
	err := json.Unmarshal(data, &rep)
	if err != nil {
		return err
	}
	r.Start, err = parseTimeOfDay(rep.Start)
	return err
}

func (r GlucoseTarget) MarshalJSON() ([]byte, error) {
	type Original GlucoseTarget
	rep := struct {
		Start string
		Original
	}{
		Start:    r.Start.String(),
		Original: Original(r),
	}
	return json.Marshal(rep)
}

func (r *GlucoseTarget) UnmarshalJSON(data []byte) error {
	type Original GlucoseTarget
	rep := struct {
		Start string
		*Original
	}{
		Original: (*Original)(r),
	}
	err := json.Unmarshal(data, &rep)
	if err != nil {
		return err
	}
	r.Start, err = parseTimeOfDay(rep.Start)
	return err
}

func (r InsulinSensitivity) MarshalJSON() ([]byte, error) {
	type Original InsulinSensitivity
	rep := struct {
		Start string
		Original
	}{
		Start:    r.Start.String(),
		Original: Original(r),
	}
	return json.Marshal(rep)
}

func (r *InsulinSensitivity) UnmarshalJSON(data []byte) error {
	type Original InsulinSensitivity
	rep := struct {
		Start string
		*Original
	}{
		Original: (*Original)(r),
	}
	err := json.Unmarshal(data, &rep)
	if err != nil {
		return err
	}
	r.Start, err = parseTimeOfDay(rep.Start)
	return err
}

func (r UnabsorbedBolus) MarshalJSON() ([]byte, error) {
	type Original UnabsorbedBolus
	rep := struct {
		Age string
		Original
	}{
		Age:      r.Age.String(),
		Original: Original(r),
	}
	return json.Marshal(rep)
}

func (r *UnabsorbedBolus) UnmarshalJSON(data []byte) error {
	type Original UnabsorbedBolus
	rep := struct {
		Age string
		*Original
	}{
		Original: (*Original)(r),
	}
	err := json.Unmarshal(data, &rep)
	if err != nil {
		return err
	}
	r.Age, err = time.ParseDuration(rep.Age)
	return err
}

func (r SettingsInfo) MarshalJSON() ([]byte, error) {
	type Original SettingsInfo
	rep := struct {
		AutoOff       string
		InsulinAction string
		Original
	}{
		AutoOff:       r.AutoOff.String(),
		InsulinAction: r.InsulinAction.String(),
		Original:      Original(r),
	}
	return json.Marshal(rep)
}

func (r *SettingsInfo) UnmarshalJSON(data []byte) error {
	type Original SettingsInfo
	rep := struct {
		AutoOff       string
		InsulinAction string
		*Original
	}{
		Original: (*Original)(r),
	}
	err := json.Unmarshal(data, &rep)
	if err != nil {
		return err
	}
	r.AutoOff, err = time.ParseDuration(rep.AutoOff)
	if err != nil {
		return err
	}
	r.InsulinAction, err = time.ParseDuration(rep.InsulinAction)
	return err
}

func (r TempBasalInfo) MarshalJSON() ([]byte, error) {
	type Original TempBasalInfo
	rep := struct {
		Duration string
		Original
	}{
		Duration: r.Duration.String(),
		Original: Original(r),
	}
	return json.Marshal(rep)
}

func (r *TempBasalInfo) UnmarshalJSON(data []byte) error {
	type Original TempBasalInfo
	rep := struct {
		Duration string
		*Original
	}{
		Original: (*Original)(r),
	}
	err := json.Unmarshal(data, &rep)
	if err != nil {
		return err
	}
	r.Duration, err = time.ParseDuration(rep.Duration)
	return err
}

func (r CarbUnitsType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, r)), nil
}

func (r *CarbUnitsType) UnmarshalJSON(data []byte) error {
	err := error(nil)
	switch string(data) {
	case `"Grams"`:
		*r = Grams
	case `"Exchanges"`:
		*r = Exchanges
	default:
		err = fmt.Errorf("unknown CarbUnitsType (%s)", data)
	}
	return err
}

func (r GlucoseUnitsType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, r)), nil
}

func (r *GlucoseUnitsType) UnmarshalJSON(data []byte) error {
	err := error(nil)
	switch string(data) {
	case `"mg/dL"`:
		*r = MgPerDeciLiter
	case `"μmol/L"`:
		*r = MicromolPerLiter
	default:
		err = fmt.Errorf("unknown GlucoseUnitsType (%s)", data)
	}
	return err
}

func (r TempBasalType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, r)), nil
}

func (r *TempBasalType) UnmarshalJSON(data []byte) error {
	err := error(nil)
	switch string(data) {
	case `"Absolute"`:
		*r = Absolute
	case `"Percent"`:
		*r = Percent
	default:
		err = fmt.Errorf("unknown TempBasalType (%s)", data)
	}
	return err
}

func (r Insulin) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(r) / 1000)
}

func (r *Insulin) UnmarshalJSON(data []byte) error {
	v := 0.0
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*r = Insulin(1000*v + 0.5)
	return nil
}

func (r Voltage) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(r) / 1000)
}

func (r *Voltage) UnmarshalJSON(data []byte) error {
	v := 0.0
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*r = Voltage(1000*v + 0.5)
	return nil
}
