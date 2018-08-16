package medtronic

import (
	"bytes"
	"encoding/base64"
	"log"
	"time"
)

// findHistory retrieves history records from the pump
// until it encounters one that satisfies the given predicate,
// in which case that record will be the final element of the result.
// If the predicate is never satisfied, the entire pump history is returned.
// The two cases can be distinguished by checking whether
// the final element of the result satisfies the predicate.
// The records are retrieved and returned in reverse chronological order.
func (pump *Pump) findHistory(check func(HistoryRecord) bool) History {
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
		for i, r := range records {
			if check(r) {
				return append(results, records[:i+1]...)
			}
		}
		results = append(results, records...)
	}
	return results
}

// History returns the history records since the specified time.
// Note that the results may include records with a zero timestamp or
// an earlier timestamp than the cutoff (in the case of DailyTotal records).
func (pump *Pump) History(since time.Time) History {
	check := func(r HistoryRecord) bool {
		return checkSince(r, since)
	}
	results := pump.findHistory(check)
	n := len(results)
	if n == 0 {
		return nil
	}
	r := results[n-1]
	if checkSince(r, since) {
		log.Printf("stopping pump history scan at %s", r.Time.Format(UserTimeLayout))
		return results[:n-1]
	}
	return results
}

// checkSince returns true if r occurred no later than the cutoff.
func checkSince(r HistoryRecord, cutoff time.Time) bool {
	switch r.Type() {
	// Don't use DailyTotal timestamps to decide when to stop,
	// because they appear out of order (at the end of the day).
	case DailyTotal, DailyTotal515, DailyTotal522, DailyTotal523:
		return false
	}
	t := r.Time
	return !t.IsZero() && !t.After(cutoff)
}

// HistoryFrom returns the history records since the specified record ID
// along with a bool indicating whether it was found. If the record ID
// was not found, the result will contain the entire pump history.
func (pump *Pump) HistoryFrom(id []byte) (History, bool) {
	check := func(r HistoryRecord) bool {
		return checkID(r, id)
	}
	results := pump.findHistory(check)
	n := len(results)
	if n == 0 {
		return nil, false
	}
	r := results[n-1]
	if checkID(r, id) {
		log.Printf("stopping pump history scan at record %s", base64.StdEncoding.EncodeToString(id))
		return results[:n-1], true
	}
	return results, false
}

// checkID returns true if r has a given record ID.
func checkID(r HistoryRecord, id []byte) bool {
	return bytes.Equal(r.Data, id)
}
