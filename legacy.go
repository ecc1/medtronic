package medtronic

type (
	// Entry represents data for the Nightscout entries API.
	LegacyEntry struct {
		Name string  `json:"name"`
		Type string  `json:"data_type"`
		SGV  int `json:"sgv"`
		Date string   `json:"date"` // Unix time in milliseconds
		Size int     `json:"packet_size"`

		// _tell
		// op
	}

	// Entries represents a sequence of Entry values.
	LegacyEntries []LegacyEntry
)

func FormatToOAPS(data CGMHistory) LegacyEntries {
	var entries LegacyEntries
	for _, r := range data {
		if r.Type != CGMGlucose {
			continue
		}
		t := r.Time
		e := LegacyEntry{
			Name: "GlucoseSensorData",
			Type: "prevTimestamp",
			SGV:  r.Glucose,
			Date: t.Format("2006-01-02T15:04:05"),
			Size: 0,
		}
		entries = append(entries, e)
	}
	return entries
}
