# This is a jq filter that converts the output of
# https://github.com/ecc1/medtronic/cmd/pumphistory
# to the schema that openaps expects.
# Usage:
#   pumphistory | jq -f openaps.jq > history.json

# Convert a Go duration to minutes.
def duration_to_minutes:
  if test("^(\\d+h)?\\d+m0s$") then
    capture("^((?<h>\\d+)h)?(?<m>\\d+)m0s$") |
    60*(.h // 0 | tonumber) + (.m | tonumber)
  else
    ("unexpected duration: " + .) | error
  end;

# Output a single JSON array.
[
  # Perform the following on each element of the input array.
  .[] |
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
      amount: .Info.Amount
    }
  elif .Type == "Prime" then
    {
      amount: .Info.Manual,
      fixed: .Info.Fixed,
      type: "fixed"
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
