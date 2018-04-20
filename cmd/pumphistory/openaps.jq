# This is a jq filter that converts the output of
# https://github.com/ecc1/medtronic/cmd/pumphistory
# to the schema that openaps expects.
# Usage:
#   pumphistory | jq -f openaps.jq > history.json

# Convert a Go duration to minutes.
def duration_to_minutes:
  if test("^(\\d+h)?(\\d+m)?0s$") then
    capture("^((?<h>\\d+)h)?((?<m>\\d+)m)?0s$") |
    60*(.h // 0 | tonumber) + (.m // 0 | tonumber)
  else
    ("unexpected duration: " + .) | error
  end;

# Output a single JSON array.
[
  # Treat null as an empty array.
  (. // []) |
  # Perform the following on each element of the input array.
  .[] |
  # Filter out fake meter entries.
  select((.Type == "BGReceived" and .Info.MeterID == "000000") | not) |
  # Start with the timestamp field, common to all record types,
  # and the type, which is the same as the decocare type in many cases.
  {
    timestamp: .Time,
    _type: .Type
  } +
  # Add type-specific fields.
  if .Type == "TempBasalDuration" then
    {
      "duration (min)": .Info | duration_to_minutes
    }
  elif .Type == "TempBasalRate" then
    {
      _type: "TempBasal",
      temp: .Info.Type | ascii_downcase,
      rate: .Info.Value
    }
  elif .Type == "Bolus" then
    {
      amount: .Info.Amount,
      programmed: .Info.Programmed,
      unabsorbed: .Info.Unabsorbed,
      duration: .Info.Duration | duration_to_minutes
    }
  elif .Type == "BolusWizard" then
    {
      bg: (.Info.Glucose // 0),
      carb_input: (.Info.CarbInput // 0),
      bg_target_low: .Info.TargetLow,
      bg_target_high: .Info.TargetHigh,
      sensitivity: .Info.Sensitivity,
      carb_ratio: .Info.CarbRatio,
      correction_estimate: .Info.Correction,
      food_estimate: .Info.Food,
      unabsorbed_insulin_total: .Info.Unabsorbed,
      bolus_estimate: .Info.Bolus
    }
  elif .Type == "Prime" then
    {
      amount: .Info.Manual,
      fixed: .Info.Fixed,
      type: (if .Info.Fixed == 0 then "manual" else "fixed" end),
    }
  elif .Type == "BGReceived" then
    {
      glucose: .Info.Glucose,
      meterID: .Info.MeterID
    }
  elif .Type == "SuspendPump" then
    {
      _type: "PumpSuspend"
    }
  elif .Type == "ResumePump" then
    {
      _type: "PumpResume"
    }
  elif .Type == "BatteryChange" then
    {
      _type: "Battery"
    }
  elif .Type == "Rewind" then
    {
    }
  # Add additional cases here as needed.
  else
    # Warn about record types being skipped.
    ("skipping " + .Type + " record") | debug | empty
  end
]
