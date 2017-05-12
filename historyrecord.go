package medtronic

import (
	"fmt"
)

// HistoryRecordType represents a history record type.
type HistoryRecordType byte

//go:generate stringer -type HistoryRecordType

// Events stored in the pump's history pages.
const (
	Bolus                   HistoryRecordType = 0x01
	Prime                   HistoryRecordType = 0x03
	Alarm                   HistoryRecordType = 0x06
	DailyTotal              HistoryRecordType = 0x07
	BasalProfileBefore      HistoryRecordType = 0x08
	BasalProfileAfter       HistoryRecordType = 0x09
	BGCapture               HistoryRecordType = 0x0A
	SensorAlarm             HistoryRecordType = 0x0B
	ClearAlarm              HistoryRecordType = 0x0C
	ChangeBasalPattern      HistoryRecordType = 0x14
	TempBasalDuration       HistoryRecordType = 0x16
	ChangeTime              HistoryRecordType = 0x17
	NewTime                 HistoryRecordType = 0x18
	LowBattery              HistoryRecordType = 0x19
	BatteryChange           HistoryRecordType = 0x1A
	SetAutoOff              HistoryRecordType = 0x1B
	SuspendPump             HistoryRecordType = 0x1E
	ResumePump              HistoryRecordType = 0x1F
	SelfTest                HistoryRecordType = 0x20
	Rewind                  HistoryRecordType = 0x21
	ClearSettings           HistoryRecordType = 0x22
	EnableChildBlock        HistoryRecordType = 0x23
	MaxBolus                HistoryRecordType = 0x24
	EnableRemote            HistoryRecordType = 0x26
	MaxBasal                HistoryRecordType = 0x2C
	EnableBolusWizard       HistoryRecordType = 0x2D
	ChangeBGReminder        HistoryRecordType = 0x31
	SetAlarmClockTime       HistoryRecordType = 0x32
	TempBasalRate           HistoryRecordType = 0x33
	LowReservoir            HistoryRecordType = 0x34
	AlarmClock              HistoryRecordType = 0x35
	ChangeMeterID           HistoryRecordType = 0x36
	SensorStatus            HistoryRecordType = 0x3B
	EnableMeter             HistoryRecordType = 0x3C
	BGReceived              HistoryRecordType = 0x3F
	MealMarker              HistoryRecordType = 0x40
	ExerciseMarker          HistoryRecordType = 0x41
	InsulinMarker           HistoryRecordType = 0x42
	OtherMarker             HistoryRecordType = 0x43
	ChangeBolusWizardSetup  HistoryRecordType = 0x4F
	ChangeGlucoseUnits      HistoryRecordType = 0x56
	BolusWizardSetup        HistoryRecordType = 0x5A
	BolusWizard             HistoryRecordType = 0x5B
	UnabsorbedInsulin       HistoryRecordType = 0x5C
	SaveSettings            HistoryRecordType = 0x5D
	EnableVariableBolus     HistoryRecordType = 0x5E
	ChangeEasyBolus         HistoryRecordType = 0x5F
	EnableBGReminder        HistoryRecordType = 0x60
	EnableAlarmClock        HistoryRecordType = 0x61
	ChangeTempBasalType     HistoryRecordType = 0x62
	ChangeAlarmType         HistoryRecordType = 0x63
	ChangeTimeFormat        HistoryRecordType = 0x64
	ChangeReservoirWarning  HistoryRecordType = 0x65
	EnableBolusReminder     HistoryRecordType = 0x66
	SetBolusReminderTime    HistoryRecordType = 0x67
	DeleteBolusReminderTime HistoryRecordType = 0x68
	BolusReminder           HistoryRecordType = 0x69
	DeleteAlarmClockTime    HistoryRecordType = 0x6A
	DailyTotal515           HistoryRecordType = 0x6C
	DailyTotal522           HistoryRecordType = 0x6D
	DailyTotal523           HistoryRecordType = 0x6E
	ChangeCarbUnits         HistoryRecordType = 0x6F
	BasalProfileStart       HistoryRecordType = 0x7B
	ConnectOtherDevices     HistoryRecordType = 0x7C
	ChangeOtherDevice       HistoryRecordType = 0x7D
	ChangeMarriage          HistoryRecordType = 0x81
	DeleteOtherDevice       HistoryRecordType = 0x82
	EnableCaptureEvent      HistoryRecordType = 0x83
)

