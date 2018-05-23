package medtronic

import (
	"log"
	"time"
)

// CGMHistory returns the CGM records since the specified time.
func (pump *Pump) CGMHistory(since time.Time) CGMHistory {
	i, j := pump.CGMPageRange()
	if pump.Error() != nil {
		return nil
	}
	var results CGMHistory
	for page := i; page < j && pump.Error() == nil; page++ {
		data := pump.GlucosePage(page)
		records, err := DecodeCGMHistory(data)
		if err != nil {
			pump.SetError(err)
		}
		i := findCGMSince(records, since)
		results = append(results, records[:i]...)
		if i < len(records) {
			break
		}
	}
	return results
}

// findCGMSince finds the first record that did not occur after the cutoff and returns its index,
// or len(records) if all the records occur more recently.
func findCGMSince(records CGMHistory, cutoff time.Time) int {
	for i, r := range records {
		t := r.Time
		if !t.IsZero() && !t.After(cutoff) {
			log.Printf("stopping CGM history scan at %s", t.Format(UserTimeLayout))
			return i
		}
	}
	return len(records)
}
