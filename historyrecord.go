package medtronic

import (
	"fmt"
	"time"
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
	PrepareInsulinChange    HistoryRecordType = 0x1C
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
	Unknown2E               HistoryRecordType = 0x2E
	BolusWizard512          HistoryRecordType = 0x2F
	UnabsorbedInsulin512    HistoryRecordType = 0x30
	ChangeBGReminder        HistoryRecordType = 0x31
	SetAlarmClockTime       HistoryRecordType = 0x32
	TempBasalRate           HistoryRecordType = 0x33
	LowReservoir            HistoryRecordType = 0x34
	AlarmClock              HistoryRecordType = 0x35
	ChangeMeterID           HistoryRecordType = 0x36
	BGReceived512           HistoryRecordType = 0x39
	ConfirmInsulinChange    HistoryRecordType = 0x3A
	SensorStatus            HistoryRecordType = 0x3B
	EnableMeter             HistoryRecordType = 0x3C
	BGReceived              HistoryRecordType = 0x3F
	MealMarker              HistoryRecordType = 0x40
	ExerciseMarker          HistoryRecordType = 0x41
	InsulinMarker           HistoryRecordType = 0x42
	OtherMarker             HistoryRecordType = 0x43
	EnableSensorAutoCal     HistoryRecordType = 0x44
	ChangeBolusWizardSetup  HistoryRecordType = 0x4F
	SensorSetup             HistoryRecordType = 0x50
	Sensor51                HistoryRecordType = 0x51
	Sensor52                HistoryRecordType = 0x52
	ChangeSensorAlarm       HistoryRecordType = 0x53
	Sensor54                HistoryRecordType = 0x54
	Sensor55                HistoryRecordType = 0x55
	ChangeSensorAlert       HistoryRecordType = 0x56
	ChangeBolusStep         HistoryRecordType = 0x57
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
	PrepareInsulinChange:    decodePrepareInsulinChange,
	SuspendPump:             decodeSuspendPump,
	ResumePump:              decodeResumePump,
	SelfTest:                decodeSelfTest,
	Rewind:                  decodeRewind,
	ClearSettings:           decodeClearSettings,
	EnableChildBlock:        decodeEnableChildBlock,
	MaxBolus:                decodeMaxBolus,
	EnableRemote:            decodeEnableRemote,
	MaxBasal:                decodeMaxBasal,
	EnableBolusWizard:       decodeEnableBolusWizard,
	Unknown2E:               decodeUnknown2E,
	BolusWizard512:          decodeBolusWizard512,
	UnabsorbedInsulin512:    decodeUnabsorbedInsulin512,
	ChangeBGReminder:        decodeChangeBGReminder,
	SetAlarmClockTime:       decodeSetAlarmClockTime,
	TempBasalRate:           decodeTempBasalRate,
	LowReservoir:            decodeLowReservoir,
	AlarmClock:              decodeAlarmClock,
	ChangeMeterID:           decodeChangeMeterID,
	BGReceived512:           decodeBGReceived512,
	ConfirmInsulinChange:    decodeConfirmInsulinChange,
	SensorStatus:            decodeSensorStatus,
	EnableMeter:             decodeEnableMeter,
	BGReceived:              decodeBGReceived,
	MealMarker:              decodeMealMarker,
	ExerciseMarker:          decodeExerciseMarker,
	InsulinMarker:           decodeInsulinMarker,
	OtherMarker:             decodeOtherMarker,
	EnableSensorAutoCal:     decodeEnableSensorAutoCal,
	ChangeBolusWizardSetup:  decodeChangeBolusWizardSetup,
	SensorSetup:             decodeSensorSetup,
	Sensor51:                decodeSensor51,
	Sensor52:                decodeSensor52,
	ChangeSensorAlarm:       decodeChangeSensorAlarm,
	Sensor54:                decodeSensor54,
	Sensor55:                decodeSensor55,
	ChangeSensorAlert:       decodeChangeSensorAlert,
	ChangeBolusStep:         decodeChangeBolusStep,
	BolusWizardSetup:        decodeBolusWizardSetup,
	BolusWizard:             decodeBolusWizard,
	UnabsorbedInsulin:       decodeUnabsorbedInsulin,
	SaveSettings:            decodeSaveSettings,
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
	BolusReminder:           decodeBolusReminder,
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

//AlarmCode represents a pump alarm code.
type AlarmCode byte

//go:generate stringer -type AlarmCode

const (
	BatteryOutLimitExceeded AlarmCode = 0x03
	NoDelivery              AlarmCode = 0x04
	BatteryDepleted         AlarmCode = 0x05
	AutoOff                 AlarmCode = 0x06
	DeviceReset             AlarmCode = 0x10
	ReprogramError          AlarmCode = 0x3D
	EmptyReservoir          AlarmCode = 0x3E
)

// nolint
type (
	decoder func([]byte, Family) HistoryRecord

	HistoryRecord struct {
		Data []byte
		Time time.Time
		Info interface{} `json:",omitempty"`
	}

	History []HistoryRecord

	GlucoseRecord struct {
		Units   GlucoseUnitsType
		Glucose Glucose
		MeterID string `json:",omitempty"`
	}

	CarbRecord struct {
		Units CarbUnitsType
		Carbs Carbs
	}

	TempBasalRecord struct {
		Type  TempBasalType
		Value interface{}
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

func decodeBase(data []byte, family Family) HistoryRecord {
	return HistoryRecord{
		Time: decodeTime(data[2:7]),
		Data: data[:7],
	}
}

func decodeEnable(data []byte, family Family) HistoryRecord {
	r := decodeBase(data, family)
	r.Info = data[1] != 0
	return r
}

func decodeDailyTotalDate(data []byte, family Family) HistoryRecord {
	return HistoryRecord{
		Time: decodeDate(data[1:3]),
		Data: data[:3],
	}
}

func extendDecoder(orig decoder) func(int) decoder {
	return func(length int) decoder {
		return func(data []byte, family Family) HistoryRecord {
			r := orig(data, family)
			r.Data = data[:length]
			return r
		}
	}
}

var decodeBaseN = extendDecoder(decodeBase)

var decodeEnableN = extendDecoder(decodeEnable)

var decodeDailyTotalN = extendDecoder(decodeDailyTotalDate)

func decodeValue(data []byte, family Family) HistoryRecord {
	r := decodeBase(data, family)
	r.Info = int(data[1])
	return r
}

func decodeInsulin(data []byte, family Family) HistoryRecord {
	r := decodeBase(data, family)
	r.Info = byteToInsulin(data[1], 23)
	return r
}

func decodeBolus(data []byte, family Family) HistoryRecord {
	if family <= 22 {
		return HistoryRecord{
			Info: BolusRecord{
				Programmed: byteToInsulin(data[1], family),
				Amount:     byteToInsulin(data[2], family),
				Duration:   halfHoursToDuration(data[3]),
			},
			Time: decodeTime(data[4:9]),
			Data: data[:9],
		}
	}
	return HistoryRecord{
		Info: BolusRecord{
			Programmed: twoByteInsulin(data[1:3], family),
			Amount:     twoByteInsulin(data[3:5], family),
			Unabsorbed: twoByteInsulin(data[5:7], family),
			Duration:   halfHoursToDuration(data[7]),
		},
		Time: decodeTime(data[8:13]),
		Data: data[:13],
	}
}

func decodePrime(data []byte, family Family) HistoryRecord {
	return HistoryRecord{
		Info: PrimeRecord{
			Fixed:  byteToInsulin(data[2], 22),
			Manual: byteToInsulin(data[4], 22),
		},
		Time: decodeTime(data[5:10]),
		Data: data[:10],
	}
}

func decodeAlarm(data []byte, family Family) HistoryRecord {
	r := HistoryRecord{
		Time: decodeTime(data[4:9]),
		Data: data[:9],
	}
	r.Info = AlarmCode(data[1])
	return r
}

func decodeDailyTotal(data []byte, family Family) HistoryRecord {
	t := decodeDate(data[5:7])
	total := twoByteInsulin(data[3:5], 23)
	if family <= 22 {
		return HistoryRecord{
			Time: t,
			Info: total,
			Data: data[:7],
		}
	}
	return HistoryRecord{
		Time: t,
		Info: total,
		Data: data[:10],
	}
}

// Note that this is a different format than the response to BasalRates.
func decodeBasalRate(data []byte) BasalRate {
	return BasalRate{
		Start: halfHoursToTimeOfDay(data[0]),
		Rate:  twoByteInsulinLE(data[1:3]),
	}
}

func decodeBasalProfile(data []byte, family Family) HistoryRecord {
	r := decodeBase(data, family)
	body := data[7:]
	var sched BasalRateSchedule
	for i := 0; i < 144; i += 3 {
		if body[i] == 0x3F {
			break
		}
		b := decodeBasalRate(body[i : i+3])
		// Don't stop if the 00:00 rate happens to be zero.
		if i > 0 && b.Start == 0 && b.Rate == 0 {
			break
		}
		sched = append(sched, b)
	}
	r.Info = sched
	r.Data = data[:152]
	return r
}

var decodeBasalProfileBefore = decodeBasalProfile

var decodeBasalProfileAfter = decodeBasalProfile

func decodeBGCapture(data []byte, family Family) HistoryRecord {
	r := decodeBase(data, family)
	units := GlucoseUnitsType((data[4] >> 5) & 0x3)
	r.Info = GlucoseRecord{
		Units:   units,
		Glucose: intToGlucose(int(data[4]>>7)<<9|int(data[6]>>7)<<8|int(data[1]), units),
	}
	return r
}

func decodeSensorAlarm(data []byte, family Family) HistoryRecord {
	r := HistoryRecord{
		Time: decodeTime(data[3:8]),
		Data: data[:8],
	}
	r.Info = int(data[1])
	return r
}

var decodeClearAlarm = decodeValue

var decodeChangeBasalPattern = decodeValue

func decodeTempBasalDuration(data []byte, family Family) HistoryRecord {
	r := decodeBase(data, family)
	r.Info = halfHoursToDuration(data[1])
	return r
}

var decodeChangeTime = decodeBase

var decodeNewTime = decodeBase

var decodeLowBattery = decodeBase

var decodeBatteryChange = decodeBase

func decodeSetAutoOff(data []byte, family Family) HistoryRecord {
	r := decodeBase(data, family)
	r.Info = hoursToDuration(data[1])
	return r
}

var decodePrepareInsulinChange = decodeBase

var decodeSuspendPump = decodeBase

var decodeResumePump = decodeBase

var decodeSelfTest = decodeBase

var decodeRewind = decodeBase

var decodeClearSettings = decodeBase

var decodeEnableChildBlock = decodeEnable

var decodeMaxBolus = decodeInsulin

var decodeEnableRemote = decodeEnableN(21)

var decodeMaxBasal = decodeInsulin

var decodeEnableBolusWizard = decodeEnable

var decodeUnknown2E = decodeBaseN(107)

func decodeBolusWizard512(data []byte, family Family) HistoryRecord {
	r := decodeBase(data, family)
	bg := int(data[1])
	body := data[7:]
	bgU := GlucoseUnitsType((body[1] >> 6) & 0x3)
	carbU := CarbUnitsType((body[1] >> 4) & 0x3)
	info := BolusWizardRecord{
		GlucoseInput: intToGlucose(int(body[1]&0x3)<<8|bg, bgU),
		CarbInput:    Carbs(int(body[1]&0xC)<<6 | int(body[0])),
		GlucoseUnits: bgU,
		CarbUnits:    carbU,
		TargetLow:    byteToGlucose(body[4], bgU),
		Sensitivity:  byteToGlucose(body[3], bgU),
		CarbRatio:    intToRatio(int(body[2]), carbU, family),
		Correction:   intToInsulin(int(body[7])+int(body[5]&0xF), 12),
		Food:         byteToInsulin(body[6], 12),
		Unabsorbed:   byteToInsulin(body[9], 12),
		Bolus:        byteToInsulin(body[11], 12),
	}
	info.TargetHigh = info.TargetLow
	r.Info = info
	r.Data = data[:19]
	return r
}

var decodeUnabsorbedInsulin512 = decodeUnabsorbedInsulin

var decodeChangeBGReminder = decodeBase

var decodeSetAlarmClockTime = decodeBase

func decodeTempBasalRate(data []byte, family Family) HistoryRecord {
	r := decodeBase(data, family)
	tb := TempBasalRecord{
		Type: TempBasalType(data[7] >> 3),
	}
	if tb.Type == Absolute {
		tb.Value = intToInsulin(int(data[7]&0x7)<<8|int(data[1]), 23)
	} else {
		tb.Value = int(data[1])
	}
	r.Info = tb
	r.Data = data[:8]
	return r
}

func decodeLowReservoir(data []byte, family Family) HistoryRecord {
	r := decodeBase(data, family)
	r.Info = byteToInsulin(data[1], 22)
	return r
}

var decodeAlarmClock = decodeBase

var decodeChangeMeterID = decodeBaseN(21)

var decodeBGReceived512 = decodeBGReceived

var decodeConfirmInsulinChange = decodeEnable

var decodeSensorStatus = decodeEnable

var decodeEnableMeter = decodeEnableN(21)

func decodeBGReceived(data []byte, family Family) HistoryRecord {
	r := decodeBase(data, family)
	r.Data = data[:10]
	units := MgPerDeciLiter // are units encoded somehow?
	r.Info = GlucoseRecord{
		Units:   units,
		Glucose: intToGlucose(int(data[1])<<3|int(data[4]>>5), units),
		MeterID: fmt.Sprintf("%06X", data[7:10]),
	}
	return r
}

func decodeMealMarker(data []byte, family Family) HistoryRecord {
	r := decodeBase(data, family)
	r.Data = data[:9]
	r.Info = CarbRecord{
		Carbs: Carbs(int(data[1])<<8 | int(data[7])),
		Units: CarbUnitsType(data[8] & 0x3),
	}
	return r
}

var decodeExerciseMarker = decodeBaseN(8)

func decodeInsulinMarker(data []byte, family Family) HistoryRecord {
	r := decodeBase(data, family)
	r.Data = data[:8]
	r.Info = intToInsulin(int(data[4]&0x60)<<3|int(data[1]), 22)
	return r
}

var decodeOtherMarker = decodeBase

var decodeEnableSensorAutoCal = decodeBase

var decodeChangeBolusWizardSetup = decodeBaseN(39)

func decodeSensorSetup(data []byte, family Family) HistoryRecord {
	var d decoder
	if family >= 51 {
		d = decodeBaseN(41)
	} else {
		d = decodeBaseN(37)
	}
	return d(data, family)
}

var decodeSensor51 = decodeBase

var decodeSensor52 = decodeBase

var decodeChangeSensorAlarm = decodeBaseN(8)

var decodeSensor54 = decodeBaseN(64)

var decodeSensor55 = decodeBaseN(55)

var decodeChangeSensorAlert = decodeBaseN(12)

var decodeChangeBolusStep = decodeBase

func decodeBolusWizardConfig(data []byte, family Family) BolusWizardConfig {
	const numEntries = 8
	r := BolusWizardConfig{}
	carbUnits := CarbUnitsType(data[0] & 0x3)
	if carbUnits != Exchanges {
		// Pumps can return 0 when first turned on.
		carbUnits = Grams
	}
	bgUnits := GlucoseUnitsType((data[0] >> 2) & 0x3)
	if bgUnits != MMolPerLiter {
		// Pumps can return 0 when first turned on.
		bgUnits = MgPerDeciLiter
	}
	step := carbRatioStep(family)
	r.Ratios = decodeCarbRatioSchedule(data[2:2+numEntries*step], carbUnits, family)
	data = data[2+numEntries*step:]
	r.Sensitivities = decodeInsulinSensitivitySchedule(data[:numEntries*2], bgUnits)
	if family <= 22 {
		data = data[numEntries*2:]
	} else {
		data = data[numEntries*2+2:]
	}
	r.Targets = decodeGlucoseTargetSchedule(data[:numEntries*3], bgUnits, family)
	return r
}

func decodeBolusWizardSetup(data []byte, family Family) HistoryRecord {
	r := decodeBase(data, family)
	if family <= 22 {
		r.Data = data[:124]
	} else {

		r.Data = data[:144]
	}
	n := len(r.Data) - 1
	body := data[7:n]
	half := (n - 7) / 2
	setup := BolusWizardSetupRecord{
		Before: decodeBolusWizardConfig(body[:half], family),
		After:  decodeBolusWizardConfig(body[half:], family),
	}
	setup.Before.InsulinAction = hoursToDuration(data[n] & 0xF)
	setup.After.InsulinAction = hoursToDuration(data[n] >> 4)
	r.Info = setup
	return r
}

func decodeBolusWizard(data []byte, family Family) HistoryRecord {
	r := decodeBase(data, family)
	bg := int(data[1])
	body := data[7:]
	bgU := GlucoseUnitsType((body[1] >> 6) & 0x3)
	carbU := CarbUnitsType((body[1] >> 4) & 0x3)
	if family <= 22 {
		r.Info = BolusWizardRecord{
			GlucoseInput: intToGlucose(int(body[1]&0x3)<<8|bg, bgU),
			CarbInput:    Carbs(int(body[1]&0xC)<<6 | int(body[0])),
			GlucoseUnits: bgU,
			CarbUnits:    carbU,
			TargetLow:    byteToGlucose(body[4], bgU),
			TargetHigh:   byteToGlucose(body[12], bgU),
			Sensitivity:  byteToGlucose(body[3], bgU),
			CarbRatio:    intToRatio(int(body[2]), carbU, family),
			Correction:   intToInsulin(int(body[7])+int(body[5]&0xF), family),
			Food:         byteToInsulin(body[6], family),
			Unabsorbed:   byteToInsulin(body[9], family),
			Bolus:        byteToInsulin(body[11], family),
		}
		r.Data = data[:20]
	} else {
		r.Info = BolusWizardRecord{
			GlucoseInput: intToGlucose(int(body[1]&0x3)<<8|bg, bgU),
			CarbInput:    Carbs(int(body[1]&0xC)<<6 | int(body[0])),
			GlucoseUnits: bgU,
			CarbUnits:    carbU,
			TargetLow:    byteToGlucose(body[5], bgU),
			TargetHigh:   byteToGlucose(body[14], bgU),
			Sensitivity:  byteToGlucose(body[4], bgU),
			CarbRatio:    intToRatio(int(body[2]&0xF)<<8|int(body[3]), carbU, family),
			Correction:   intToInsulin(int(body[9]&0x38)<<5|int(body[6]), family),
			Food:         twoByteInsulin(body[7:9], family),
			Unabsorbed:   twoByteInsulin(body[10:12], family),
			Bolus:        twoByteInsulin(body[12:14], family),
		}
		r.Data = data[:22]
	}
	return r
}

func decodeUnabsorbedInsulin(data []byte, family Family) HistoryRecord {
	n := int(data[1]) - 2
	body := data[2:]
	var unabsorbed UnabsorbedBolusHistory
	for i := 0; i < n; i += 3 {
		amount := byteToInsulin(body[i], 23)
		curve := body[i+2]
		age := minutesToDuration(body[i+1] + (curve&0x30)<<4)
		unabsorbed = append(unabsorbed, UnabsorbedBolus{
			Bolus: amount,
			Age:   age,
		})
	}
	return HistoryRecord{
		Data: data[:n+2],
		Info: unabsorbed,
	}
}

var decodeSaveSettings = decodeBase

var decodeEnableVariableBolus = decodeEnable

var decodeChangeEasyBolus = decodeBase

var decodeEnableBGReminder = decodeEnable

var decodeEnableAlarmClock = decodeEnable

func decodeChangeTempBasalType(data []byte, family Family) HistoryRecord {
	r := decodeBase(data, family)
	r.Info = TempBasalType(data[1])
	return r
}

var decodeChangeAlarmType = decodeValue

var decodeChangeTimeFormat = decodeValue

func decodeChangeReservoirWarning(data []byte, family Family) HistoryRecord {
	r := decodeBase(data, family)
	v := data[1]
	if v&0x1 == 0 {
		r.Info = Insulin(1000 * int(v>>2))
	} else {
		r.Info = halfHoursToDuration(v >> 2)
	}
	return r
}

var decodeEnableBolusReminder = decodeEnable

var decodeSetBolusReminderTime = decodeEnableN(9)

var decodeDeleteBolusReminderTime = decodeEnableN(9)

var decodeBolusReminder = decodeBaseN(9)

var decodeDeleteAlarmClockTime = decodeBase

var decodeDailyTotal515 = decodeDailyTotalN(38)

var decodeDailyTotal522 = decodeDailyTotalN(44)

var decodeDailyTotal523 = decodeDailyTotalN(52)

var decodeChangeCarbUnits = decodeValue

func decodeBasalProfileStart(data []byte, family Family) HistoryRecord {
	r := decodeBase(data, family)
	r.Info = BasalProfileStartRecord{
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
func DecodeHistoryRecord(data []byte, family Family) (HistoryRecord, error) {
	if len(data) == 0 {
		return HistoryRecord{}, fmt.Errorf("empty history record")
	}
	decoder := decode[HistoryRecordType(data[0])]
	if decoder == nil {
		return HistoryRecord{}, unknownRecord(data)
	}
	return decoder(data, family), nil
}

func (e UnknownRecordTypeError) Error() string {
	return fmt.Sprintf("unknown record type here: % X", e.Data)
}

func unknownRecord(data []byte) error {
	return UnknownRecordTypeError{
		Data: data,
	}
}

// DecodeHistory decodes the records in a page of data and
// returns them in reverse chronological order (most recent first),
// to match the order of the history pages themselves.
func DecodeHistory(data []byte, family Family) (History, error) {
	var results History
	var r HistoryRecord
	var err error
	for !allZero(data) {
		r, err = DecodeHistoryRecord(data, family)
		if err != nil {
			break
		}
		results = append(results, r)
		data = data[len(r.Data):]
	}
	ReverseHistory(results)
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

// ReverseHistory reverses a slice of history records.
func ReverseHistory(a History) {
	for i, j := 0, len(a)-1; i < len(a)/2; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}