var decode = map[HistoryRecordType]decoder{
	Bolus:                   decodeBolus,
	Prime:                   decodePrime,
	Alarm:                   decodeAlarm,
	DailyTotal:              decodeDailyTotal,
	BasalProfileBefore:      decodeBasalProfileBefore,
	BasalProfileAfter:       decodeBasalProfileAfter,
	BGCapture:               decodeBGCapture,
	SensorAlarm:             decodeSensorAlarm,
	ClearAlarm:              decodeClearAlarm,
	ChangeBasalPattern:      decodeChangeBasalPattern,
	TempBasalDuration:       decodeTempBasalDuration,
	ChangeTime:              decodeChangeTime,
	NewTime:                 decodeNewTime,
	LowBattery:              decodeLowBattery,
	BatteryChange:           decodeBatteryChange,
	SetAutoOff:              decodeSetAutoOff,
	SuspendPump:             decodeSuspendPump,
	ResumePump:              decodeResumePump,
	SelfTest:                decodeSelfTest,
	Rewind:                  decodeRewind,
	EnableChildBlock:        decodeEnableChildBlock,
	MaxBolus:                decodeMaxBolus,
	EnableRemote:            decodeEnableRemote,
	MaxBasal:                decodeMaxBasal,
	EnableBolusWizard:       decodeEnableBolusWizard,
	ChangeBGReminder:        decodeChangeBGReminder,
	SetAlarmClockTime:       decodeSetAlarmClockTime,
	TempBasalRate:           decodeTempBasalRate,
	LowReservoir:            decodeLowReservoir,
	AlarmClock:              decodeAlarmClock,
	SensorStatus:            decodeSensorStatus,
	EnableMeter:             decodeEnableMeter,
	BGReceived:              decodeBGReceived,
	MealMarker:              decodeMealMarker,
	ExerciseMarker:          decodeExerciseMarker,
	InsulinMarker:           decodeInsulinMarker,
	OtherMarker:             decodeOtherMarker,
	ChangeBolusWizardSetup:  decodeChangeBolusWizardSetup,
	ChangeGlucoseUnits:      decodeChangeGlucoseUnits,
	BolusWizardSetup:        decodeBolusWizardSetup,
	BolusWizard:             decodeBolusWizard,
	UnabsorbedInsulin:       decodeUnabsorbedInsulin,
	EnableVariableBolus:     decodeEnableVariableBolus,
	ChangeEasyBolus:         decodeChangeEasyBolus,
	EnableBGReminder:        decodeEnableBGReminder,
	EnableAlarmClock:        decodeEnableAlarmClock,
	ChangeTempBasalType:     decodeChangeTempBasalType,
	ChangeAlarmType:         decodeChangeAlarmType,
	ChangeTimeFormat:        decodeChangeTimeFormat,
	ChangeReservoirWarning:  decodeChangeReservoirWarning,
	EnableBolusReminder:     decodeEnableBolusReminder,
	SetBolusReminderTime:    decodeSetBolusReminderTime,
	DeleteBolusReminderTime: decodeDeleteBolusReminderTime,
	DeleteAlarmClockTime:    decodeDeleteAlarmClockTime,
	DailyTotal515:           decodeDailyTotal515,
	DailyTotal522:           decodeDailyTotal522,
	DailyTotal523:           decodeDailyTotal523,
	ChangeCarbUnits:         decodeChangeCarbUnits,
	BasalProfileStart:       decodeBasalProfileStart,
	ConnectOtherDevices:     decodeConnectOtherDevices,
	ChangeOtherDevice:       decodeChangeOtherDevice,
	ChangeMarriage:          decodeChangeMarriage,
	DeleteOtherDevice:       decodeDeleteOtherDevice,
	EnableCaptureEvent:      decodeEnableCaptureEvent,
}

