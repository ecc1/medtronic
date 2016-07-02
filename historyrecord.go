package medtronic

import (
	"fmt"
	"log"
	"time"
)

type HistoryRecordType byte

//go:generate stringer -type HistoryRecordType

type HistoryRecord interface {
	Type() HistoryRecordType
	Time() time.Time
	Decode([]byte, bool)
	Length() int
	Bytes() []byte
}

type BaseRecord struct {
	Value     int
	Timestamp time.Time
	Data      []byte
}

func (r *BaseRecord) Type() HistoryRecordType { return HistoryRecordType(r.Data[0]) }
func (r *BaseRecord) Time() time.Time         { return r.Timestamp }
func (r *BaseRecord) Length() int             { return len(r.Data) }
func (r *BaseRecord) Bytes() []byte           { return r.Data }

func (r *BaseRecord) Decode(data []byte, newerPump bool) {
	r.Value = int(data[1])
	r.Timestamp = decodeTimestamp(data[2:7])
	r.Data = data[:7]
}

const Bolus HistoryRecordType = 0x01

type BolusRecord struct {
	BaseRecord
	Programmed int           // milliUnits
	Amount     int           // milliUnits
	Unabsorbed int           // milliUnits
	Duration   time.Duration // non-zero for square wave bolus
}

func (r *BolusRecord) Decode(data []byte, newerPump bool) {
	if newerPump {
		r.Programmed = twoByteInt(data[1:3]) * 25
		r.Amount = twoByteInt(data[3:5]) * 25
		r.Unabsorbed = twoByteInt(data[5:7]) * 25
		r.Duration = scheduleToDuration(data[7])
		r.Timestamp = decodeTimestamp(data[8:13])
		r.Data = data[:13]
	} else {
		r.Programmed = int(data[1]) * 100
		r.Amount = int(data[2]) * 100
		r.Duration = scheduleToDuration(data[3])
		r.Timestamp = decodeTimestamp(data[4:9])
		r.Data = data[:9]
	}
}

const Prime HistoryRecordType = 0x03

type PrimeRecord struct {
	BaseRecord     // Value = manual amount in milliUnits
	Fixed      int // milliUnits
}

func (r *PrimeRecord) Decode(data []byte, newerPump bool) {
	r.BaseRecord.Decode(data[3:], newerPump)
	r.Value *= 100
	r.Fixed = int(data[2]) * 100
	r.Data = data[:10]
}

const Alarm HistoryRecordType = 0x06

type AlarmRecord struct {
	BaseRecord
	AlarmType int
}

func (r *AlarmRecord) Decode(data []byte, newerPump bool) {
	r.BaseRecord.Decode(data[2:], newerPump)
	r.AlarmType = int(data[1])
	// data[1:3] = line number (?)
	r.Data = data[:9]
}

const DailyTotal HistoryRecordType = 0x07

type DailyTotalRecord struct {
	BaseRecord // Value = daily total in milliUnits
}

func (r *DailyTotalRecord) Decode(data []byte, newerPump bool) {
	r.Value = twoByteInt(data[3:5]) * 25
	r.Timestamp = decodeDate(data[5:7])
	if newerPump {
		r.Data = data[:10]
	} else {
		r.Data = data[:7]
	}
}

type BasalProfileRecord struct {
	BaseRecord
	Rates BasalRateSchedule
}

// Note that this is a different format than the response to BasalRates.
func decodeBasalRate(data []byte) BasalRate {
	return BasalRate{
		Start: scheduleToDuration(data[0]),
		Rate:  int(data[1]) * 25,
		// data[2] unused
	}
}

func (r *BasalProfileRecord) Decode(data []byte, newerPump bool) {
	r.BaseRecord.Decode(data, newerPump)
	body := data[7:]
	for i := 0; i < 144; i += 3 {
		b := decodeBasalRate(body[i : i+3])
		// Don't stop if the 00:00 rate happens to be zero.
		if i > 0 && b.Start == 0 && b.Rate == 0 {
			break
		}
		r.Rates.Schedule = append(r.Rates.Schedule, b)
	}
	r.Data = data[:152]
}

const BasalProfileBefore HistoryRecordType = 0x08

type BasalProfileBeforeRecord struct {
	BasalProfileRecord
}

const BasalProfileAfter HistoryRecordType = 0x09

type BasalProfileAfterRecord struct {
	BasalProfileRecord
}

