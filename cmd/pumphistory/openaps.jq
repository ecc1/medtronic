# This is a jq filter that converts the output of
# https://github.com/ecc1/medtronic/cmd/pumphistory
# to the schema that openaps expects.
# Usage:
#   pumphistory | jq -f openaps.jq > history.json

# Convert a Go duration of the form Nm0s to minutes.
# This will fail (in the tonumber conversion)
# if the input has an hours or nonzero seconds component.
def duration_to_minutes:
  . | rtrimstr("m0s") | tonumber;

# Output a single JSON array.
[
  # Perform the following on each element of the input array.
  .[] |
  # Start with the timestamp field, common to all record types.
  {
    timestamp: .Time,
  } +
  # Add type-specific fields.
  if .Type == "TempBasalDuration" then
    {
      _type: "TempBasalDuration",
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
      _type: "Bolus",
      amount: .Info.Amount
    }
  # Add additional cases here as needed.
  else
    # Warn about record types being skipped.
    ("skipping " + .Type + " record") | debug | empty
  end
]
