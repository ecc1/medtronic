package medtronic

import (
	"encoding/base64"
	"testing"
	"time"
)

func TestPumpHistory(t *testing.T) {
	tc := testCase{"pump-records", 523, 0}
	file := testFileName(tc) + ".json"
	family := testPumpFamily(tc)
	records, err := decodeFromData(file, family)
	if err != nil {
		t.Error(err)
		return
	}
	cases := []struct {
		cutoff string
		index  int
	}{
		{"2020-02-25T11:00", 0},
		{"2020-02-25T09:00", 8},
		{"2020-02-24T23:00", 54},
		{"2020-02-24T01:00", 297},
		{"2000-01-01T00:00", 297},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			results := findSince(records, parseTime(c.cutoff))
			if len(results) != c.index {
				t.Errorf("findSince(%s) returned %d records, want %d", c.cutoff, len(results), c.index)
			}
		})
	}
}

func TestPumpHistoryFrom(t *testing.T) {
	tc := testCase{"pump-records", 523, 0}
	file := testFileName(tc) + ".json"
	family := testPumpFamily(tc)
	records, err := decodeFromData(file, family)
	if err != nil {
		t.Error(err)
		return
	}
	cases := []struct {
		id    string
		index int
	}{
		{"ewgAgAsZFBYjAA==", 0},
		{"AwAAAAAytygZFA==", 8},
		{"ewAyuwEYFAAYAA==", 296},
		{"xyzzy", -1},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			recordID, err := base64.StdEncoding.DecodeString(c.id)
			if err != nil {
				// Allow illegal base64 strings for nonexistent record IDs.
				recordID = nil
			}
			results, found := findID(records, recordID)
			if !found {
				if c.index != -1 {
					t.Errorf("findID(%s) not found, want %d", c.id, c.index)
					return
				} else if len(results) != len(records) {
					t.Errorf("findID(%s) returned %d records when not found, want %d", c.id, len(results), len(records))
				}
				return
			}
			if len(results) != c.index {
				t.Errorf("findID(%s) returned %d records, want %d", c.id, len(results), c.index)
			}
		})
	}
}

// Mimic the behavior of Pump.findRecords.
func findRecords(records History, check func(HistoryRecord) bool) (History, bool) {
	results := records
	for i, r := range records {
		if check(r) {
			results = records[:i+1]
			break
		}
	}
	n := len(results)
	if n == 0 {
		return nil, false
	}
	r := results[n-1]
	if check(r) {
		return results[:n-1], true
	}
	return results, false
}

func findSince(records History, since time.Time) History {
	check := func(r HistoryRecord) bool {
		return checkBefore(r, since)
	}
	results, _ := findRecords(records, check)
	return results
}

func findID(records History, id []byte) (History, bool) {
	check := func(r HistoryRecord) bool {
		return checkID(r, id)
	}
	return findRecords(records, check)
}