// nolint
type (
	decoder func([]byte, bool) HistoryRecord

	HistoryRecord struct {
		Data              []byte                   `json:",omitempty"`
		Time              Time                     `json:",omitempty"`
		Duration          *Duration                `json:",omitempty"`
		Enabled           *bool                    `json:",omitempty"`
		Glucose           *Glucose                 `json:",omitempty"`
		GlucoseUnits      *GlucoseUnitsType        `json:",omitempty"`
		Insulin           *Insulin                 `json:",omitempty"`
		Carbs             *Carbs                   `json:",omitempty"`
		CarbUnits         *CarbUnitsType           `json:",omitempty"`
		TempBasalType     *TempBasalType           `json:",omitempty"`
		Value             *int                     `json:",omitempty"`
		BasalProfile      BasalRateSchedule        `json:",omitempty"`
		BasalProfileStart *BasalProfileStartRecord `json:",omitempty"`
		Bolus             *BolusRecord             `json:",omitempty"`
		BolusWizard       *BolusWizardRecord       `json:",omitempty"`
		BolusWizardSetup  *BolusWizardSetupRecord  `json:",omitempty"`
		Prime             *PrimeRecord             `json:",omitempty"`
		UnabsorbedInsulin UnabsorbedBolusHistory   `json:",omitempty"`
	}

	BasalProfileStartRecord struct {
		ProfileIndex int
		BasalRate    BasalRate
	}

	PrimeRecord struct {
		Fixed  Insulin
		Manual Insulin
	}

	BolusRecord struct {
		Programmed Insulin
		Amount     Insulin
		Unabsorbed Insulin
		Duration   Duration // non-zero for square wave bolus
	}

	BolusWizardRecord struct {
		GlucoseInput Glucose
		CarbInput    Carbs
		GlucoseUnits GlucoseUnitsType
		CarbUnits    CarbUnitsType
		TargetLow    Glucose
		TargetHigh   Glucose
		Sensitivity  Glucose // glucose reduction per insulin unit
		CarbRatio    Ratio   // 10x grams/unit or 1000x units/exchange
		Correction   Insulin
		Food         Insulin
		Unabsorbed   Insulin
		Bolus        Insulin
	}

	BolusWizardConfig struct {
		Ratios        CarbRatioSchedule
		Sensitivities InsulinSensitivitySchedule
		Targets       GlucoseTargetSchedule
		InsulinAction Duration
	}

	BolusWizardSetupRecord struct {
		Before BolusWizardConfig
		After  BolusWizardConfig
	}

	UnabsorbedBolus struct {
		Bolus Insulin
		Age   Duration
	}

	UnabsorbedBolusHistory []UnabsorbedBolus

	UnknownRecordTypeError struct {
		Data []byte
	}
)

// Type returns the history record type.
func (r HistoryRecord) Type() HistoryRecordType {
	return HistoryRecordType(r.Data[0])
}

func decodeBase(data []byte, newerPump bool) HistoryRecord {
	return HistoryRecord{
		Time: decodeTime(data[2:7]),
		Data: data[:7],
	}
}

func decodeEnable(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	enabled := false
	if data[1] != 0 {
		enabled = true
	}
	r.Enabled = &enabled
	return r
}

func decodeDailyTotalDate(data []byte, newerPump bool) HistoryRecord {
	return HistoryRecord{
		Time: decodeDate(data[1:3]),
		Data: data[:3],
	}
}

func extendDecoder(orig decoder) func(int) decoder {
	return func(length int) decoder {
		return func(data []byte, newerPump bool) HistoryRecord {
			r := orig(data, newerPump)
			r.Data = data[:length]
			return r
		}
	}
}

var decodeBaseN = extendDecoder(decodeBase)

var decodeEnableN = extendDecoder(decodeEnable)

var decodeDailyTotalN = extendDecoder(decodeDailyTotalDate)

func decodeValue(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	v := int(data[1])
	r.Value = &v
	return r
}

func decodeInsulin(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	i := byteToInsulin(data[1], true)
	r.Insulin = &i
	return r
}

func decodeBolus(data []byte, newerPump bool) HistoryRecord {
	switch newerPump {
	case true:
		return HistoryRecord{
			Bolus: &BolusRecord{
				Programmed: twoByteInsulin(data[1:3], true),
				Amount:     twoByteInsulin(data[3:5], true),
				Unabsorbed: twoByteInsulin(data[5:7], true),
				Duration:   halfHoursToDuration(data[7]),
			},
			Time: decodeTime(data[8:13]),
			Data: data[:13],
		}
	case false:
		return HistoryRecord{
			Bolus: &BolusRecord{
				Programmed: byteToInsulin(data[1], false),
				Amount:     byteToInsulin(data[2], false),
				Duration:   halfHoursToDuration(data[3]),
			},
			Time: decodeTime(data[4:9]),
			Data: data[:9],
		}
	}
	panic("unreachable")
}

