package medtronic

import (
	"log"
	"time"
)

// HistoryRecords returns the history records since the specified time.
// Note that the results may include records with a zero timestamp or
// an earlier timestamp than the cutoff (in the case of DailyTotal records).
func (pump *Pump) HistoryRecords(since time.Time) []HistoryRecord {
	newer := pump.Family() >= 23
	lastPage := pump.LastHistoryPage()
	if pump.Error() != nil {
		return nil
	}
	var results []HistoryRecord
	for page := 0; page <= lastPage && pump.Error() == nil; page++ {
		data := pump.HistoryPage(page)
		records, err := DecodeHistoryRecords(data, newer)
		if err != nil {
			pump.SetError(err)
		}
		i := findOlder(records, since)
		results = append(results, records[:i]...)
		if i < len(records) {
			break
		}
	}
	return results
}

// findOlder finds the first record that is older than the given time and returns its index,
// or len(records) if all the records occur more recently.
func findOlder(records []HistoryRecord, cutoff time.Time) int {
	for i, r := range records {
		// Don't use DailyTotal timestamps to decide when to stop,
		// because they appear out of order (at the end of the day).
		switch r.Type() {
		case DailyTotal:
		case DailyTotal522:
		case DailyTotal523:
		default:
			t := time.Time(r.Time)
			if !t.IsZero() && t.Before(cutoff) {
				log.Printf("stopping pump history scan at %s", t.Format(UserTimeLayout))
				return i
			}
		}
	}
	return len(records)
}
