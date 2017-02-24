package medtronic

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestHistoryRecord(t *testing.T) {
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
			   "Data": "Cm5l2S4EEA=="
			 }`,
			false,
		},
		{
			`{
			   "Type": "BgCapture",
			   "Time": "2016-05-19T03:49:06-04:00",
			   "Glucose": 500,
			   "Data": "CvRGcSMTkA=="
			 }`,
			false,
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
				   "CarbRatio": 15,
				   "Units": "Grams"
				 },
				 {
				   "Start": "04:00",
				   "CarbRatio": 20,
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
				   "CarbRatio": 10,
				   "Units": "Exchanges"
				 },
				 {
				   "Start": "01:00",
				   "CarbRatio": 15,
				   "Units": "Exchanges"
				 },
				 {
				   "Start": "02:00",
				   "CarbRatio": 12,
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
			   "Time": "2016-05-09T12:07:45-04:00",
			   "BolusWizardSetup": {
			     "Before": {
			       "InsulinAction": "4h0m0s",
			       "Ratios": [
				 {
				   "Start": "00:00",
				   "CarbRatio": 6,
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
				   "CarbRatio": 6,
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
			     "GlucoseInput": 110,
			     "TargetLow": 99,
			     "TargetHigh": 101,
			     "Sensitivity": 45,
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
			     "GlucoseInput": 131,
			     "TargetLow": 100,
			     "TargetHigh": 100,
			     "Sensitivity": 35,
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
			t.Errorf("decodeHistoryRecord(%X, %v) returned %v, want %v", r1.Data, c.newer, err, r1)
			continue
		}
		if !r1.equal(r2) {
			t.Errorf("decodeHistoryRecord(%X, %v) == %v, want %v", r1.Data, c.newer, r2, r1)
		}
	}
}

func (r1 HistoryRecord) equal(r2 HistoryRecord) bool {
	return reflect.DeepEqual(r1, r2)
}

func (r HistoryRecord) String() string {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}