func decodePrime(data []byte, newerPump bool) HistoryRecord {
	return HistoryRecord{
		Prime: &PrimeRecord{
			Fixed:  byteToInsulin(data[2], false),
			Manual: byteToInsulin(data[4], false),
		},
		Time: decodeTime(data[5:10]),
		Data: data[:10],
	}
}

func decodeAlarm(data []byte, newerPump bool) HistoryRecord {
	r := HistoryRecord{
		Time: decodeTime(data[4:9]),
		Data: data[:9],
	}
	alarm := int(data[1])
	r.Value = &alarm
	return r
}

func decodeDailyTotal(data []byte, newerPump bool) HistoryRecord {
	t := decodeDate(data[5:7])
	total := twoByteInsulin(data[3:5], true)
	switch newerPump {
	case true:
		return HistoryRecord{
			Time:    t,
			Insulin: &total,
			Data:    data[:10],
		}
	case false:
		return HistoryRecord{
			Time:    t,
			Insulin: &total,
			Data:    data[:7],
		}
	}
	panic("unreachable")
}

// Note that this is a different format than the response to BasalRates.
func decodeBasalRate(data []byte) BasalRate {
	return BasalRate{
		Start: halfHoursToTimeOfDay(data[0]),
		Rate:  byteToInsulin(data[1], true),
		// data[2] unused
	}
}

func decodeBasalProfile(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	body := data[7:]
	var sched []BasalRate
	for i := 0; i < 144; i += 3 {
		b := decodeBasalRate(body[i : i+3])
		// Don't stop if the 00:00 rate happens to be zero.
		if i > 0 && b.Start == 0 && b.Rate == 0 {
			break
		}
		sched = append(sched, b)
	}
	r.BasalProfile = sched
	r.Data = data[:152]
	return r
}

var decodeBasalProfileBefore = decodeBasalProfile

var decodeBasalProfileAfter = decodeBasalProfile

func decodeBGCapture(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	bgU := GlucoseUnitsType((data[4] >> 5) & 0x3)
	bg := intToGlucose(int(data[4]>>7)<<9|int(data[6]>>7)<<8|int(data[1]), MgPerDeciLiter)
	r.GlucoseUnits = &bgU
	r.Glucose = &bg
	return r
}

func decodeSensorAlarm(data []byte, newerPump bool) HistoryRecord {
	r := HistoryRecord{
		Time: decodeTime(data[3:8]),
		Data: data[:8],
	}
	alarm := int(data[1])
	r.Value = &alarm
	return r
}

func decodeClearAlarm(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	alarm := int(data[1])
	r.Value = &alarm
	return r
}

var decodeChangeBasalPattern = decodeValue

func decodeTempBasalDuration(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	d := halfHoursToDuration(data[1])
	r.Duration = &d
	return r
}

var decodeChangeTime = decodeBase

var decodeNewTime = decodeBase

var decodeLowBattery = decodeBase

var decodeBatteryChange = decodeBase

func decodeSetAutoOff(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	d := hoursToDuration(data[1])
	r.Duration = &d
	return r
}

var decodeSuspendPump = decodeBase

var decodeResumePump = decodeBase

var decodeSelfTest = decodeBase

var decodeRewind = decodeBase

var decodeEnableChildBlock = decodeEnable

var decodeMaxBolus = decodeInsulin

var decodeEnableRemote = decodeEnableN(21)

var decodeMaxBasal = decodeInsulin

var decodeEnableBolusWizard = decodeEnable

var decodeChangeBGReminder = decodeBase

var decodeSetAlarmClockTime = decodeBase

func decodeTempBasalRate(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	tempBasalType := TempBasalType(data[7] >> 3)
	if tempBasalType == Absolute {
		rate := intToInsulin(int(data[7]&0x7)<<8|int(data[1]), true)
		r.Insulin = &rate
	} else {
		rate := int(data[1])
		r.Value = &rate
	}
	r.TempBasalType = &tempBasalType
	r.Data = data[:8]
	return r
}

func decodeLowReservoir(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	amount := byteToInsulin(data[1], false)
	r.Insulin = &amount
	return r
}

var decodeAlarmClock = decodeBase

var decodeSensorStatus = decodeEnable

