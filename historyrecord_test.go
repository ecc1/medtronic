package medtronic

import (
	"encoding/json"
	"log"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestDecodeHistoryRecord(t *testing.T) {
	cases := []struct {
		test  string
		newer bool
	}{
		{
			`{
			   "Type": "Bolus",
			   "Time": "2016-07-04T14:25:46-04:00",
			   "Bolus": {
			     "Duration": "0",
			     "Programmed": 0.5,
			     "Amount": 0.5,
			     "Unabsorbed": 0
			   },
			   "Data": "AQUFAG7ZTgQQ"
			 }`,
			false,
		},
		{
			`{
			   "Type": "Bolus",
			   "Time": "2016-07-04T14:27:08-04:00",
			   "Bolus": {
			     "Duration": "30m0s",
			     "Programmed": 0.1,
			     "Amount": 0.1,
			     "Unabsorbed": 0
			   },
			   "Data": "AQEBAUjbbgQQ"
			 }`,
			false,
		},
		{
			`{
			   "Type": "Bolus",
			   "Time": "2016-07-08T07:51:06-04:00",
			   "Bolus": {
			     "Duration": "0",
			     "Programmed": 2.4,
			     "Amount": 2.4,
			     "Unabsorbed": 0
			   },
			   "Data": "AQBgAGAAAABG80doEA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "Bolus",
			   "Time": "2016-05-25T14:27:43-04:00",
			   "Bolus": {
			     "Duration": "1h0m0s",
			     "Programmed": 4.6,
			     "Amount": 4.6,
			     "Unabsorbed": 1.1
			   },
			   "Data": "AQC4ALgALAJrW655EA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "Prime",
			   "Time": "2016-07-07T15:25:05-04:00",
			   "Prime": {
			     "Fixed": 0.5,
			     "Manual": 0.5
			   },
			   "Data": "AwAFAAVF2Q8HEA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "Alarm",
			   "Time": "2016-07-07T15:19:07-04:00",
			   "Value": 4,
			   "Data": "BgQMW0fTT0cQ"
			 }`,
			true,
		},
		{
			`{
			   "Type": "DailyTotal",
			   "Time": "2016-07-01T00:00:00-04:00",
			   "Insulin": 11.5,
			   "Data": "BwAAAcxhkA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "DailyTotal",
			   "Time": "2016-07-01T00:00:00-04:00",
			   "Insulin": 51.9,
			   "Data": "BwAACBxhkAAAAA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "BasalProfileBefore",
			   "Time": "2016-06-15T16:49:16-04:00",
			   "BasalProfile": [
			     {
			       "Start": "00:00",
			       "Rate": 1
			     },
			     {
			       "Start": "03:00",
			       "Rate": 0.8
			     },
			     {
			       "Start": "07:00",
			       "Rate": 1.1
			     },
			     {
			       "Start": "10:00",
			       "Rate": 1.2
			     }
			   ],
			   "Data": "CARQsRAPEAAoAAYgAA4sABQwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
			 }`,
			false,
		},
		{
			`{
			   "Type": "BasalProfileAfter",
			   "Time": "2016-05-22T23:39:53-04:00",
			   "BasalProfile": [
			     {
			       "Start": "00:00",
			       "Rate": 0.8
			     },
			     {
			       "Start": "03:00",
			       "Rate": 0.8
			     },
			     {
			       "Start": "06:00",
			       "Rate": 1.1
			     },
			     {
			       "Start": "10:00",
			       "Rate": 1.2
			     },
			     {
			       "Start": "22:00",
			       "Rate": 1.2
			     }
			   ],
			   "Data": "CQV1ZxcWEAAgAAYgAAwsABQwACwwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
			 }`,
			true,
		},
		{
			`{
			   "Type": "BgCapture",
			   "Time": "2016-07-04T14:25:37-04:00",
			   "Glucose": 110,
			   "GlucoseUnits": "mg/dL",
			   "Data": "Cm5l2S4EEA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "BgCapture",
			   "Time": "2016-05-19T03:49:06-04:00",
			   "Glucose": 500,
			   "GlucoseUnits": "mg/dL",
			   "Data": "CvRGcSMTkA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "BGCapture",
			   "Time": "2017-02-26T17:12:08-05:00",
			   "Glucose": 83,
			   "GlucoseUnits": "μmol/L",
			   "Data": "ClMIjFEaEQ=="
			 }`,
			true,
		},

		{
			`{
			   "Type": "SensorAlarm",
			   "Time": "2016-07-04T14:33:52-04:00",
			   "Value": 113,
			   "Data": "C3EAdOEupBA="
			 }`,
			false,
		},
		{
			`{
			   "Type": "ClearAlarm",
			   "Time": "2005-01-01T00:00:31-05:00",
			   "Value": 55,
			   "Data": "DDcfQAABBQ=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "ChangeBasalPattern",
			   "Time": "2016-07-04T14:11:18-04:00",
			   "Value": 2,
			   "Data": "FAJSyw4EEA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "TempBasalDuration",
			   "Time": "2016-06-14T22:47:04-04:00",
			   "Duration": "1h30m0s",
			   "Data": "FgNErxZOEA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "ChangeTime",
			   "Time": "2016-07-05T13:08:24-04:00",
			   "Data": "FwBYyA0FEA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "NewTime",
			   "Time": "2016-06-11T20:42:45-04:00",
			   "Data": "GABtqhRLEA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "LowBattery",
			   "Time": "2016-05-23T23:01:25-04:00",
			   "Data": "GQBZQRcXEA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "BatteryChange",
			   "Time": "2016-05-23T23:12:05-04:00",
			   "Data": "GgFFTBcXEA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "SetAutoOff",
			   "Time": "2016-07-02T13:53:28-04:00",
			   "Duration": "13h0m0s",
			   "Data": "Gw1c9Q0CEA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "SuspendPump",
			   "Time": "2016-05-17T13:06:32-04:00",
			   "Data": "HgFgRg0REA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "ResumePump",
			   "Time": "2016-07-03T16:33:47-04:00",
			   "Data": "HwBv4RBDEA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "SelfTest",
			   "Time": "2016-07-04T14:18:24-04:00",
			   "Data": "IABY0g4EEA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "Rewind",
			   "Time": "2016-07-04T14:11:44-04:00",
			   "Data": "IQBsyw4EEA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "EnableChildBlock",
			   "Time": "2016-07-04T14:17:46-04:00",
			   "Enabled": true,
			   "Data": "IwFu0Q4EEA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "MaxBolus",
			   "Time": "2016-07-04T13:59:52-04:00",
			   "Insulin": 3,
			   "Data": "JHh0+w0EEA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "EnableRemote",
			   "Time": "2016-06-05T14:47:32-04:00",
			   "Enabled": true,
			   "Data": "JgFgrw4FECcB4kAAAAAoAAAAAAAA"
			 }`,
			false,
		},
		{
			`{
			   "Type": "MaxBasal",
			   "Time": "2016-07-04T14:08:36-04:00",
			   "Insulin": 2.5,
			   "Data": "LGRkyA4EEA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "EnableBolusWizard",
			   "Time": "2016-07-05T15:44:59-04:00",
			   "Enabled": true,
			   "Data": "LQF77A8FEA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "SetAlarmClockTime",
			   "Time": "2016-07-04T14:15:31-04:00",
			   "Data": "Mhdfzw4EEA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "TempBasalRate",
			   "Time": "2016-06-04T18:52:45-04:00",
			   "Insulin": 1.25,
			   "TempBasalType": "Absolute",
			   "Data": "MzJttBIEEAA="
			 }`,
			false,
		},
		{
			`{
			   "Type": "TempBasalRate",
			   "Time": "2016-07-05T03:13:09-04:00",
			   "Insulin": 0,
			   "TempBasalType": "Absolute",
			   "Data": "MwBJzQMFEAA="
			 }`,
			true,
		},
		{
			`{
			   "Type": "LowReservoir",
			   "Time": "2016-05-12T16:41:13-04:00",
			   "Insulin": 10,
			   "Data": "NGRNaRAMEA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "SensorStatus",
			   "Time": "2016-06-24T18:05:44-04:00",
			   "Enabled": true,
			   "Data": "O4NshRIYEA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "EnableMeter",
			   "Time": "2016-06-16T14:42:52-04:00",
			   "Enabled": true,
			   "Data": "PAF0qg4QED0SNFYAAAA+AAAAAAAA"
			 }`,
			true,
		},
		{
			`{
			   "Type": "BGReceived",
			   "Time": "2016-07-10T20:37:26-04:00",
			   "Data": "PxRa5XRqEM4mkg=="
			 }`,
			true,
		},

		{
			`{
			   "Type": "MealMarker",
			   "Time": "2016-07-11T05:09:29-04:00",
			   "Carbs": 12,
			   "CarbUnits": "Grams",
			   "Data": "QABdyQULEAwB"
			 }`,
			true,
		},
		{
			`{
			   "Type": "ChangeBolusWizardSetup",
			   "Time": "2016-07-04T14:06:42-04:00",
			   "Data": "TwBqxg4EEEARAG8nFh4APBQAHjwg9PdBUQBuIxYKAB4eAB48EtaH"
			 }`,
			false,
		},
		{
			`{
			   "Type": "BolusWizardSetup",
			   "Time": "2016-07-04T14:03:03-04:00",
			   "BolusWizardSetup": {
			     "Before": {
			       "InsulinAction": "4h0m0s",
			       "Ratios": [
				 {
				   "Start": "00:00",
				   "Ratio": 15,
				   "Units": "Grams"
				 },
				 {
				   "Start": "04:00",
				   "Ratio": 20,
				   "Units": "Grams"
				 }
			       ],
			       "Sensitivities": [
				 {
				   "Start": "00:00",
				   "Sensitivity": 50,
				   "Units": "mg/dL"
				 },
				 {
				   "Start": "06:00",
				   "Sensitivity": 40,
				   "Units": "mg/dL"
				 }
			       ],
			       "Targets": [
				 {
				   "Start": "00:00",
				   "Low": 95,
				   "High": 105,
				   "Units": "mg/dL"
				 },
				 {
				   "Start": "09:00",
				   "Low": 90,
				   "High": 110,
				   "Units": "mg/dL"
				 },
				 {
				   "Start": "17:00",
				   "Low": 80,
				   "High": 120,
				   "Units": "mg/dL"
				 }
			       ]
			     },
			     "After": {
			       "InsulinAction": "3h0m0s",
			       "Ratios": [
				 {
				   "Start": "00:00",
				   "Ratio": 1,
				   "Units": "Exchanges"
				 },
				 {
				   "Start": "01:00",
				   "Ratio": 1.5,
				   "Units": "Exchanges"
				 },
				 {
				   "Start": "02:00",
				   "Ratio": 1.2,
				   "Units": "Exchanges"
				 }
			       ],
			       "Sensitivities": [
				 {
				   "Start": "00:00",
				   "Sensitivity": 3000,
				   "Units": "μmol/L"
				 },
				 {
				   "Start": "02:00",
				   "Sensitivity": 2500,
				   "Units": "μmol/L"
				 }
			       ],
			       "Targets": [
				 {
				   "Start": "00:00",
				   "Low": 5500,
				   "High": 5700,
				   "Units": "μmol/L"
				 },
				 {
				   "Start": "01:00",
				   "Low": 5600,
				   "High": 5600,
				   "Units": "μmol/L"
				 }
			       ]
			     }
			   },
			   "Data": "Wg9Dww4EECUyAA8IFAAAAAAAAAAAAAAAAAAyDCgAAAAAAAAAAAAAAAAAX2kSWm4iUHgAAAAAAAAAAAAAAAAAAAA6IgAKAg8EDAAAAAAAAAAAAAAAHgQZAAAAAAAAAAAAAAAAADc5Ajg4AAAAAAAAAAAAAAAAAAAAAAAANA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "BolusWizardSetup",
			   "Time": "2017-02-26T17:17:16-05:00",
			   "Data": "Wg8QkREaERkRAAA8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEwAAAAAAAAAAAAAAAAAAADJDAAAAAAAAAAAAAAAAAAAAAAAAAAAAGhEACcQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAATAAAAAAAAAAAAAAAAAAAAMkMAAAAAAAAAAAAAAAAAAAAAAAAAAAAz",
			   "BolusWizardSetup": {
			     "Before": {
			       "InsulinAction": "3h0m0s",
			       "Ratios": [
				 {
				   "Start": "00:00",
				   "Ratio": 6,
				   "Units": "Grams"
				 }
			       ],
			       "Sensitivities": [
				 {
				   "Start": "00:00",
				   "Sensitivity": 0,
				   "Units": "μmol/L"
				 }
			       ],
			       "Targets": [
				 {
				   "Start": "00:00",
				   "Low": 5000,
				   "High": 6700,
				   "Units": "μmol/L"
				 }
			       ]
			     },
			     "After": {
			       "InsulinAction": "3h0m0s",
			       "Ratios": [
				 {
				   "Start": "00:00",
				   "Ratio": 2.5,
				   "Units": "Exchanges"
				 }
			       ],
			       "Sensitivities": [
				 {
				   "Start": "00:00",
				   "Sensitivity": 0,
				   "Units": "μmol/L"
				 }
			       ],
			       "Targets": [
				 {
				   "Start": "00:00",
				   "Low": 5000,
				   "High": 6700,
				   "Units": "μmol/L"
				 }
			       ]
			     }
			   }
			 }`,
			true,
		},
		{
			`{
			   "Type": "BolusWizardSetup",
			   "Time": "2016-05-09T12:07:45-04:00",
			   "BolusWizardSetup": {
			     "Before": {
			       "InsulinAction": "4h0m0s",
			       "Ratios": [
				 {
				   "Start": "00:00",
				   "Ratio": 6,
				   "Units": "Grams"
				 }
			       ],
			       "Sensitivities": [
				 {
				   "Start": "00:00",
				   "Sensitivity": 0,
				   "Units": "mg/dL"
				 }
			       ],
			       "Targets": [
				 {
				   "Start": "00:00",
				   "Low": 100,
				   "High": 100,
				   "Units": "mg/dL"
				 }
			       ]
			     },
			     "After": {
			       "InsulinAction": "3h0m0s",
			       "Ratios": [
				 {
				   "Start": "00:00",
				   "Ratio": 6,
				   "Units": "Grams"
				 }
			       ],
			       "Sensitivities": [
				 {
				   "Start": "00:00",
				   "Sensitivity": 0,
				   "Units": "mg/dL"
				 }
			       ],
			       "Targets": [
				 {
				   "Start": "00:00",
				   "Low": 100,
				   "High": 100,
				   "Units": "mg/dL"
				 }
			       ]
			     }
			   },
			   "Data": "Wg9tRwwJEBURAAA8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIwAAAAAAAAAAAAAAAAAAAGRkAAAAAAAAAAAAAAAAAAAAAAAAAAAAFREAADwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAjAAAAAAAAAAAAAAAAAAAAZGQAAAAAAAAAAAAAAAAAAAAAAAAAAAA0"
			 }`,
			true,
		},
		{
			`{
			   "Type": "BolusWizard",
			   "Time": "2016-07-04T14:25:46-04:00",
			   "BolusWizard": {
			     "GlucoseUnits": "mg/dL",
			     "GlucoseInput": 110,
			     "TargetLow": 99,
			     "TargetHigh": 101,
			     "Sensitivity": 45,
			     "CarbUnits": "Grams",
			     "CarbInput": 5,
			     "CarbRatio": 12,
			     "Unabsorbed": 0.7,
			     "Correction": 0.2,
			     "Food": 0.4,
			     "Bolus": 0.4
			   },
			   "Data": "W25u2Q4EEAVQDC1jAgQAAAcABGU="
			 }`,
			false,
		},
		{
			`{
			   "Type": "BolusWizard",
			   "Time": "2016-07-07T09:59:43-04:00",
			   "BolusWizard": {
			     "GlucoseUnits": "mg/dL",
			     "GlucoseInput": 131,
			     "TargetLow": 100,
			     "TargetHigh": 100,
			     "Sensitivity": 35,
			     "CarbUnits": "Grams",
			     "CarbInput": 40,
			     "CarbRatio": 6,
			     "Unabsorbed": 0,
			     "Correction": 0.8,
			     "Food": 6.6,
			     "Bolus": 7.4
			   },
			   "Data": "W4Nr+wlnEChQADwjZCABCAAAAAEoZA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "BolusWizard",
			   "Time": "2017-02-25T12:41:11-05:00",
			   "BolusWizard": {
			     "GlucoseUnits": "mg/dL",
			     "GlucoseInput": 114,
			     "TargetLow": 100,
			     "TargetHigh": 100,
			     "Sensitivity": 50,
			     "CarbUnits": "Grams",
			     "CarbInput": 10,
			     "CarbRatio": 8,
			     "Unabsorbed": 1.3,
			     "Correction": 0.2,
			     "Food": 1.2,
			     "Bolus": 1.2
			   },
			   "Data": "W3ILqQwZEQpQCDJkAgwAAA0ADGQ="
			 }`,
			false,
		},
		{
			`{
			   "Type": "BolusWizard",
			   "Time": "2017-02-25T12:35:38-05:00",
			   "BolusWizard": {
			     "GlucoseUnits": "μmol/L",
			     "GlucoseInput": 6000,
			     "TargetLow": 5600,
			     "TargetHigh": 5600,
			     "Sensitivity": 2800,
			     "CarbUnits": "Grams",
			     "CarbInput": 10,
			     "CarbRatio": 8,
			     "Unabsorbed": 0,
			     "Correction": 0.1,
			     "Food": 1.2,
			     "Bolus": 1.3
			   },
			   "Data": "WzwmowwZEQqQCBw4AQwAAAAADTg="
			 }`,
			false,
		},
		{
			`{
			   "Type": "BolusWizard",
			   "Time": "2017-02-25T13:09:48-05:00",
			   "BolusWizard": {
			     "GlucoseUnits": "mg/dL",
			     "GlucoseInput": 60,
			     "TargetLow": 100,
			     "TargetHigh": 100,
			     "Sensitivity": 50,
			     "CarbUnits": "Grams",
			     "CarbInput": 10,
			     "CarbRatio": 8,
			     "Unabsorbed": 3.6,
			     "Correction": 24.8,
			     "Food": 1.2,
			     "Bolus": 0.4
			   },
			   "Data": "WzwwiQ0ZEQpQCDJk+AzwACQABGQ="
			 }`,
			false,
		},
		{
			`{
			   "Type": "BolusWizard",
			   "Time": "2017-02-25T13:03:10-05:00",
			   "BolusWizard": {
			     "GlucoseUnits": "μmol/L",
			     "GlucoseInput": 6000,
			     "TargetLow": 5600,
			     "TargetHigh": 5600,
			     "Sensitivity": 2800,
			     "CarbUnits": "Grams",
			     "CarbInput": 10,
			     "CarbRatio": 8,
			     "Unabsorbed": 2.4,
			     "Correction": 0.1,
			     "Food": 1.2,
			     "Bolus": 1.2
			   },
			   "Data": "WzwKgw0ZEQqQCBw4AQwAABgADDg="
			 }`,
			false,
		},
		{
			`{
			   "Type": "BolusWizard",
			   "Time": "2017-02-26T16:10:51-05:00",
			   "BolusWizard": {
			     "GlucoseUnits": "mg/dL",
			     "GlucoseInput": 150,
			     "TargetLow": 90,
			     "TargetHigh": 120,
			     "Sensitivity": 35,
			     "CarbUnits": "Grams",
			     "CarbInput": 20,
			     "CarbRatio": 6,
			     "Unabsorbed": 0,
			     "Correction": 0.8,
			     "Food": 3.3,
			     "Bolus": 4.1
			   },
			   "Data": "W5YzihB6ERRQADwjWiAAhAAAAACkeA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "BolusWizard",
			   "Time": "2017-02-28T20:03:55-05:00",
			   "Data": "W7w3gxQcEQBgGSNkGQAAAAAAGWQ=",
			   "BolusWizard": {
			     "CarbRatio": 2.5,
			     "GlucoseUnits": "mg/dL",
			     "GlucoseInput": 188,
			     "TargetLow": 100,
			     "TargetHigh": 100,
			     "Sensitivity": 35,
			     "CarbUnits": "Exchanges",
			     "CarbInput": 0,
			     "Correction": 0.9,
			     "Food": 0,
			     "Unabsorbed": 0,
			     "Bolus": 2.5
			   }
			 }`,
			false,
		},
		{
			`{
			   "Type": "BolusWizard",
			   "Time": "2017-02-26T17:27:03-05:00",
			   "BolusWizard": {
			     "GlucoseUnits": "mg/dL",
			     "GlucoseInput": 150,
			     "TargetLow": 90,
			     "TargetHigh": 120,
			     "Sensitivity": 35,
			     "CarbUnits": "Exchanges",
			     "CarbInput": 30,
			     "CarbRatio": 2.5,
			     "Unabsorbed": 13.1,
			     "Correction": 0.8,
			     "Food": 7.5,
			     "Bolus": 7.5
			   },
			   "Data": "W5YDmxF6ER5gCcQjWiABLAACDAEseA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "BolusWizard",
			   "Time": "2017-02-26T17:13:16-05:00",
			   "BolusWizard": {
			     "GlucoseUnits": "μmol/L",
			     "GlucoseInput": 8300,
			     "TargetLow": 5000,
			     "TargetHigh": 6700,
			     "Sensitivity": 1900,
			     "CarbUnits": "Grams",
			     "CarbInput": 20,
			     "CarbRatio": 6,
			     "Unabsorbed": 2.9,
			     "Correction": 0.8,
			     "Food": 3.3,
			     "Bolus": 3.3
			   },
			   "Data": "W1MQjRF6ERSQADwTMiAAhAAAdACEQw=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "BolusWizard",
			   "Time": "2017-02-26T17:18:47-05:00",
			   "BolusWizard": {
			     "GlucoseUnits": "μmol/L",
			     "GlucoseInput": 8300,
			     "TargetLow": 5000,
			     "TargetHigh": 6700,
			     "Sensitivity": 1900,
			     "CarbUnits": "Exchanges",
			     "CarbInput": 30,
			     "CarbRatio": 2.5,
			     "Unabsorbed": 6,
			     "Correction": 0.8,
			     "Food": 7.5,
			     "Bolus": 7.5
			   },
			   "Data": "W1MvkhF6ER6gCcQTMiABLAAA8AEsQw=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "UnabsorbedInsulin",
			   "UnabsorbedInsulin": [
			     {
			       "Age": "1h0m0s",
			       "Bolus": 1.15
			     },
			     {
			       "Age": "1h10m0s",
			       "Bolus": 0.85
			     },
			     {
			       "Age": "2h30m0s",
			       "Bolus": 1
			     },
			     {
			       "Age": "3h40m0s",
			       "Bolus": 0.4
			     },
			     {
			       "Age": "3h50m0s",
			       "Bolus": 0.6
			     },
			     {
			       "Age": "4h0m0s",
			       "Bolus": 0.55
			     },
			     {
			       "Age": "4h10m0s",
			       "Bolus": 0.55
			     },
			     {
			       "Age": "4m0s",
			       "Bolus": 0.55
			     },
			     {
			       "Age": "14m0s",
			       "Bolus": 0.55
			     },
			     {
			       "Age": "24m0s",
			       "Bolus": 0.55
			     },
			     {
			       "Age": "34m0s",
			       "Bolus": 0.55
			     },
			     {
			       "Age": "44m0s",
			       "Bolus": 0.55
			     },
			     {
			       "Age": "54m0s",
			       "Bolus": 5.15
			     },
			     {
			       "Age": "1h54m0s",
			       "Bolus": 0.45
			     },
			     {
			       "Age": "2h4m0s",
			       "Bolus": 0.55
			     },
			     {
			       "Age": "2h34m0s",
			       "Bolus": 0.75
			     },
			     {
			       "Age": "2h44m0s",
			       "Bolus": 1.45
			     }
			   ],
			   "Data": "XDUuPMAiRsAolsAQ3MAY5sAW8MAW+sAWBNAWDtAWGNAWItAWLNDONtASctAWfNAemtA6pNA="
			 }`,
			true,
		},
		{
			`{
			   "Type": "EnableVariableBolus",
			   "Time": "2016-07-04T14:00:00-04:00",
			   "Enabled": false,
			   "Data": "XgBAwA4EEA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "ChangeEasyBolus",
			   "Time": "2016-07-04T14:00:22-04:00",
			   "Data": "XxRWwA4EEA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "EnableBgReminder",
			   "Time": "2016-07-05T15:43:10-04:00",
			   "Enabled": false,
			   "Data": "YABK6w8FEA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "EnableAlarmClock",
			   "Time": "2016-07-05T13:09:47-04:00",
			   "Enabled": false,
			   "Data": "YQBvyQ0FEA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "ChangeTempBasalType",
			   "Time": "2016-06-04T18:52:10-04:00",
			   "TempBasalType": "Absolute",
			   "Data": "YgBKtBIEEA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "ChangeTempBasalType",
			   "Time": "2016-06-27T13:18:45-04:00",
			   "TempBasalType": "Percent",
			   "Data": "YgFtkg0bEA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "ChangeAlarmType",
			   "Time": "2016-07-04T14:13:37-04:00",
			   "Value": 4,
			   "Data": "YwRlzQ4EEA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "ChangeTimeFormat",
			   "Time": "2005-01-01T00:01:44-05:00",
			   "Value": 0,
			   "Data": "ZAAsQQABBQ=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "ChangeReservoirWarning",
			   "Time": "2016-07-05T13:07:10-04:00",
			   "Insulin": 20,
			   "Data": "ZVBKxw0FEA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "ChangeReservoirWarning",
			   "Time": "2016-07-05T13:07:03-04:00",
			   "Duration": "8h0m0s",
			   "Data": "ZUFDxw0FEA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "EnableBolusReminder",
			   "Time": "2016-07-05T15:43:31-04:00",
			   "Enabled": true,
			   "Data": "ZgFf6w8FEA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "SetBolusReminderTime",
			   "Time": "2016-07-05T15:43:31-04:00",
			   "Enabled": true,
			   "Data": "ZwFf6w8FEAIA"
			 }`,
			true,
		},
		{
			`{
			   "Type": "DeleteBolusReminderTime",
			   "Time": "2016-07-05T15:43:38-04:00",
			   "Enabled": true,
			   "Data": "aAFm6w8FEAIA"
			 }`,
			true,
		},
		{
			`{
			   "Type": "DeleteAlarmClockTime",
			   "Time": "2016-07-04T14:15:40-04:00",
			   "Data": "agFozw4EEA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "DailyTotal522",
			   "Time": "2016-07-03T00:00:00-04:00",
			   "Data": "bWOQBQwA6AAAAAAB2AHYZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMAOgAAAA="
			 }`,
			false,
		},
		{
			`{
			   "Type": "DailyTotal523",
			   "Time": "2016-07-04T00:00:00-04:00",
			   "Data": "bmSQBQEMAAAIAAAKnAQMJgaQPgBLAQgESAFAAAABBgEABAAAAAAAAAAAtHsAAAAAAAAAAA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "BasalProfileStart",
			   "Time": "2016-06-24T10:00:00-04:00",
			   "BasalProfileStart": {
			     "ProfileIndex": 3,
			     "BasalRate": {
			       "Start": "10:00",
			       "Rate": 1.2
			     }
			   },
			   "Data": "ewNAgAoYEBQwAA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "ConnectOtherDevices",
			   "Time": "2016-06-02T11:37:50-04:00",
			   "Enabled": true,
			   "Data": "fAFypQsCEA=="
			 }`,
			true,
		},
		{
			`{
			   "Type": "ChangeOtherDevice",
			   "Time": "2016-07-05T13:12:11-04:00",
			   "Data": "gQFLzA0FEAARERER"
			 }`,
			true,
		},
		{
			`{
			   "Type": "ChangeMarriage",
			   "Time": "2016-07-05T13:12:11-04:00",
			   "Data": "gQFLzA0FEAARERER"
			 }`,
			true,
		},
		{
			`{
			   "Type": "DeleteOtherDevice",
			   "Time": "2016-07-05T13:12:19-04:00",
			   "Data": "ggFTzA0FEAARERER"
			 }`,
			true,
		},
	}
	for _, c := range cases {
		var r1, r2 HistoryRecord
		err := json.Unmarshal([]byte(c.test), &r1)
		if err != nil {
			t.Errorf("json.Unmarshal(%s) returned %v", c.test, err)
			continue
		}
		r2, err = DecodeHistoryRecord(r1.Data, c.newer)
		if err != nil {
			t.Errorf("DecodeHistoryRecord(%X, %v) returned %v, want %v", r1.Data, c.newer, err, r1)
			continue
		}
		if !reflect.DeepEqual(r1, r2) {
			t.Errorf("DecodeHistoryRecord(%X, %v) == %v, want %v", r1.Data, c.newer, r2, r1)
		}
	}
}

func (r HistoryRecord) String() string {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func TestDecodeHistoryRecords(t *testing.T) {
	cases := []struct {
		page    string
		results string
		newer   bool
	}{
		{
			"6E 21 90 05 00 00 00 00 00 00 00 02 BE 02 BE 64 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 7B 01 00 DE 08 02 10 11 22 00 7B 02 00 C0 16 02 10 2C 1C 00 7B 00 00 C0 00 03 10 00 16 00 07 00 00 02 BE 22 90 00 00 00 6E 22 90 05 00 00 00 00 00 00 00 02 BE 02 BE 64 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 7B 01 00 DE 08 03 10 11 22 00 34 64 0E DC 12 03 10 81 01 23 EF 12 03 10 00 10 11 11 11 7D 02 23 EF 12 03 10 00 A2 CE 8A A0 00 10 11 11 11 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 7B 02 00 C0 16 03 10 2C 1C 00 21 00 1C CE 16 03 10 03 00 00 00 10 0F D0 36 03 10 7B 02 1E DB 16 03 10 2C 1C 00 03 00 01 00 01 1C DB 16 03 10 82 01 06 DC 16 03 10 00 A2 CE 8A A0 82 01 08 DC 16 03 10 00 10 11 11 11 7B 00 00 C0 00 04 10 00 16 00 07 00 00 02 B8 23 90 00 00 00 6E 23 90 05 00 00 00 00 00 00 00 02 B8 02 B8 64 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 7B 01 00 DE 08 04 10 11 22 00 7B 02 00 C0 16 04 10 2C 1C 00 7B 00 00 C0 00 05 10 00 16 00 07 00 00 02 BE 24 90 00 00 00 6E 24 90 05 00 00 00 00 00 00 00 02 BE 02 BE 64 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 7B 01 00 DE 08 05 10 11 22 00 7B 02 00 C0 16 05 10 2C 1C 00 7B 00 00 C0 00 06 10 00 16 00 07 00 00 02 BE 25 90 00 00 00 6E 25 90 05 00 00 00 00 00 00 00 02 BE 02 BE 64 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 7B 01 00 DE 08 06 10 11 22 00 7B 02 00 C0 16 06 10 2C 1C 00 7B 00 00 C0 00 07 10 00 16 00 07 00 00 02 BE 26 90 00 00 00 6E 26 90 05 00 00 00 00 00 00 00 02 BE 02 BE 64 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 7B 01 00 DE 08 07 10 11 22 00 81 01 0B EC 0A 07 10 00 A2 CE 8A A0 7D 01 0B EC 0A 07 10 00 A2 CE 8A A0 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 82 01 08 C4 0B 07 10 00 A2 CE 8A A0 81 01 0C C4 0B 07 10 00 A2 CE 8A A0 7D 01 0C C4 0B 07 10 00 A2 CE 8A A0 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00",
			`[
			   {
			     "Type": "ChangeOtherDevice",
			     "Time": "2016-03-07T11:04:12-05:00",
			     "Data": "fQEMxAsHEACizoqgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="
			   },
			   {
			     "Type": "ChangeMarriage",
			     "Time": "2016-03-07T11:04:12-05:00",
			     "Data": "gQEMxAsHEACizoqg"
			   },
			   {
			     "Type": "DeleteOtherDevice",
			     "Time": "2016-03-07T11:04:08-05:00",
			     "Data": "ggEIxAsHEACizoqg"
			   },
			   {
			     "Type": "ChangeOtherDevice",
			     "Time": "2016-03-07T10:44:11-05:00",
			     "Data": "fQEL7AoHEACizoqgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="
			   },
			   {
			     "Type": "ChangeMarriage",
			     "Time": "2016-03-07T10:44:11-05:00",
			     "Data": "gQEL7AoHEACizoqg"
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-03-07T08:30:00-05:00",
			     "Data": "ewEA3ggHEBEiAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 1,
			       "BasalRate": {
				 "Start": "08:30",
				 "Rate": 0.85
			       }
			     }
			   },
			   {
			     "Type": "DailyTotal523",
			     "Time": "2016-03-06T00:00:00-05:00",
			     "Data": "biaQBQAAAAAAAAACvgK+ZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="
			   },
			   {
			     "Type": "DailyTotal",
			     "Time": "2016-03-06T00:00:00-05:00",
			     "Data": "BwAAAr4mkAAAAA==",
			     "Insulin": 17.55
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-03-07T00:00:00-05:00",
			     "Data": "ewAAwAAHEAAWAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 0,
			       "BasalRate": {
				 "Start": "00:00",
				 "Rate": 0.55
			       }
			     }
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-03-06T22:00:00-05:00",
			     "Data": "ewIAwBYGECwcAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 2,
			       "BasalRate": {
				 "Start": "22:00",
				 "Rate": 0.7
			       }
			     }
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-03-06T08:30:00-05:00",
			     "Data": "ewEA3ggGEBEiAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 1,
			       "BasalRate": {
				 "Start": "08:30",
				 "Rate": 0.85
			       }
			     }
			   },
			   {
			     "Type": "DailyTotal523",
			     "Time": "2016-03-05T00:00:00-05:00",
			     "Data": "biWQBQAAAAAAAAACvgK+ZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="
			   },
			   {
			     "Type": "DailyTotal",
			     "Time": "2016-03-05T00:00:00-05:00",
			     "Data": "BwAAAr4lkAAAAA==",
			     "Insulin": 17.55
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-03-06T00:00:00-05:00",
			     "Data": "ewAAwAAGEAAWAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 0,
			       "BasalRate": {
				 "Start": "00:00",
				 "Rate": 0.55
			       }
			     }
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-03-05T22:00:00-05:00",
			     "Data": "ewIAwBYFECwcAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 2,
			       "BasalRate": {
				 "Start": "22:00",
				 "Rate": 0.7
			       }
			     }
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-03-05T08:30:00-05:00",
			     "Data": "ewEA3ggFEBEiAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 1,
			       "BasalRate": {
				 "Start": "08:30",
				 "Rate": 0.85
			       }
			     }
			   },
			   {
			     "Type": "DailyTotal523",
			     "Time": "2016-03-04T00:00:00-05:00",
			     "Data": "biSQBQAAAAAAAAACvgK+ZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="
			   },
			   {
			     "Type": "DailyTotal",
			     "Time": "2016-03-04T00:00:00-05:00",
			     "Data": "BwAAAr4kkAAAAA==",
			     "Insulin": 17.55
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-03-05T00:00:00-05:00",
			     "Data": "ewAAwAAFEAAWAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 0,
			       "BasalRate": {
				 "Start": "00:00",
				 "Rate": 0.55
			       }
			     }
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-03-04T22:00:00-05:00",
			     "Data": "ewIAwBYEECwcAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 2,
			       "BasalRate": {
				 "Start": "22:00",
				 "Rate": 0.7
			       }
			     }
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-03-04T08:30:00-05:00",
			     "Data": "ewEA3ggEEBEiAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 1,
			       "BasalRate": {
				 "Start": "08:30",
				 "Rate": 0.85
			       }
			     }
			   },
			   {
			     "Type": "DailyTotal523",
			     "Time": "2016-03-03T00:00:00-05:00",
			     "Data": "biOQBQAAAAAAAAACuAK4ZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="
			   },
			   {
			     "Type": "DailyTotal",
			     "Time": "2016-03-03T00:00:00-05:00",
			     "Data": "BwAAArgjkAAAAA==",
			     "Insulin": 17.4
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-03-04T00:00:00-05:00",
			     "Data": "ewAAwAAEEAAWAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 0,
			       "BasalRate": {
				 "Start": "00:00",
				 "Rate": 0.55
			       }
			     }
			   },
			   {
			     "Type": "DeleteOtherDevice",
			     "Time": "2016-03-03T22:28:08-05:00",
			     "Data": "ggEI3BYDEAAQERER"
			   },
			   {
			     "Type": "DeleteOtherDevice",
			     "Time": "2016-03-03T22:28:06-05:00",
			     "Data": "ggEG3BYDEACizoqg"
			   },
			   {
			     "Type": "Prime",
			     "Time": "2016-03-03T22:27:28-05:00",
			     "Data": "AwABAAEc2xYDEA==",
			     "Prime": {
			       "Fixed": 0.1,
			       "Manual": 0.1
			     }
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-03-03T22:27:30-05:00",
			     "Data": "ewIe2xYDECwcAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 2,
			       "BasalRate": {
				 "Start": "22:00",
				 "Rate": 0.7
			       }
			     }
			   },
			   {
			     "Type": "Prime",
			     "Time": "2016-03-03T22:16:15-05:00",
			     "Data": "AwAAABAP0DYDEA==",
			     "Prime": {
			       "Fixed": 0,
			       "Manual": 1.6
			     }
			   },
			   {
			     "Type": "Rewind",
			     "Time": "2016-03-03T22:14:28-05:00",
			     "Data": "IQAczhYDEA=="
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-03-03T22:00:00-05:00",
			     "Data": "ewIAwBYDECwcAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 2,
			       "BasalRate": {
				 "Start": "22:00",
				 "Rate": 0.7
			       }
			     }
			   },
			   {
			     "Type": "ChangeOtherDevice",
			     "Time": "2016-03-03T18:47:35-05:00",
			     "Data": "fQIj7xIDEACizoqgABAREREAAAAAAAAAAAAAAAAAAAAAAAAAAA=="
			   },
			   {
			     "Type": "ChangeMarriage",
			     "Time": "2016-03-03T18:47:35-05:00",
			     "Data": "gQEj7xIDEAAQERER"
			   },
			   {
			     "Type": "LowReservoir",
			     "Time": "2016-03-03T18:28:14-05:00",
			     "Data": "NGQO3BIDEA==",
			     "Insulin": 10
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-03-03T08:30:00-05:00",
			     "Data": "ewEA3ggDEBEiAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 1,
			       "BasalRate": {
				 "Start": "08:30",
				 "Rate": 0.85
			       }
			     }
			   },
			   {
			     "Type": "DailyTotal523",
			     "Time": "2016-03-02T00:00:00-05:00",
			     "Data": "biKQBQAAAAAAAAACvgK+ZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="
			   },
			   {
			     "Type": "DailyTotal",
			     "Time": "2016-03-02T00:00:00-05:00",
			     "Data": "BwAAAr4ikAAAAA==",
			     "Insulin": 17.55
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-03-03T00:00:00-05:00",
			     "Data": "ewAAwAADEAAWAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 0,
			       "BasalRate": {
				 "Start": "00:00",
				 "Rate": 0.55
			       }
			     }
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-03-02T22:00:00-05:00",
			     "Data": "ewIAwBYCECwcAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 2,
			       "BasalRate": {
				 "Start": "22:00",
				 "Rate": 0.7
			       }
			     }
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-03-02T08:30:00-05:00",
			     "Data": "ewEA3ggCEBEiAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 1,
			       "BasalRate": {
				 "Start": "08:30",
				 "Rate": 0.85
			       }
			     }
			   },
			   {
			     "Type": "DailyTotal523",
			     "Time": "2016-03-01T00:00:00-05:00",
			     "Data": "biGQBQAAAAAAAAACvgK+ZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="
			   }
			 ]`,
			true,
		},
		{
			"01 00 52 00 52 00 00 00 42 22 54 65 10 5B 00 44 3B 14 65 10 0E 50 00 78 4B 50 00 00 2E 00 00 00 00 2E 78 5C 0E 1A 18 C0 38 22 C0 16 58 D0 02 DA D0 01 00 2E 00 2E 00 4C 00 44 3B 54 65 10 0A 33 61 02 36 05 90 5B 33 72 02 16 65 10 00 51 00 B4 55 50 58 00 00 00 00 42 00 16 78 5C 0E 2E 43 C0 1A 57 C0 38 61 C0 16 97 D0 01 00 16 00 16 00 42 00 72 02 56 65 10 7B 00 40 00 00 06 10 00 10 00 07 00 00 02 CB 45 10 00 00 00 6E 45 10 05 00 E7 00 00 02 00 00 02 CB 01 51 2F 01 7A 35 00 7B 00 D6 00 16 00 8E 00 00 04 01 01 00 04 00 00 00 00 00 00 00 00 9B 33 00 00 00 00 00 00 00 00 7B 01 40 00 01 06 10 02 0C 00 0A CA 48 09 21 06 10 5B CA 4B 09 01 06 10 00 50 00 C8 50 50 28 00 00 00 00 00 00 28 78 5C 0E 16 C2 C0 2E FE C0 1A 12 D0 38 1C D0 01 00 28 00 28 00 00 00 4C 09 41 06 10 7B 02 40 00 04 06 10 08 0D 00 7B 03 40 00 06 06 10 0C 10 00 5B 00 6A 13 09 66 10 14 50 00 6E 4B 50 00 00 48 00 00 00 00 48 78 01 00 48 00 48 00 00 00 6A 13 49 66 10 7B 04 40 00 0A 06 10 14 0B 00 0A FA 41 36 2B 06 10 5B FA 67 36 0B 66 10 47 50 00 B4 4B 50 44 00 9C 00 00 08 00 D8 78 5C 05 48 9F C0 01 00 68 00 68 00 08 04 5D 38 AB 66 10 01 00 70 00 70 00 08 00 67 36 8B 66 10 7B 05 40 00 0C 06 10 18 0A 00 7B 06 40 00 10 06 10 20 0E 00 5B 00 43 19 12 66 10 3C 50 00 78 4B 50 00 00 C8 00 00 00 00 C8 78 5C 2C 02 0E D0 08 18 D0 08 22 D0 0A 2C D0 08 36 D0 08 40 D0 0A 4A D0 08 54 D0 08 5E D0 0A 68 D0 08 72 D0 08 7C D0 62 86 D0 16 90 D0 01 00 64 00 1C 00 00 06 6B 1A B2 66 10 01 00 64 00 64 00 00 00 43 19 92 66 10 7B 07 40 00 13 06 10 26 10 00 1E 01 65 12 13 06 10 7B 07 6E 12 13 06 10 26 10 00 1F 20 6E 12 13 06 10 21 00 74 12 13 06 10 1A 00 62 14 13 06 10 1A 01 77 14 13 06 10 03 00 00 00 43 6A 15 33 06 10 7B 07 66 1B 13 06 10 26 10 00 03 00 03 00 03 5C 1B 13 06 10 7B 08 40 1E 14 06 10 29 15 00 0A 82 5B 10 35 06 90 5B 82 60 10 15 06 10 00 51 00 B4 55 50 7C 00 00 00 00 0C 00 70 78 5C 20 02 79 C0 04 83 C0 06 8D C0 06 97 C0 06 A1 C0 68 AB C0 02 B9 D0 08 C3 D0 08 CD D0 0A D7 D0 5B 82 62 10 15 06 10 00 51 00 B4 55 50 7C 00 00 00 00 0C 00 70 78 5C 20 02 79 C0 04 83 C0 06 8D C0 06 97 C0 06 A1 C0 68 AB C0 02 B9 D0 08 C3 D0 08 CD D0 0A D7 D0 01 00 70 00 70 00 0C 00 62 10 55 06 10",
			`[
			   {
			     "Type": "Bolus",
			     "Time": "2016-04-06T21:16:34-04:00",
			     "Data": "AQBwAHAADABiEFUGEA==",
			     "Bolus": {
			       "Duration": "0s",
			       "Programmed": 2.8,
			       "Amount": 2.8,
			       "Unabsorbed": 0.3
			     }
			   },
			   {
			     "Type": "UnabsorbedInsulin",
			     "Data": "XCACecAEg8AGjcAGl8AGocBoq8ACudAIw9AIzdAK19A=",
			     "UnabsorbedInsulin": [
			       {
				 "Age": "2h1m0s",
				 "Bolus": 0.05
			       },
			       {
				 "Age": "2h11m0s",
				 "Bolus": 0.1
			       },
			       {
				 "Age": "2h21m0s",
				 "Bolus": 0.15
			       },
			       {
				 "Age": "2h31m0s",
				 "Bolus": 0.15
			       },
			       {
				 "Age": "2h41m0s",
				 "Bolus": 0.15
			       },
			       {
				 "Age": "2h51m0s",
				 "Bolus": 2.6
			       },
			       {
				 "Age": "3h5m0s",
				 "Bolus": 0.05
			       },
			       {
				 "Age": "3h15m0s",
				 "Bolus": 0.2
			       },
			       {
				 "Age": "3h25m0s",
				 "Bolus": 0.2
			       },
			       {
				 "Age": "3h35m0s",
				 "Bolus": 0.25
			       }
			     ]
			   },
			   {
			     "Type": "BolusWizard",
			     "Time": "2016-04-06T21:16:34-04:00",
			     "Data": "W4JiEBUGEABRALRVUHwAAAAADABweA==",
			     "BolusWizard": {
			       "GlucoseUnits": "mg/dL",
			       "GlucoseInput": 386,
			       "TargetLow": 80,
			       "TargetHigh": 120,
			       "Sensitivity": 85,
			       "CarbUnits": "Grams",
			       "CarbInput": 0,
			       "CarbRatio": 18,
			       "Unabsorbed": 0.3,
			       "Correction": 3.1,
			       "Food": 0,
			       "Bolus": 2.8
			     }
			   },
			   {
			     "Type": "UnabsorbedInsulin",
			     "Data": "XCACecAEg8AGjcAGl8AGocBoq8ACudAIw9AIzdAK19A=",
			     "UnabsorbedInsulin": [
			       {
				 "Age": "2h1m0s",
				 "Bolus": 0.05
			       },
			       {
				 "Age": "2h11m0s",
				 "Bolus": 0.1
			       },
			       {
				 "Age": "2h21m0s",
				 "Bolus": 0.15
			       },
			       {
				 "Age": "2h31m0s",
				 "Bolus": 0.15
			       },
			       {
				 "Age": "2h41m0s",
				 "Bolus": 0.15
			       },
			       {
				 "Age": "2h51m0s",
				 "Bolus": 2.6
			       },
			       {
				 "Age": "3h5m0s",
				 "Bolus": 0.05
			       },
			       {
				 "Age": "3h15m0s",
				 "Bolus": 0.2
			       },
			       {
				 "Age": "3h25m0s",
				 "Bolus": 0.2
			       },
			       {
				 "Age": "3h35m0s",
				 "Bolus": 0.25
			       }
			     ]
			   },
			   {
			     "Type": "BolusWizard",
			     "Time": "2016-04-06T21:16:32-04:00",
			     "Data": "W4JgEBUGEABRALRVUHwAAAAADABweA==",
			     "BolusWizard": {
			       "GlucoseUnits": "mg/dL",
			       "GlucoseInput": 386,
			       "TargetLow": 80,
			       "TargetHigh": 120,
			       "Sensitivity": 85,
			       "CarbUnits": "Grams",
			       "CarbInput": 0,
			       "CarbRatio": 18,
			       "Unabsorbed": 0.3,
			       "Correction": 3.1,
			       "Food": 0,
			       "Bolus": 2.8
			     }
			   },
			   {
			     "Type": "BGCapture",
			     "Time": "2016-04-06T21:16:27-04:00",
			     "Data": "CoJbEDUGkA==",
			     "Glucose": 386,
			     "GlucoseUnits": "mg/dL"
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-04-06T20:30:00-04:00",
			     "Data": "ewhAHhQGECkVAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 8,
			       "BasalRate": {
				 "Start": "20:30",
				 "Rate": 0.525
			       }
			     }
			   },
			   {
			     "Type": "Prime",
			     "Time": "2016-04-06T19:27:28-04:00",
			     "Data": "AwADAANcGxMGEA==",
			     "Prime": {
			       "Fixed": 0.3,
			       "Manual": 0.3
			     }
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-04-06T19:27:38-04:00",
			     "Data": "ewdmGxMGECYQAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 7,
			       "BasalRate": {
				 "Start": "19:00",
				 "Rate": 0.4
			       }
			     }
			   },
			   {
			     "Type": "Prime",
			     "Time": "2016-04-06T19:21:42-04:00",
			     "Data": "AwAAAENqFTMGEA==",
			     "Prime": {
			       "Fixed": 0,
			       "Manual": 6.7
			     }
			   },
			   {
			     "Type": "BatteryChange",
			     "Time": "2016-04-06T19:20:55-04:00",
			     "Data": "GgF3FBMGEA=="
			   },
			   {
			     "Type": "BatteryChange",
			     "Time": "2016-04-06T19:20:34-04:00",
			     "Data": "GgBiFBMGEA=="
			   },
			   {
			     "Type": "Rewind",
			     "Time": "2016-04-06T19:18:52-04:00",
			     "Data": "IQB0EhMGEA=="
			   },
			   {
			     "Type": "ResumePump",
			     "Time": "2016-04-06T19:18:46-04:00",
			     "Data": "HyBuEhMGEA=="
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-04-06T19:18:46-04:00",
			     "Data": "ewduEhMGECYQAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 7,
			       "BasalRate": {
				 "Start": "19:00",
				 "Rate": 0.4
			       }
			     }
			   },
			   {
			     "Type": "SuspendPump",
			     "Time": "2016-04-06T19:18:37-04:00",
			     "Data": "HgFlEhMGEA=="
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-04-06T19:00:00-04:00",
			     "Data": "ewdAABMGECYQAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 7,
			       "BasalRate": {
				 "Start": "19:00",
				 "Rate": 0.4
			       }
			     }
			   },
			   {
			     "Type": "Bolus",
			     "Time": "2016-04-06T18:25:03-04:00",
			     "Data": "AQBkAGQAAABDGZJmEA==",
			     "Bolus": {
			       "Duration": "0s",
			       "Programmed": 2.5,
			       "Amount": 2.5,
			       "Unabsorbed": 0
			     }
			   },
			   {
			     "Type": "Bolus",
			     "Time": "2016-04-06T18:26:43-04:00",
			     "Data": "AQBkABwAAAZrGrJmEA==",
			     "Bolus": {
			       "Duration": "3h0m0s",
			       "Programmed": 2.5,
			       "Amount": 0.7,
			       "Unabsorbed": 0
			     }
			   },
			   {
			     "Type": "UnabsorbedInsulin",
			     "Data": "XCwCDtAIGNAIItAKLNAINtAIQNAKStAIVNAIXtAKaNAIctAIfNBihtAWkNA=",
			     "UnabsorbedInsulin": [
			       {
				 "Age": "14m0s",
				 "Bolus": 0.05
			       },
			       {
				 "Age": "24m0s",
				 "Bolus": 0.2
			       },
			       {
				 "Age": "34m0s",
				 "Bolus": 0.2
			       },
			       {
				 "Age": "44m0s",
				 "Bolus": 0.25
			       },
			       {
				 "Age": "54m0s",
				 "Bolus": 0.2
			       },
			       {
				 "Age": "1h4m0s",
				 "Bolus": 0.2
			       },
			       {
				 "Age": "1h14m0s",
				 "Bolus": 0.25
			       },
			       {
				 "Age": "1h24m0s",
				 "Bolus": 0.2
			       },
			       {
				 "Age": "1h34m0s",
				 "Bolus": 0.2
			       },
			       {
				 "Age": "1h44m0s",
				 "Bolus": 0.25
			       },
			       {
				 "Age": "1h54m0s",
				 "Bolus": 0.2
			       },
			       {
				 "Age": "2h4m0s",
				 "Bolus": 0.2
			       },
			       {
				 "Age": "2h14m0s",
				 "Bolus": 2.45
			       },
			       {
				 "Age": "2h24m0s",
				 "Bolus": 0.55
			       }
			     ]
			   },
			   {
			     "Type": "BolusWizard",
			     "Time": "2016-04-06T18:25:03-04:00",
			     "Data": "WwBDGRJmEDxQAHhLUAAAyAAAAADIeA==",
			     "BolusWizard": {
			       "GlucoseUnits": "mg/dL",
			       "GlucoseInput": 0,
			       "TargetLow": 80,
			       "TargetHigh": 120,
			       "Sensitivity": 75,
			       "CarbUnits": "Grams",
			       "CarbInput": 60,
			       "CarbRatio": 12,
			       "Unabsorbed": 0,
			       "Correction": 0,
			       "Food": 5,
			       "Bolus": 5
			     }
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-04-06T16:00:00-04:00",
			     "Data": "ewZAABAGECAOAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 6,
			       "BasalRate": {
				 "Start": "16:00",
				 "Rate": 0.35
			       }
			     }
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-04-06T12:00:00-04:00",
			     "Data": "ewVAAAwGEBgKAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 5,
			       "BasalRate": {
				 "Start": "12:00",
				 "Rate": 0.25
			       }
			     }
			   },
			   {
			     "Type": "Bolus",
			     "Time": "2016-04-06T11:54:39-04:00",
			     "Data": "AQBwAHAACABnNotmEA==",
			     "Bolus": {
			       "Duration": "0s",
			       "Programmed": 2.8,
			       "Amount": 2.8,
			       "Unabsorbed": 0.2
			     }
			   },
			   {
			     "Type": "Bolus",
			     "Time": "2016-04-06T11:56:29-04:00",
			     "Data": "AQBoAGgACARdOKtmEA==",
			     "Bolus": {
			       "Duration": "2h0m0s",
			       "Programmed": 2.6,
			       "Amount": 2.6,
			       "Unabsorbed": 0.2
			     }
			   },
			   {
			     "Type": "UnabsorbedInsulin",
			     "Data": "XAVIn8A=",
			     "UnabsorbedInsulin": [
			       {
				 "Age": "2h39m0s",
				 "Bolus": 1.8
			       }
			     ]
			   },
			   {
			     "Type": "BolusWizard",
			     "Time": "2016-04-06T11:54:39-04:00",
			     "Data": "W/pnNgtmEEdQALRLUEQAnAAACADYeA==",
			     "BolusWizard": {
			       "GlucoseUnits": "mg/dL",
			       "GlucoseInput": 250,
			       "TargetLow": 80,
			       "TargetHigh": 120,
			       "Sensitivity": 75,
			       "CarbUnits": "Grams",
			       "CarbInput": 71,
			       "CarbRatio": 18,
			       "Unabsorbed": 0.2,
			       "Correction": 1.7,
			       "Food": 3.9,
			       "Bolus": 5.4
			     }
			   },
			   {
			     "Type": "BGCapture",
			     "Time": "2016-04-06T11:54:01-04:00",
			     "Data": "CvpBNisGEA==",
			     "Glucose": 250,
			     "GlucoseUnits": "mg/dL"
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-04-06T10:00:00-04:00",
			     "Data": "ewRAAAoGEBQLAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 4,
			       "BasalRate": {
				 "Start": "10:00",
				 "Rate": 0.275
			       }
			     }
			   },
			   {
			     "Type": "Bolus",
			     "Time": "2016-04-06T09:19:42-04:00",
			     "Data": "AQBIAEgAAABqE0lmEA==",
			     "Bolus": {
			       "Duration": "0s",
			       "Programmed": 1.8,
			       "Amount": 1.8,
			       "Unabsorbed": 0
			     }
			   },
			   {
			     "Type": "BolusWizard",
			     "Time": "2016-04-06T09:19:42-04:00",
			     "Data": "WwBqEwlmEBRQAG5LUAAASAAAAABIeA==",
			     "BolusWizard": {
			       "GlucoseUnits": "mg/dL",
			       "GlucoseInput": 0,
			       "TargetLow": 80,
			       "TargetHigh": 120,
			       "Sensitivity": 75,
			       "CarbUnits": "Grams",
			       "CarbInput": 20,
			       "CarbRatio": 11,
			       "Unabsorbed": 0,
			       "Correction": 0,
			       "Food": 1.8,
			       "Bolus": 1.8
			     }
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-04-06T06:00:00-04:00",
			     "Data": "ewNAAAYGEAwQAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 3,
			       "BasalRate": {
				 "Start": "06:00",
				 "Rate": 0.4
			       }
			     }
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-04-06T04:00:00-04:00",
			     "Data": "ewJAAAQGEAgNAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 2,
			       "BasalRate": {
				 "Start": "04:00",
				 "Rate": 0.325
			       }
			     }
			   },
			   {
			     "Type": "Bolus",
			     "Time": "2016-04-06T01:09:12-04:00",
			     "Data": "AQAoACgAAABMCUEGEA==",
			     "Bolus": {
			       "Duration": "0s",
			       "Programmed": 1,
			       "Amount": 1,
			       "Unabsorbed": 0
			     }
			   },
			   {
			     "Type": "UnabsorbedInsulin",
			     "Data": "XA4WwsAu/sAaEtA4HNA=",
			     "UnabsorbedInsulin": [
			       {
				 "Age": "3h14m0s",
				 "Bolus": 0.55
			       },
			       {
				 "Age": "4h14m0s",
				 "Bolus": 1.15
			       },
			       {
				 "Age": "18m0s",
				 "Bolus": 0.65
			       },
			       {
				 "Age": "28m0s",
				 "Bolus": 1.4
			       }
			     ]
			   },
			   {
			     "Type": "BolusWizard",
			     "Time": "2016-04-06T01:09:11-04:00",
			     "Data": "W8pLCQEGEABQAMhQUCgAAAAAAAAoeA==",
			     "BolusWizard": {
			       "GlucoseUnits": "mg/dL",
			       "GlucoseInput": 202,
			       "TargetLow": 80,
			       "TargetHigh": 120,
			       "Sensitivity": 80,
			       "CarbUnits": "Grams",
			       "CarbInput": 0,
			       "CarbRatio": 20,
			       "Unabsorbed": 0,
			       "Correction": 1,
			       "Food": 0,
			       "Bolus": 1
			     }
			   },
			   {
			     "Type": "BGCapture",
			     "Time": "2016-04-06T01:09:08-04:00",
			     "Data": "CspICSEGEA==",
			     "Glucose": 202,
			     "GlucoseUnits": "mg/dL"
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-04-06T01:00:00-04:00",
			     "Data": "ewFAAAEGEAIMAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 1,
			       "BasalRate": {
				 "Start": "01:00",
				 "Rate": 0.3
			       }
			     }
			   },
			   {
			     "Type": "DailyTotal523",
			     "Time": "2016-04-05T00:00:00-04:00",
			     "Data": "bkUQBQDnAAACAAACywFRLwF6NQB7ANYAFgCOAAAEAQEABAAAAAAAAAAAmzMAAAAAAAAAAA=="
			   },
			   {
			     "Type": "DailyTotal",
			     "Time": "2016-04-05T00:00:00-04:00",
			     "Data": "BwAAAstFEAAAAA==",
			     "Insulin": 17.875
			   },
			   {
			     "Type": "BasalProfileStart",
			     "Time": "2016-04-06T00:00:00-04:00",
			     "Data": "ewBAAAAGEAAQAA==",
			     "BasalProfileStart": {
			       "ProfileIndex": 0,
			       "BasalRate": {
				 "Start": "00:00",
				 "Rate": 0.4
			       }
			     }
			   },
			   {
			     "Type": "Bolus",
			     "Time": "2016-04-05T22:02:50-04:00",
			     "Data": "AQAWABYAQgByAlZlEA==",
			     "Bolus": {
			       "Duration": "0s",
			       "Programmed": 0.55,
			       "Amount": 0.55,
			       "Unabsorbed": 1.65
			     }
			   },
			   {
			     "Type": "UnabsorbedInsulin",
			     "Data": "XA4uQ8AaV8A4YcAWl9A=",
			     "UnabsorbedInsulin": [
			       {
				 "Age": "1h7m0s",
				 "Bolus": 1.15
			       },
			       {
				 "Age": "1h27m0s",
				 "Bolus": 0.65
			       },
			       {
				 "Age": "1h37m0s",
				 "Bolus": 1.4
			       },
			       {
				 "Age": "2h31m0s",
				 "Bolus": 0.55
			       }
			     ]
			   },
			   {
			     "Type": "BolusWizard",
			     "Time": "2016-04-05T22:02:50-04:00",
			     "Data": "WzNyAhZlEABRALRVUFgAAAAAQgAWeA==",
			     "BolusWizard": {
			       "GlucoseUnits": "mg/dL",
			       "GlucoseInput": 307,
			       "TargetLow": 80,
			       "TargetHigh": 120,
			       "Sensitivity": 85,
			       "CarbUnits": "Grams",
			       "CarbInput": 0,
			       "CarbRatio": 18,
			       "Unabsorbed": 1.65,
			       "Correction": 2.2,
			       "Food": 0,
			       "Bolus": 0.55
			     }
			   },
			   {
			     "Type": "BGCapture",
			     "Time": "2016-04-05T22:02:33-04:00",
			     "Data": "CjNhAjYFkA==",
			     "Glucose": 307,
			     "GlucoseUnits": "mg/dL"
			   },
			   {
			     "Type": "Bolus",
			     "Time": "2016-04-05T20:59:04-04:00",
			     "Data": "AQAuAC4ATABEO1RlEA==",
			     "Bolus": {
			       "Duration": "0s",
			       "Programmed": 1.15,
			       "Amount": 1.15,
			       "Unabsorbed": 1.9
			     }
			   },
			   {
			     "Type": "UnabsorbedInsulin",
			     "Data": "XA4aGMA4IsAWWNAC2tA=",
			     "UnabsorbedInsulin": [
			       {
				 "Age": "24m0s",
				 "Bolus": 0.65
			       },
			       {
				 "Age": "34m0s",
				 "Bolus": 1.4
			       },
			       {
				 "Age": "1h28m0s",
				 "Bolus": 0.55
			       },
			       {
				 "Age": "3h38m0s",
				 "Bolus": 0.05
			       }
			     ]
			   },
			   {
			     "Type": "BolusWizard",
			     "Time": "2016-04-05T20:59:04-04:00",
			     "Data": "WwBEOxRlEA5QAHhLUAAALgAAAAAueA==",
			     "BolusWizard": {
			       "GlucoseUnits": "mg/dL",
			       "GlucoseInput": 0,
			       "TargetLow": 80,
			       "TargetHigh": 120,
			       "Sensitivity": 75,
			       "CarbUnits": "Grams",
			       "CarbInput": 14,
			       "CarbRatio": 12,
			       "Unabsorbed": 0,
			       "Correction": 0,
			       "Food": 1.15,
			       "Bolus": 1.15
			     }
			   },
			   {
			     "Type": "Bolus",
			     "Time": "2016-04-05T20:34:02-04:00",
			     "Data": "AQBSAFIAAABCIlRlEA==",
			     "Bolus": {
			       "Duration": "0s",
			       "Programmed": 2.05,
			       "Amount": 2.05,
			       "Unabsorbed": 0
			     }
			   }
			 ]`,
			true,
		},
		{
			"06 15 05 0D 00 40 60 01 05 06 2F 18 73 00 40 20 C1 05 06 2F 0C 5B 00 40 20 C1 05 06 2F 0C 6C 00 40 20 C1 05 06 2F 0C 7D 00 40 20 C1 05 06 2F 0C 8E 00 40 20 C1 05 06 2F 0C 9F 00 40 20 C1 05 06 2F 0C C3 00 40 20 C1 05 06 2F 0C D4 00 40 20 C1 05 06 11 0C E3 00 40 20 C1 05 06 15 04 E9 00 40 40 A1 05 0C 15 16 40 00 01 05 64 00 19 43 00 01 05 17 00 1D 43 00 01 05 18 00 00 43 00 01 05 21 00 03 43 00 01 05 03 00 00 00 12 08 43 20 01 05 2C 50 18 52 A0 01 05 24 E5 30 53 00 01 05 1A 00 0D 57 02 01 05 06 03 05 0D 0D 57 62 01 05 0C 03 0C 40 00 01 05 64 00 0D 40 00 01 05 17 00 22 40 00 01 05 18 00 80 4A 02 17 10 07 00 00 00 36 01 85 6D 01 85 05 0C 00 E8 00 00 00 00 00 36 00 36 64 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 0C 00 E8 00 00 00 07 00 00 02 50 97 90 6D 97 90 05 0C 00 E8 00 00 00 00 02 50 02 50 64 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 0C 00 E8 00 00 00 60 01 9D 60 00 18 10 60 01 A0 60 00 18 10 60 01 A2 61 00 18 10 31 02 B6 57 01 18 10 5B 00 B6 57 01 18 10 01 50 08 21 78 00 01 00 00 00 00 01 78 01 01 01 00 B6 57 21 18 10 01 0A 0A 00 BA 5C 21 18 10 5E 01 AB 5D 01 18 10 31 02 86 5E 01 18 10 35 00 80 40 02 18 10 01 0A 0A 0F 86 5E 61 18 10 33 75 A5 65 0C 18 10 08 16 01 A5 65 0C 18 10 33 00 B4 65 0C 18 10 08 16 00 B4 65 0C 18 10 62 00 B8 65 0C 18 10 33 34 85 66 0C 18 10 00 16 01 85 66 0C 18 10 5B 00 B3 6E 0E 18 10 02 50 08 21 78 00 02 00 00 00 00 02 78 5C 14 02 67 14 02 85 14 02 99 14 02 AD 14 02 C1 14 02 DF 14 32 00 90 7A 0E 18 10 61 01 90 7A 0E 18 10 6A 00 B0 7A 0E 18 10 61 00 B0 7A 0E 18 10 32 4F 8C 40 0F 18 10 61 01 8C 40 0F 18 10 35 00 80 42 2F 18 10 1E 00 AF 45 0F 18 10 01 03 01 04 B3 6E 6E 18 10 1F 00 BA 45 0F 18 10 31 02 93 6D 14 18 10 5B 00 93 6D 14 18 10 06 50 08 21 78 00 07 00 00 00 00 07 78 5C 05 02 66 14 01 07 07 00 93 6D 54 18 10 01 0D 0D 00 8D 49 55 18 10 35 00 80 4F 15 18 10 01 1B 1B 00 98 43 56 18 10 33 00 82 7B 16 58 10 00 16 01 82 7B 16 58 10 6A 4F B8 61 17 18 10 61 00 B8 61 17 18 10 33 4A 85 63 17 18 10 05 16 01 85 63 17 18 10 33 00 B3 6A 17 18 10 00 16 00 B3 6A 17 18 10 33 90 98 6B 17 18 10 01 16 01 98 6B 17 18 10 33 00 A6 6D 17 18 10 00 16 00 A6 6D 17 18 10 33 C8 B7 6D 17 18 10 00 16 01 B7 6D 17 18 10 33 00 B1 73 17 18 10 00 16 00 B1 73 17 18 10 33 90 81 74 17 18 10 01 16 01 81 74 17 18 10 33 00 89 74 17 18 10 00 16 00 89 74 17 18 10 33 20 9A 74 17 18 10 03 16 01 9A 74 17 18 10 33 00 A6 74 17 18 10 00 16 00 A6 74 17 18 10 33 50 AE 74 17 18 10 05 16 01 AE 74 17 18 10",
			`[
			   {
			     "Type": "TempBasalDuration",
			     "Time": "2016-09-24T23:52:46-04:00",
			     "Duration": "30m0s",
			     "Data": "FgGudBcYEA=="
			   },
			   {
			     "Type": "TempBasalRate",
			     "Time": "2016-09-24T23:52:46-04:00",
			     "Data": "M1CudBcYEAU=",
			     "Insulin": 34,
			     "TempBasalType": "Absolute"
			   },
			   {
			     "Type": "TempBasalDuration",
			     "Time": "2016-09-24T23:52:38-04:00",
			     "Duration": "0s",
			     "Data": "FgCmdBcYEA=="
			   },
			   {
			     "Type": "TempBasalRate",
			     "Time": "2016-09-24T23:52:38-04:00",
			     "Data": "MwCmdBcYEAA=",
			     "Insulin": 0,
			     "TempBasalType": "Absolute"
			   },
			   {
			     "Type": "TempBasalDuration",
			     "Time": "2016-09-24T23:52:26-04:00",
			     "Duration": "30m0s",
			     "Data": "FgGadBcYEA=="
			   },
			   {
			     "Type": "TempBasalRate",
			     "Time": "2016-09-24T23:52:26-04:00",
			     "Data": "MyCadBcYEAM=",
			     "Insulin": 20,
			     "TempBasalType": "Absolute"
			   },
			   {
			     "Type": "TempBasalDuration",
			     "Time": "2016-09-24T23:52:09-04:00",
			     "Duration": "0s",
			     "Data": "FgCJdBcYEA=="
			   },
			   {
			     "Type": "TempBasalRate",
			     "Time": "2016-09-24T23:52:09-04:00",
			     "Data": "MwCJdBcYEAA=",
			     "Insulin": 0,
			     "TempBasalType": "Absolute"
			   },
			   {
			     "Type": "TempBasalDuration",
			     "Time": "2016-09-24T23:52:01-04:00",
			     "Duration": "30m0s",
			     "Data": "FgGBdBcYEA=="
			   },
			   {
			     "Type": "TempBasalRate",
			     "Time": "2016-09-24T23:52:01-04:00",
			     "Data": "M5CBdBcYEAE=",
			     "Insulin": 10,
			     "TempBasalType": "Absolute"
			   },
			   {
			     "Type": "TempBasalDuration",
			     "Time": "2016-09-24T23:51:49-04:00",
			     "Duration": "0s",
			     "Data": "FgCxcxcYEA=="
			   },
			   {
			     "Type": "TempBasalRate",
			     "Time": "2016-09-24T23:51:49-04:00",
			     "Data": "MwCxcxcYEAA=",
			     "Insulin": 0,
			     "TempBasalType": "Absolute"
			   },
			   {
			     "Type": "TempBasalDuration",
			     "Time": "2016-09-24T23:45:55-04:00",
			     "Duration": "30m0s",
			     "Data": "FgG3bRcYEA=="
			   },
			   {
			     "Type": "TempBasalRate",
			     "Time": "2016-09-24T23:45:55-04:00",
			     "Data": "M8i3bRcYEAA=",
			     "Insulin": 5,
			     "TempBasalType": "Absolute"
			   },
			   {
			     "Type": "TempBasalDuration",
			     "Time": "2016-09-24T23:45:38-04:00",
			     "Duration": "0s",
			     "Data": "FgCmbRcYEA=="
			   },
			   {
			     "Type": "TempBasalRate",
			     "Time": "2016-09-24T23:45:38-04:00",
			     "Data": "MwCmbRcYEAA=",
			     "Insulin": 0,
			     "TempBasalType": "Absolute"
			   },
			   {
			     "Type": "TempBasalDuration",
			     "Time": "2016-09-24T23:43:24-04:00",
			     "Duration": "30m0s",
			     "Data": "FgGYaxcYEA=="
			   },
			   {
			     "Type": "TempBasalRate",
			     "Time": "2016-09-24T23:43:24-04:00",
			     "Data": "M5CYaxcYEAE=",
			     "Insulin": 10,
			     "TempBasalType": "Absolute"
			   },
			   {
			     "Type": "TempBasalDuration",
			     "Time": "2016-09-24T23:42:51-04:00",
			     "Duration": "0s",
			     "Data": "FgCzahcYEA=="
			   },
			   {
			     "Type": "TempBasalRate",
			     "Time": "2016-09-24T23:42:51-04:00",
			     "Data": "MwCzahcYEAA=",
			     "Insulin": 0,
			     "TempBasalType": "Absolute"
			   },
			   {
			     "Type": "TempBasalDuration",
			     "Time": "2016-09-24T23:35:05-04:00",
			     "Duration": "30m0s",
			     "Data": "FgGFYxcYEA=="
			   },
			   {
			     "Type": "TempBasalRate",
			     "Time": "2016-09-24T23:35:05-04:00",
			     "Data": "M0qFYxcYEAU=",
			     "Insulin": 33.85,
			     "TempBasalType": "Absolute"
			   },
			   {
			     "Type": "EnableAlarmClock",
			     "Time": "2016-09-24T23:33:56-04:00",
			     "Data": "YQC4YRcYEA==",
			     "Enabled": false
			   },
			   {
			     "Type": "DeleteAlarmClockTime",
			     "Time": "2016-09-24T23:33:56-04:00",
			     "Data": "ak+4YRcYEA=="
			   },
			   {
			     "Type": "TempBasalDuration",
			     "Time": "2016-09-24T22:59:02-04:00",
			     "Duration": "30m0s",
			     "Data": "FgGCexZYEA=="
			   },
			   {
			     "Type": "TempBasalRate",
			     "Time": "2016-09-24T22:59:02-04:00",
			     "Data": "MwCCexZYEAA=",
			     "Insulin": 0,
			     "TempBasalType": "Absolute"
			   },
			   {
			     "Type": "Bolus",
			     "Time": "2016-09-24T22:03:24-04:00",
			     "Data": "ARsbAJhDVhgQ",
			     "Bolus": {
			       "Duration": "0s",
			       "Programmed": 2.7,
			       "Amount": 2.7,
			       "Unabsorbed": 0
			     }
			   },
			   {
			     "Type": "AlarmClock",
			     "Time": "2016-09-24T21:15:00-04:00",
			     "Data": "NQCATxUYEA=="
			   },
			   {
			     "Type": "Bolus",
			     "Time": "2016-09-24T21:09:13-04:00",
			     "Data": "AQ0NAI1JVRgQ",
			     "Bolus": {
			       "Duration": "0s",
			       "Programmed": 1.3,
			       "Amount": 1.3,
			       "Unabsorbed": 0
			     }
			   },
			   {
			     "Type": "Bolus",
			     "Time": "2016-09-24T20:45:19-04:00",
			     "Data": "AQcHAJNtVBgQ",
			     "Bolus": {
			       "Duration": "0s",
			       "Programmed": 0.7,
			       "Amount": 0.7,
			       "Unabsorbed": 0
			     }
			   },
			   {
			     "Type": "UnabsorbedInsulin",
			     "Data": "XAUCZhQ=",
			     "UnabsorbedInsulin": [
			       {
				 "Age": "1h42m0s",
				 "Bolus": 0.05
			       }
			     ]
			   },
			   {
			     "Type": "BolusWizard",
			     "Time": "2016-09-24T20:45:19-04:00",
			     "Data": "WwCTbRQYEAZQCCF4AAcAAAAAB3g=",
			     "BolusWizard": {
			       "GlucoseUnits": "mg/dL",
			       "GlucoseInput": 0,
			       "TargetLow": 120,
			       "TargetHigh": 120,
			       "Sensitivity": 33,
			       "CarbUnits": "Grams",
			       "CarbInput": 6,
			       "CarbRatio": 8,
			       "Unabsorbed": 0,
			       "Correction": 0,
			       "Food": 0.7,
			       "Bolus": 0.7
			     }
			   },
			   {
			     "Type": "ChangeBGReminder",
			     "Time": "2016-09-24T20:45:19-04:00",
			     "Data": "MQKTbRQYEA=="
			   },
			   {
			     "Type": "ResumePump",
			     "Time": "2016-09-24T15:05:58-04:00",
			     "Data": "HwC6RQ8YEA=="
			   },
			   {
			     "Type": "Bolus",
			     "Time": "2016-09-24T14:46:51-04:00",
			     "Data": "AQMBBLNubhgQ",
			     "Bolus": {
			       "Duration": "2h0m0s",
			       "Programmed": 0.3,
			       "Amount": 0.1,
			       "Unabsorbed": 0
			     }
			   },
			   {
			     "Type": "SuspendPump",
			     "Time": "2016-09-24T15:05:47-04:00",
			     "Data": "HgCvRQ8YEA=="
			   },
			   {
			     "Type": "AlarmClock",
			     "Time": "2016-09-24T15:02:00-04:00",
			     "Data": "NQCAQi8YEA=="
			   },
			   {
			     "Type": "EnableAlarmClock",
			     "Time": "2016-09-24T15:00:12-04:00",
			     "Data": "YQGMQA8YEA==",
			     "Enabled": true
			   },
			   {
			     "Type": "SetAlarmClockTime",
			     "Time": "2016-09-24T15:00:12-04:00",
			     "Data": "Mk+MQA8YEA=="
			   },
			   {
			     "Type": "EnableAlarmClock",
			     "Time": "2016-09-24T14:58:48-04:00",
			     "Data": "YQCweg4YEA==",
			     "Enabled": false
			   },
			   {
			     "Type": "DeleteAlarmClockTime",
			     "Time": "2016-09-24T14:58:48-04:00",
			     "Data": "agCweg4YEA=="
			   },
			   {
			     "Type": "EnableAlarmClock",
			     "Time": "2016-09-24T14:58:16-04:00",
			     "Data": "YQGQeg4YEA==",
			     "Enabled": true
			   },
			   {
			     "Type": "SetAlarmClockTime",
			     "Time": "2016-09-24T14:58:16-04:00",
			     "Data": "MgCQeg4YEA=="
			   },
			   {
			     "Type": "UnabsorbedInsulin",
			     "Data": "XBQCZxQChRQCmRQCrRQCwRQC3xQ=",
			     "UnabsorbedInsulin": [
			       {
				 "Age": "1h43m0s",
				 "Bolus": 0.05
			       },
			       {
				 "Age": "2h13m0s",
				 "Bolus": 0.05
			       },
			       {
				 "Age": "2h33m0s",
				 "Bolus": 0.05
			       },
			       {
				 "Age": "2h53m0s",
				 "Bolus": 0.05
			       },
			       {
				 "Age": "3h13m0s",
				 "Bolus": 0.05
			       },
			       {
				 "Age": "3h43m0s",
				 "Bolus": 0.05
			       }
			     ]
			   },
			   {
			     "Type": "BolusWizard",
			     "Time": "2016-09-24T14:46:51-04:00",
			     "Data": "WwCzbg4YEAJQCCF4AAIAAAAAAng=",
			     "BolusWizard": {
			       "GlucoseUnits": "mg/dL",
			       "GlucoseInput": 0,
			       "TargetLow": 120,
			       "TargetHigh": 120,
			       "Sensitivity": 33,
			       "CarbUnits": "Grams",
			       "CarbInput": 2,
			       "CarbRatio": 8,
			       "Unabsorbed": 0,
			       "Correction": 0,
			       "Food": 0.2,
			       "Bolus": 0.2
			     }
			   },
			   {
			     "Type": "TempBasalDuration",
			     "Time": "2016-09-24T12:38:05-04:00",
			     "Duration": "30m0s",
			     "Data": "FgGFZgwYEA=="
			   },
			   {
			     "Type": "TempBasalRate",
			     "Time": "2016-09-24T12:38:05-04:00",
			     "Data": "MzSFZgwYEAA=",
			     "Insulin": 1.3,
			     "TempBasalType": "Absolute"
			   },
			   {
			     "Type": "ChangeTempBasalType",
			     "Time": "2016-09-24T12:37:56-04:00",
			     "Data": "YgC4ZQwYEA==",
			     "TempBasalType": "Absolute"
			   },
			   {
			     "Type": "TempBasalDuration",
			     "Time": "2016-09-24T12:37:52-04:00",
			     "Duration": "0s",
			     "Data": "FgC0ZQwYEA=="
			   },
			   {
			     "Type": "TempBasalRate",
			     "Time": "2016-09-24T12:37:52-04:00",
			     "Data": "MwC0ZQwYEAg=",
			     "TempBasalType": "Percent",
			     "Value": 0
			   },
			   {
			     "Type": "TempBasalDuration",
			     "Time": "2016-09-24T12:37:37-04:00",
			     "Duration": "30m0s",
			     "Data": "FgGlZQwYEA=="
			   },
			   {
			     "Type": "TempBasalRate",
			     "Time": "2016-09-24T12:37:37-04:00",
			     "Data": "M3WlZQwYEAg=",
			     "TempBasalType": "Percent",
			     "Value": 117
			   },
			   {
			     "Type": "Bolus",
			     "Time": "2016-09-24T01:30:06-04:00",
			     "Data": "AQoKD4ZeYRgQ",
			     "Bolus": {
			       "Duration": "7h30m0s",
			       "Programmed": 1,
			       "Amount": 1,
			       "Unabsorbed": 0
			     }
			   },
			   {
			     "Type": "AlarmClock",
			     "Time": "2016-09-24T02:00:00-04:00",
			     "Data": "NQCAQAIYEA=="
			   },
			   {
			     "Type": "ChangeBGReminder",
			     "Time": "2016-09-24T01:30:06-04:00",
			     "Data": "MQKGXgEYEA=="
			   },
			   {
			     "Type": "EnableVariableBolus",
			     "Time": "2016-09-24T01:29:43-04:00",
			     "Data": "XgGrXQEYEA==",
			     "Enabled": true
			   },
			   {
			     "Type": "Bolus",
			     "Time": "2016-09-24T01:28:58-04:00",
			     "Data": "AQoKALpcIRgQ",
			     "Bolus": {
			       "Duration": "0s",
			       "Programmed": 1,
			       "Amount": 1,
			       "Unabsorbed": 0
			     }
			   },
			   {
			     "Type": "Bolus",
			     "Time": "2016-09-24T01:23:54-04:00",
			     "Data": "AQEBALZXIRgQ",
			     "Bolus": {
			       "Duration": "0s",
			       "Programmed": 0.1,
			       "Amount": 0.1,
			       "Unabsorbed": 0
			     }
			   },
			   {
			     "Type": "BolusWizard",
			     "Time": "2016-09-24T01:23:54-04:00",
			     "Data": "WwC2VwEYEAFQCCF4AAEAAAAAAXg=",
			     "BolusWizard": {
			       "GlucoseUnits": "mg/dL",
			       "GlucoseInput": 0,
			       "TargetLow": 120,
			       "TargetHigh": 120,
			       "Sensitivity": 33,
			       "CarbUnits": "Grams",
			       "CarbInput": 1,
			       "CarbRatio": 8,
			       "Unabsorbed": 0,
			       "Correction": 0,
			       "Food": 0.1,
			       "Bolus": 0.1
			     }
			   },
			   {
			     "Type": "ChangeBGReminder",
			     "Time": "2016-09-24T01:23:54-04:00",
			     "Data": "MQK2VwEYEA=="
			   },
			   {
			     "Type": "EnableBGReminder",
			     "Time": "2016-09-24T00:33:34-04:00",
			     "Data": "YAGiYQAYEA==",
			     "Enabled": true
			   },
			   {
			     "Type": "EnableBGReminder",
			     "Time": "2016-09-24T00:32:32-04:00",
			     "Data": "YAGgYAAYEA==",
			     "Enabled": true
			   },
			   {
			     "Type": "EnableBGReminder",
			     "Time": "2016-09-24T00:32:29-04:00",
			     "Data": "YAGdYAAYEA==",
			     "Enabled": true
			   },
			   {
			     "Type": "DailyTotal522",
			     "Time": "2016-09-23T00:00:00-04:00",
			     "Data": "bZeQBQwA6AAAAAACUAJQZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMAOgAAAA="
			   },
			   {
			     "Type": "DailyTotal",
			     "Time": "2016-09-23T00:00:00-04:00",
			     "Data": "BwAAAlCXkA==",
			     "Insulin": 14.8
			   },
			   {
			     "Type": "DailyTotal522",
			     "Time": "2005-01-01T00:00:00-05:00",
			     "Data": "bQGFBQwA6AAAAAAANgA2ZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMAOgAAAA="
			   },
			   {
			     "Type": "DailyTotal",
			     "Time": "2005-01-01T00:00:00-05:00",
			     "Data": "BwAAADYBhQ==",
			     "Insulin": 1.35
			   },
			   {
			     "Type": "NewTime",
			     "Time": "2016-09-23T02:10:00-04:00",
			     "Data": "GACASgIXEA=="
			   },
			   {
			     "Type": "ChangeTime",
			     "Time": "2005-01-01T00:00:34-05:00",
			     "Data": "FwAiQAABBQ=="
			   },
			   {
			     "Type": "ChangeTimeFormat",
			     "Time": "2005-01-01T00:00:13-05:00",
			     "Data": "ZAANQAABBQ==",
			     "Value": 0
			   },
			   {
			     "Type": "ClearAlarm",
			     "Time": "2005-01-01T00:00:12-05:00",
			     "Data": "DAMMQAABBQ==",
			     "Value": 3
			   },
			   {
			     "Type": "Alarm",
			     "Time": "2005-01-01T02:23:13-05:00",
			     "Data": "BgMFDQ1XYgEF",
			     "Value": 3
			   },
			   {
			     "Type": "BatteryChange",
			     "Time": "2005-01-01T02:23:13-05:00",
			     "Data": "GgANVwIBBQ=="
			   },
			   {
			     "Type": "MaxBolus",
			     "Time": "2005-01-01T00:19:48-05:00",
			     "Data": "JOUwUwABBQ==",
			     "Insulin": 5.725
			   },
			   {
			     "Type": "MaxBasal",
			     "Time": "2005-01-01T00:18:24-05:00",
			     "Data": "LFAYUqABBQ==",
			     "Insulin": 2
			   },
			   {
			     "Type": "Prime",
			     "Time": "2005-01-01T00:03:08-05:00",
			     "Data": "AwAAABIIQyABBQ==",
			     "Prime": {
			       "Fixed": 0,
			       "Manual": 1.8
			     }
			   },
			   {
			     "Type": "Rewind",
			     "Time": "2005-01-01T00:03:03-05:00",
			     "Data": "IQADQwABBQ=="
			   },
			   {
			     "Type": "NewTime",
			     "Time": "2005-01-01T00:03:00-05:00",
			     "Data": "GAAAQwABBQ=="
			   },
			   {
			     "Type": "ChangeTime",
			     "Time": "2005-01-01T00:03:29-05:00",
			     "Data": "FwAdQwABBQ=="
			   },
			   {
			     "Type": "ChangeTimeFormat",
			     "Time": "2005-01-01T00:03:25-05:00",
			     "Data": "ZAAZQwABBQ==",
			     "Value": 0
			   },
			   {
			     "Type": "ClearAlarm",
			     "Time": "2005-01-01T00:00:22-05:00",
			     "Data": "DBUWQAABBQ==",
			     "Value": 21
			   },
			   {
			     "Type": "Alarm",
			     "Time": "2005-01-01T00:00:00-05:00",
			     "Data": "BhUE6QBAQKEF",
			     "Value": 21
			   },
			   {
			     "Type": "Alarm",
			     "Time": "2005-01-01T00:00:00-05:00",
			     "Data": "BhEM4wBAIMEF",
			     "Value": 17
			   },
			   {
			     "Type": "Alarm",
			     "Time": "2005-01-01T00:00:00-05:00",
			     "Data": "Bi8M1ABAIMEF",
			     "Value": 47
			   },
			   {
			     "Type": "Alarm",
			     "Time": "2005-01-01T00:00:00-05:00",
			     "Data": "Bi8MwwBAIMEF",
			     "Value": 47
			   },
			   {
			     "Type": "Alarm",
			     "Time": "2005-01-01T00:00:00-05:00",
			     "Data": "Bi8MnwBAIMEF",
			     "Value": 47
			   },
			   {
			     "Type": "Alarm",
			     "Time": "2005-01-01T00:00:00-05:00",
			     "Data": "Bi8MjgBAIMEF",
			     "Value": 47
			   },
			   {
			     "Type": "Alarm",
			     "Time": "2005-01-01T00:00:00-05:00",
			     "Data": "Bi8MfQBAIMEF",
			     "Value": 47
			   },
			   {
			     "Type": "Alarm",
			     "Time": "2005-01-01T00:00:00-05:00",
			     "Data": "Bi8MbABAIMEF",
			     "Value": 47
			   },
			   {
			     "Type": "Alarm",
			     "Time": "2005-01-01T00:00:00-05:00",
			     "Data": "Bi8MWwBAIMEF",
			     "Value": 47
			   },
			   {
			     "Type": "Alarm",
			     "Time": "2005-01-01T00:00:00-05:00",
			     "Data": "Bi8YcwBAIMEF",
			     "Value": 47
			   },
			   {
			     "Type": "Alarm",
			     "Time": "2005-01-01T00:00:00-05:00",
			     "Data": "BhUFDQBAYAEF",
			     "Value": 21
			   }
			 ]`,
			false,
		},
	}
	for _, c := range cases {
		data := readBytes(c.page)
		var records []HistoryRecord
		err := json.Unmarshal([]byte(c.results), &records)
		if err != nil {
			t.Errorf("json.Unmarshal(%s) returned %v", c.results, err)
			continue
		}
		decoded, err := DecodeHistoryRecords(data, c.newer)
		if err != nil {
			t.Errorf("DecodeHistoryRecords(%X, %v) returned %v", data, c.newer, err)
			continue
		}
		if !reflect.DeepEqual(decoded, records) {
			t.Errorf("DecodeHistoryRecords(%X, %v) == %v, want %v", data, c.newer, decoded, records)
		}
	}
}

func readBytes(hex string) []byte {
	f := strings.Fields(hex)
	data := make([]byte, len(f))
	for i, s := range f {
		b, err := strconv.ParseUint(s, 16, 8)
		if err != nil {
			log.Fatal(err)
		}
		data[i] = byte(b)
	}
	return data
}
