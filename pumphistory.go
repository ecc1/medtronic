package medtronic

import (
	"log"
	"time"
)

// History returns the history records since the specified time.
// Note that the results may include records with a zero timestamp or
// an earlier timestamp than the cutoff (in the case of DailyTotal records).
func (pump *Pump) History(since time.Time) History {
	count := pump.HistoryPageCount()
	if pump.Error() != nil {
		return nil
	}
	family := pump.Family()
	var results History
	for page := 0; page < count && pump.Error() == nil; page++ {
		data := pump.HistoryPage(page)
		records, err := DecodeHistory(data, family)
		if err != nil {
			pump.SetError(err)
		}
		i := findSince(records, since)
		results = append(results, records[:i]...)
		if i < len(records) {
			break
		}
	}
	return results
}

// findSince finds the first record that did not occur after the cutoff and returns its index,
// or len(records) if all the records occur more recently.
func findSince(records History, cutoff time.Time) int {
	for i, r := range records {
		switch r.Type() {
		case DailyTotal, DailyTotal515, DailyTotal522, DailyTotal523:
			// Don't use DailyTotal timestamps to decide when to stop,
			// because they appear out of order (at the end of the day).
		default:
			t := r.Time
			if !t.IsZero() && !t.After(cutoff) {
				log.Printf("stopping pump history scan at %s", t.Format(UserTimeLayout))
				return i
			}
		}
	}
	return len(records)
}