var decodeEnableMeter = decodeEnableN(21)

var decodeBGReceived = decodeBaseN(10)

func decodeMealMarker(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	r.Data = data[:9]
	carbs := Carbs(int(data[1])<<8 | int(data[7]))
	r.Carbs = &carbs
	carbU := CarbUnitsType(data[8] & 0x3)
	r.CarbUnits = &carbU
	return r
}

func decodeExerciseMarker(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	r.Data = data[:8]
	return r
}

func decodeInsulinMarker(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	r.Data = data[:8]
	amount := intToInsulin(int(data[4]&0x60)<<3|int(data[1]), false)
	r.Insulin = &amount
	return r
}

var decodeOtherMarker = decodeBase

var decodeChangeBolusWizardSetup = decodeBaseN(39)

var decodeChangeGlucoseUnits = decodeBaseN(12)

func decodeBolusWizardConfig(data []byte, newerPump bool) BolusWizardConfig {
	const numEntries = 8
	r := BolusWizardConfig{}
	carbUnits := CarbUnitsType(data[0] & 0x3)
	bgUnits := GlucoseUnitsType((data[0] >> 2) & 0x3)
	step := carbRatioStep(newerPump)
	r.Ratios = decodeCarbRatioSchedule(data[2:2+numEntries*step], carbUnits, newerPump)
	data = data[2+numEntries*step:]
	r.Sensitivities = decodeInsulinSensitivitySchedule(data[:numEntries*2], bgUnits)
	switch newerPump {
	case true:
		data = data[numEntries*2+2:]
	case false:
		data = data[numEntries*2:]
	}
	r.Targets = decodeGlucoseTargetSchedule(data[:numEntries*3], bgUnits)
	return r
}

func decodeBolusWizardSetup(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	switch newerPump {
	case true:
		r.Data = data[:144]
	case false:
		r.Data = data[:124]
	}
	n := len(r.Data) - 1
	body := data[7:n]
	half := (n - 7) / 2
	setup := &BolusWizardSetupRecord{
		Before: decodeBolusWizardConfig(body[:half], newerPump),
		After:  decodeBolusWizardConfig(body[half:], newerPump),
	}
	setup.Before.InsulinAction = hoursToDuration(data[n] & 0xF)
	setup.After.InsulinAction = hoursToDuration(data[n] >> 4)
	r.BolusWizardSetup = setup
	return r
}

func decodeBolusWizard(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	bg := int(data[1])
	body := data[7:]
	bgU := GlucoseUnitsType((body[1] >> 6) & 0x3)
	carbU := CarbUnitsType((body[1] >> 4) & 0x3)
	switch newerPump {
	case true:
		r.BolusWizard = &BolusWizardRecord{
			GlucoseInput: intToGlucose(bg|int(body[1]&0x3)<<8, bgU),
			CarbInput:    Carbs(int(body[1]&0xC)<<6 | int(body[0])),
			GlucoseUnits: bgU,
			CarbUnits:    carbU,
			TargetLow:    byteToGlucose(body[5], bgU),
			TargetHigh:   byteToGlucose(body[14], bgU),
			Sensitivity:  byteToGlucose(body[4], bgU),
			CarbRatio:    intToRatio(int(body[2]&0xF)<<8|int(body[3]), carbU, true),
			Correction:   intToInsulin(int(body[9]&0x38)<<5|int(body[6]), true),
			Food:         twoByteInsulin(body[7:9], true),
			Unabsorbed:   twoByteInsulin(body[10:12], true),
			Bolus:        twoByteInsulin(body[12:14], true),
		}
		r.Data = data[:22]
	case false:
		r.BolusWizard = &BolusWizardRecord{
			GlucoseInput: intToGlucose(bg|int(body[1]&0xF)<<8, bgU),
			CarbInput:    Carbs(body[0]),
			GlucoseUnits: bgU,
			CarbUnits:    carbU,
			TargetLow:    byteToGlucose(body[4], bgU),
			TargetHigh:   byteToGlucose(body[12], bgU),
			Sensitivity:  byteToGlucose(body[3], bgU),
			CarbRatio:    intToRatio(int(body[2]), carbU, false),
			Correction:   intToInsulin(int(body[7])+int(body[5]&0xF), false),
			Food:         byteToInsulin(body[6], false),
			Unabsorbed:   byteToInsulin(body[9], false),
			Bolus:        byteToInsulin(body[11], false),
		}
		r.Data = data[:20]
	}
	return r
}