const BgCapture HistoryRecordType = 0x0A

type BgCaptureRecord struct {
	BaseRecord // Value = BG input
}

func (r *BgCaptureRecord) Decode(data []byte, newerPump bool) {
	r.BaseRecord.Decode(data, newerPump)
	r.Value |= int(data[6]>>7) << 8
	r.Data = data[:7]
}

const ClearAlarm HistoryRecordType = 0x0C

type ClearAlarmRecord struct {
	BaseRecord // Value = alarm type
}

const TempBasalDuration HistoryRecordType = 0x16

type TempBasalDurationRecord struct {
	BaseRecord
	Duration time.Duration
}

func (r *TempBasalDurationRecord) Decode(data []byte, newerPump bool) {
	r.BaseRecord.Decode(data, newerPump)
	r.Duration = scheduleToDuration(uint8(r.Value))
	r.Data = data[:7]
}

const ChangeTime HistoryRecordType = 0x17

type ChangeTimeRecord struct {
	BaseRecord
}

const NewTime HistoryRecordType = 0x18

type NewTimeRecord struct {
	BaseRecord
}

const LowBattery HistoryRecordType = 0x19

type LowBatteryRecord struct {
	BaseRecord
}

const BatteryChange HistoryRecordType = 0x1A

type BatteryChangeRecord struct {
	BaseRecord
}

const SetAutoOff HistoryRecordType = 0x1B

type SetAutoOffRecord struct {
	BaseRecord
}

const SuspendPump HistoryRecordType = 0x1E

type SuspendPumpRecord struct {
	BaseRecord
}

const ResumePump HistoryRecordType = 0x1F

type ResumePumpRecord struct {
	BaseRecord
}

const Rewind HistoryRecordType = 0x21

type RewindRecord struct {
	BaseRecord
}

const EnableRemote HistoryRecordType = 0x26

type EnableRemoteRecord struct {
	BaseRecord // Value = 1 if enabled
}

func (r *EnableRemoteRecord) Decode(data []byte, newerPump bool) {
	r.BaseRecord.Decode(data, newerPump)
	r.Data = data[:21]
}

const TempBasalRate HistoryRecordType = 0x33

type TempBasalRateRecord struct {
	BaseRecord    // Value = rate
	TempBasalType TempBasalType
}

func (r *TempBasalRateRecord) Decode(data []byte, newerPump bool) {
	r.BaseRecord.Decode(data, newerPump)
	r.TempBasalType = TempBasalType(data[7] >> 3)
	switch r.TempBasalType {
	case Absolute:
		r.Value *= 25
	case Percent:
	default:
		log.Panicf("unexpected temp basal type %02X", r.TempBasalType)
	}
	r.Data = data[:8]
}

const LowReservoir HistoryRecordType = 0x34

type LowReservoirRecord struct {
	BaseRecord // Value = reservoir remaining in milliUnits
}

func (r *LowReservoirRecord) Decode(data []byte, newerPump bool) {
	r.BaseRecord.Decode(data, newerPump)
	r.Value *= 100
	r.Data = data[:7]
}

const SensorStatus HistoryRecordType = 0x3B

type SensorStatusRecord struct {
	BaseRecord // Value = 1 if enabled
}

const EnableMeter HistoryRecordType = 0x3C

type EnableMeterRecord struct {
	BaseRecord // Value = 1 if enabled
}

func (r *EnableMeterRecord) Decode(data []byte, newerPump bool) {
	r.BaseRecord.Decode(data, newerPump)
	r.Data = data[:21]
}

const BolusWizardSetup HistoryRecordType = 0x5A

type BolusWizardConfig struct {
	Ratios        CarbRatioSchedule
	Sensitivities InsulinSensitivitySchedule
	Targets       GlucoseTargetSchedule
}

func decodeBolusWizardConfig(data []byte, newerPump bool) BolusWizardConfig {
	const numEntries = 8
	conf := BolusWizardConfig{}
	carbUnits := Grams // FIXME
	step := carbRatioStep(newerPump)
	conf.Ratios = decodeCarbRatioSchedule(data[2:2+numEntries*step], carbUnits, newerPump)
	data = data[2+numEntries*step:]
	bgUnits := MgPerDeciLiter // FIXME
	conf.Sensitivities = decodeInsulinSensitivitySchedule(data[:numEntries*2], bgUnits)
	if newerPump {
		data = data[numEntries*2+2:]
	} else {
		data = data[numEntries*2:]
	}
	conf.Targets = decodeGlucoseTargetSchedule(data[:numEntries*3], bgUnits)
	return conf
}

