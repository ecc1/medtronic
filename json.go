package medtronic

import (
	"encoding/json"
	"fmt"
	"time"
)

func (r HistoryRecord) MarshalJSON() ([]byte, error) {
	type Original HistoryRecord
	rep := struct {
		Type string
		Time string `json:",omitempty"`
		Original
	}{
		Type:     fmt.Sprintf("%v", r.Type()),
		Original: Original(r),
	}
	t := time.Time(r.Time)
	if !t.IsZero() {
		rep.Time = t.Format(JsonTimeLayout)
	}
	return json.Marshal(rep)
}

func (r *HistoryRecord) UnmarshalJSON(data []byte) error {
	type Original HistoryRecord
	rep := struct {
		Type string
		Time string
		*Original
	}{
		Original: (*Original)(r),
	}
	err := json.Unmarshal(data, &rep)
	if err != nil {
		return err
	}
	if rep.Time != "" {
		var t time.Time
		t, err = time.Parse(JsonTimeLayout, rep.Time)
		r.Time = Time(t)
	}
	return err
}

func (r BolusWizardRecord) MarshalJSON() ([]byte, error) {
	type Original BolusWizardRecord
	rep := struct {
		CarbRatio float64
		Original
	}{
		Original: Original(r),
	}
	switch r.CarbUnits {
	case Grams:
		rep.CarbRatio = float64(r.CarbRatio) / 10
	case Exchanges:
		rep.CarbRatio = float64(r.CarbRatio) / 1000
	default:
		return nil, fmt.Errorf("unknown carb unit %d marshaling BolusWizardRecord", r.CarbUnits)
	}
	return json.Marshal(rep)
}

func (r *BolusWizardRecord) UnmarshalJSON(data []byte) error {
	type Original BolusWizardRecord
	rep := struct {
		CarbRatio float64
		*Original
	}{
		Original: (*Original)(r),
	}
	err := json.Unmarshal(data, &rep)
	if err != nil {
		return err
	}
	switch r.CarbUnits {
	case Grams:
		r.CarbRatio = Ratio(10*rep.CarbRatio + 0.5)
	case Exchanges:
		r.CarbRatio = Ratio(1000*rep.CarbRatio + 0.5)
	default:
		err = fmt.Errorf("unknown carb unit %d unmarshaling BolusWizardRecord", r.CarbUnits)
	}
	return err
}

func (r CarbRatio) MarshalJSON() ([]byte, error) {
	type Original CarbRatio
	rep := struct {
		Ratio float64
		Original
	}{
		Original: Original(r),
	}
	switch r.Units {
	case Grams:
		rep.Ratio = float64(r.Ratio) / 10
	case Exchanges:
		rep.Ratio = float64(r.Ratio) / 1000
	default:
		return nil, fmt.Errorf("unknown carb unit %d marshaling CarbRatio", r.Units)
	}
	return json.Marshal(rep)
}

func (r *CarbRatio) UnmarshalJSON(data []byte) error {
	type Original CarbRatio
	rep := struct {
		Ratio float64
		*Original
	}{
		Original: (*Original)(r),
	}
	err := json.Unmarshal(data, &rep)
	if err != nil {
		return err
	}
	switch r.Units {
	case Grams:
		r.Ratio = Ratio(10*rep.Ratio + 0.5)
	case Exchanges:
		r.Ratio = Ratio(1000*rep.Ratio + 0.5)
	default:
		err = fmt.Errorf("unknown carb unit %d unmarshaling CarbRatio", r.Units)
	}
	return err
}

func (r Ratio) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("cannot marshal carb ratio without units")
}

func (r *Ratio) UnmarshalJSON([]byte) error {
	return fmt.Errorf("cannot unmarshal carb ratio without units")
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
		*r = MMolPerLiter
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

func (r Time) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("marshaling Time value")
}

func (r *Time) UnmarshalJSON(data []byte) error {
	return fmt.Errorf("unmarshaling Time value")
}

func (r Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(r).String())
}

func (r *Duration) UnmarshalJSON(data []byte) error {
	v := ""
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	d, err := time.ParseDuration(v)
	*r = Duration(d)
	return err
}

func (r TimeOfDay) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *TimeOfDay) UnmarshalJSON(data []byte) error {
	v := ""
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*r, err = parseTimeOfDay(v)
	return err
}