func decodeUnabsorbedInsulin(data []byte, newerPump bool) HistoryRecord {
	n := int(data[1]) - 2
	body := data[2:]
	var unabsorbed []UnabsorbedBolus
	for i := 0; i < n; i += 3 {
		amount := byteToInsulin(body[i], true)
		curve := body[i+2]
		age := minutesToDuration(body[i+1] + (curve&0x30)<<4)
		unabsorbed = append(unabsorbed, UnabsorbedBolus{
			Bolus: amount,
			Age:   age,
		})
	}
	return HistoryRecord{
		Data:              data[:n+2],
		UnabsorbedInsulin: unabsorbed,
	}
}

var decodeEnableVariableBolus = decodeEnable

var decodeChangeEasyBolus = decodeBase

var decodeEnableBGReminder = decodeEnable

var decodeEnableAlarmClock = decodeEnable

func decodeChangeTempBasalType(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	tempBasalType := TempBasalType(data[1])
	r.TempBasalType = &tempBasalType
	return r
}

var decodeChangeAlarmType = decodeValue

var decodeChangeTimeFormat = decodeValue

func decodeChangeReservoirWarning(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	v := data[1]
	if v&0x1 == 0 {
		amount := Insulin(1000 * int(v>>2))
		r.Insulin = &amount
	} else {
		d := halfHoursToDuration(v >> 2)
		r.Duration = &d
	}
	return r
}

var decodeEnableBolusReminder = decodeEnable

var decodeSetBolusReminderTime = decodeEnableN(9)

var decodeDeleteBolusReminderTime = decodeEnableN(9)

var decodeDeleteAlarmClockTime = decodeBase

var decodeDailyTotal515 = decodeDailyTotalN(38)

var decodeDailyTotal522 = decodeDailyTotalN(44)

var decodeDailyTotal523 = decodeDailyTotalN(52)

var decodeChangeCarbUnits = decodeValue

func decodeBasalProfileStart(data []byte, newerPump bool) HistoryRecord {
	r := decodeBase(data, newerPump)
	r.BasalProfileStart = &BasalProfileStartRecord{
		ProfileIndex: int(data[1]),
		BasalRate:    decodeBasalRate(data[7:10]),
	}
	r.Data = data[:10]
	return r
}

var decodeConnectOtherDevices = decodeEnable

var decodeChangeOtherDevice = decodeBaseN(37)

var decodeChangeMarriage = decodeBaseN(12)

var decodeDeleteOtherDevice = decodeBaseN(12)

var decodeEnableCaptureEvent = decodeEnable

// DecodeHistoryRecord decodes a history record based on its type.
func DecodeHistoryRecord(data []byte, newerPump bool) (HistoryRecord, error) {
	if len(data) == 0 {
		return HistoryRecord{}, fmt.Errorf("empty history record")
	}
	decoder := decode[HistoryRecordType(data[0])]
	if decoder == nil {
		return HistoryRecord{}, unknownRecord(data)
	}
	return decoder(data, newerPump), nil
}

func (e UnknownRecordTypeError) Error() string {
	return fmt.Sprintf("unknown record type here: % X", e.Data)
}

func unknownRecord(data []byte) error {
	return UnknownRecordTypeError{
		Data: data,
	}
}

// DecodeHistoryRecords decodes the records in a page of data
// and returns them in reverse chronological order (most recent first),
// to match the order of the history pages themselves.
func DecodeHistoryRecords(data []byte, newerPump bool) ([]HistoryRecord, error) {
	var results []HistoryRecord
	r := HistoryRecord{}
	err := error(nil)
	for !allZero(data) {
		r, err = DecodeHistoryRecord(data, newerPump)
		if err != nil {
			break
		}
		results = append(results, r)
		data = data[len(r.Data):]
	}
	ReverseHistoryRecords(results)
	return results, err
}

// Partially filled history pages are padded to the end with zero bytes.
func allZero(data []byte) bool {
	for _, b := range data {
		if b != 0 {
			return false
		}
	}
	return true
}

// ReverseHistoryRecords reverses a slice of history records.
func ReverseHistoryRecords(a []HistoryRecord) {
	for i, j := 0, len(a)-1; i < len(a)/2; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}