type BolusWizardSetupRecord struct {
	BaseRecord
	Before BolusWizardConfig
	After  BolusWizardConfig
}

func (r *BolusWizardSetupRecord) Decode(data []byte, newerPump bool) {
	r.BaseRecord.Decode(data, newerPump)
	if newerPump {
		r.Data = data[:144]
	} else {
		r.Data = data[:124]
	}
	n := r.Length() - 1
	body := data[7:n]
	half := (n - 7) / 2
	r.Before = decodeBolusWizardConfig(body[:half], newerPump)
	r.After = decodeBolusWizardConfig(body[half:], newerPump)
}

const BolusWizard HistoryRecordType = 0x5B

type BolusWizardRecord struct {
	BaseRecord
	GlucoseInput int // mg/dL or μmol/L
	TargetLow    int // mg/dL or μmol/L
	TargetHigh   int // mg/dL or μmol/L
	Sensitivity  int // mg/dL or μmol/L reduction per insulin unit
	CarbInput    int // grams or exchanges
	CarbRatio    int // grams or exchanges covered by TEN insulin units
	Unabsorbed   int // milliUnits
	Correction   int // milliUnits
	Food         int // milliUnits
	Bolus        int // milliUnits
}

func (r *BolusWizardRecord) Decode(data []byte, newerPump bool) {
	r.BaseRecord.Decode(data, newerPump)
	body := data[7:]
	if newerPump {
		r.GlucoseInput = int(body[1]&0x3)<<8 | r.Value
		r.CarbInput = int(body[1]&0xC)<<6 + int(body[0])
		r.CarbRatio = (int(body[2]&0x7)<<8 | int(body[3]))
		r.Sensitivity = int(body[4])
		r.TargetLow = int(body[5])
		r.Correction = (int(body[9]&0x38)<<5 + int(body[6])) * 25
		r.Food = twoByteInt(body[7:9]) * 25
		r.Unabsorbed = twoByteInt(body[10:12]) * 25
		r.Bolus = twoByteInt(body[12:14]) * 25
		r.TargetHigh = int(body[14])
		r.Data = data[:22]
	} else {
		r.GlucoseInput = int(body[1]&0xF)<<8 | r.Value
		r.CarbInput = int(body[0])
		r.CarbRatio = int(body[2])
		r.Sensitivity = int(body[3])
		r.TargetLow = int(body[4])
		r.Correction = (int(body[7]) + int(body[5]&0xF)) * 100
		r.Food = int(body[6]) * 100
		r.Unabsorbed = int(body[9]) * 100
		r.Bolus = int(body[11]) * 100
		r.TargetHigh = int(body[12])
		r.Data = data[:20]
	}
}

const UnabsorbedInsulin HistoryRecordType = 0x5C

type BolusInfo struct {
	Amount        int
	Age           time.Duration
	InsulinAction time.Duration
}

type UnabsorbedInsulinRecord struct {
	BaseRecord
	Unabsorbed []BolusInfo
}

func (r *UnabsorbedInsulinRecord) Decode(data []byte, newerPump bool) {
	n := int(data[1]) - 2
	body := data[2:]
	for i := 0; i < n; i += 3 {
		amount := int(body[i]) * 25 // milliUnits
		curve := body[i+2]
		age := time.Duration(body[i+1]+(curve&0x30)<<4) * time.Minute
		action := time.Duration(curve&0xF) * time.Hour
		bolus := BolusInfo{
			Amount:        amount,
			Age:           age,
			InsulinAction: action,
		}
		r.Unabsorbed = append(r.Unabsorbed, bolus)
	}
	r.Data = data[:n+2]
}

const ChangeTempBasalType HistoryRecordType = 0x62

type ChangeTempBasalTypeRecord struct {
	BaseRecord
	TempBasalType TempBasalType
}

func (r *ChangeTempBasalTypeRecord) Decode(data []byte, newerPump bool) {
	r.BaseRecord.Decode(data, newerPump)
	r.TempBasalType = TempBasalType(r.Value)
	r.Data = data[:7]
}

const ChangeTimeDisplay HistoryRecordType = 0x64

