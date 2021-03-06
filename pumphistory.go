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
	lastPage := pump.LastHistoryPage()
	if pump.Error() != nil {
		return nil
	}
	family := pump.Family()
	var results History
	for page := 0; page <= lastPage && pump.Error() == nil; page++ {
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

func (pump *Pump) findRecords(check func(HistoryRecord) bool, msg func(HistoryRecord)) (History, bool) {
	results := pump.findHistory(check)
	n := len(results)
	if n == 0 {
		return nil, false
	}
	r := results[n-1]
	if check(r) {
		msg(r)
		return results[:n-1], true
	}
	return results, false
}

// History returns the history records since the specified time.
// Note that the results may include records with a zero timestamp or
// an earlier timestamp than the cutoff (in the case of DailyTotal records).
func (pump *Pump) History(since time.Time) History {
	check := func(r HistoryRecord) bool {
		return checkBefore(r, since)
	}
	msg := func(r HistoryRecord) {
		log.Printf("stopping pump history scan at %s", r.Time.Format(UserTimeLayout))
	}
	results, _ := pump.findRecords(check, msg)
	return results
}

// checkBefore returns true if r occurred no later than the cutoff.
// The time must also be after 2015, to avoid stopping too soon when
// history records are created before the pump clock has been set correctly.
func checkBefore(r HistoryRecord, cutoff time.Time) bool {
	switch r.Type() {
	// Don't use DailyTotal timestamps to decide when to stop,
	// because they appear out of order (at the end of the day).
	case DailyTotal, DailyTotal515, DailyTotal522, DailyTotal523:
		return false
	}
	t := r.Time
	return !t.After(cutoff) && t.Year() > 2015
}

// HistoryFrom returns the history records since the specified record ID
// along with a bool indicating whether it was found. If the record ID
// was not found, the result will contain the entire pump history.
func (pump *Pump) HistoryFrom(id []byte) (History, bool) {
	check := func(r HistoryRecord) bool {
		return checkID(r, id)
	}
	msg := func(r HistoryRecord) {
		log.Printf("stopping pump history scan at record %s", base64.StdEncoding.EncodeToString(id))
	}
	return pump.findRecords(check, msg)
}

// checkID returns true if r has a given record ID.
func checkID(r HistoryRecord, id []byte) bool {
	return bytes.Equal(r.Data, id)
}