type ChangeTimeDisplayRecord struct {
	BaseRecord
}

const DailyTotal522 HistoryRecordType = 0x6D

type DailyTotal522Record struct {
	BaseRecord
}

func (r *DailyTotal522Record) Decode(data []byte, newerPump bool) {
	r.Timestamp = decodeDate(data[1:3])
	r.Data = data[:44]
}

const DailyTotal523 HistoryRecordType = 0x6E

type DailyTotal523Record struct {
	BaseRecord
}

func (r *DailyTotal523Record) Decode(data []byte, newerPump bool) {
	r.Timestamp = decodeDate(data[1:3])
	r.Data = data[:52]
}

const BasalProfileStart HistoryRecordType = 0x7B

type BasalProfileStartRecord struct {
	BaseRecord // Value = profile index
	BasalRate  BasalRate
}

func (r *BasalProfileStartRecord) Decode(data []byte, newerPump bool) {
	r.BaseRecord.Decode(data, newerPump)
	r.BasalRate = decodeBasalRate(data[7:10])
	r.Data = data[:10]
}

const ConnectOtherDevices HistoryRecordType = 0x7C

type ConnectOtherDevicesRecord struct {
	BaseRecord // Value = 1 if enabled
}

func DecodeHistoryRecord(data []byte, newerPump bool) HistoryRecord {
	var r HistoryRecord
	switch HistoryRecordType(data[0]) {
	case Bolus:
		r = &BolusRecord{}
	case Prime:
		r = &PrimeRecord{}
	case Alarm:
		r = &AlarmRecord{}
	case DailyTotal:
		r = &DailyTotalRecord{}
	case BasalProfileBefore:
		r = &BasalProfileBeforeRecord{}
	case BasalProfileAfter:
		r = &BasalProfileAfterRecord{}
	case BgCapture:
		r = &BgCaptureRecord{}
	case ClearAlarm:
		r = &ClearAlarmRecord{}
	case TempBasalDuration:
		r = &TempBasalDurationRecord{}
	case ChangeTime:
		r = &ChangeTimeRecord{}
	case NewTime:
		r = &NewTimeRecord{}
	case LowBattery:
		r = &LowBatteryRecord{}
	case BatteryChange:
		r = &BatteryChangeRecord{}
	case SetAutoOff:
		r = &SetAutoOffRecord{}
	case SuspendPump:
		r = &SuspendPumpRecord{}
	case ResumePump:
		r = &ResumePumpRecord{}
	case Rewind:
		r = &RewindRecord{}
	case EnableRemote:
		r = &EnableRemoteRecord{}
	case TempBasalRate:
		r = &TempBasalRateRecord{}
	case LowReservoir:
		r = &LowReservoirRecord{}
	case SensorStatus:
		r = &SensorStatusRecord{}
	case EnableMeter:
		r = &EnableMeterRecord{}
	case BolusWizardSetup:
		r = &BolusWizardSetupRecord{}
	case BolusWizard:
		r = &BolusWizardRecord{}
	case UnabsorbedInsulin:
		r = &UnabsorbedInsulinRecord{}
	case ChangeTempBasalType:
		r = &ChangeTempBasalTypeRecord{}
	case ChangeTimeDisplay:
		r = &ChangeTimeDisplayRecord{}
	case DailyTotal522:
		r = &DailyTotal522Record{}
	case DailyTotal523:
		r = &DailyTotal523Record{}
	case BasalProfileStart:
		r = &BasalProfileStartRecord{}
	case ConnectOtherDevices:
		r = &ConnectOtherDevicesRecord{}
	default:
		panic(fmt.Sprintf("unknown record type here: % X", data))
	}
	r.Decode(data, newerPump)
	return r
}

// Decode records in a page of data and return them in reverse order,
// to match the order of the history pages themselves.
func DecodeHistoryRecords(data []byte, newerPump bool) []HistoryRecord {
	results := []HistoryRecord{}
	for !allZero(data) {
		r := DecodeHistoryRecord(data, newerPump)
		results = append(results, r)
		data = data[r.Length():]
	}
	reverseHistoryRecords(results)
	return results
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

func reverseHistoryRecords(a []HistoryRecord) {
	for i, j := 0, len(a)-1; i < len(a)/2; i,j = i+1,j-1 {
		a[i], a[j] = a[j], a[i]
	}
}
