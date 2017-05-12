package medtronic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"testing"

	"github.com/ecc1/nightscout"
)

func TestDecodeHistoryRecord(t *testing.T) {
	cases := []struct {
		jsonFile string
		newer    bool
	}{
		{"testdata/records-new.json", true},
		{"testdata/records-old.json", false},
		{"testdata/new1.json", true},
		{"testdata/new2.json", true},
		{"testdata/old1.json", false},
	}
	for _, c := range cases {
		records, err := readHistoryRecords(c.jsonFile)
		if err != nil {
			t.Errorf("%v", err)
			continue
		}
		decoded := make([]HistoryRecord, len(records))
		for i, r1 := range records {
			r2, err := DecodeHistoryRecord(r1.Data, c.newer)
			if err != nil {
				t.Errorf("DecodeHistoryRecord(% X, %v) returned %v", r1.Data, c.newer, err)
				continue
			}
			decoded[i] = r2
		}
		if !equalHistoryRecords(t, decoded, records, c.jsonFile) {
			continue
		}
	}
}

func TestDecodeHistoryRecords(t *testing.T) {
	cases := []struct {
		pageFile string
		jsonFile string
		newer    bool
	}{
		{"testdata/new1.data", "testdata/new1.json", true},
		{"testdata/new2.data", "testdata/new2.json", true},
		{"testdata/old1.data", "testdata/old1.json", false},
	}
	for _, c := range cases {
		data, err := readBytes(c.pageFile)
		if err != nil {
			t.Errorf("%v", err)
			continue
		}
		records, err := readHistoryRecords(c.jsonFile)
		if err != nil {
			t.Errorf("%v", err)
			continue
		}
		decoded, err := DecodeHistoryRecords(data, c.newer)
		if err != nil {
			t.Errorf("DecodeHistoryRecords(% X, %v) returned %v", data, c.newer, err)
			continue
		}
		if !equalHistoryRecords(t, decoded, records, c.jsonFile) {
			continue
		}
	}
}

func equalHistoryRecords(t *testing.T, got []HistoryRecord, want []HistoryRecord, jsonFile string) bool {
	for i, r1 := range want {
		r2 := got[i]
		if !reflect.DeepEqual(r1, r2) {
			t.Errorf("got %v, want %v", r2, r1)
			return false
		}
	}
	eq, msg := compareJSON(got, jsonFile)
	if !eq {
		t.Errorf("JSON is different:\n%s\n", msg)
		return false
	}
	return true
}

func readHistoryRecords(file string) ([]HistoryRecord, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	d := json.NewDecoder(f)
	var records []HistoryRecord
	err = d.Decode(&records)
	f.Close() // nolint
	if err != nil {
		err = fmt.Errorf("%s: %v", file, err)
	}
	return records, err
}

func readBytes(file string) ([]byte, error) {
	hex, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	fields := bytes.Fields(hex)
	data := make([]byte, len(fields))
	for i, s := range fields {
		b, err := strconv.ParseUint(string(s), 16, 8)
		if err != nil {
			return nil, err
		}
		data[i] = byte(b)
	}
	return data, nil
}

func (r HistoryRecord) String() string {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

// nolint: errcheck
func compareJSON(data interface{}, jsonFile string) (bool, string) {
	// Write data in JSON format to temporary file.
	tmpfile, err := ioutil.TempFile("", "json")
	if err != nil {
		return false, err.Error()
	}
	defer os.Remove(tmpfile.Name())
	e := json.NewEncoder(tmpfile)
	e.SetIndent("", "  ")
	err = e.Encode(data)
	tmpfile.Close()
	if err != nil {
		return false, err.Error()
	}
	// Write JSON in canonical form for comparison.
	canon1 := canonicalJSON(jsonFile)
	canon2 := canonicalJSON(tmpfile.Name())
	// Find differences.
	cmd := exec.Command("diff", "-u", "--label", jsonFile, "--label", "decoded", canon1, canon2)
	diffs, err := cmd.Output()
	os.Remove(canon1)
	os.Remove(canon2)
	return err == nil, string(diffs)
}

// canonicalJSON reads the given file and creates a temporary file
// containing equivalent JSON in canonical form
// (using the "jq" command, which must be on the user's PATH).
// It returns the temporary file name; it is the caller's responsibility
// to remove it when done.
func canonicalJSON(file string) string {
	canon, err := exec.Command("jq", "-S", ".", file).Output()
	if err != nil {
		log.Fatal(err)
	}
	tmpfile, err := ioutil.TempFile("", "json")
	if err != nil {
		log.Fatal(err)
	}
	tmpfile.Write(canon) // nolint
	tmpfile.Close()      // nolint
	return tmpfile.Name()
}

func TestTreatments(t *testing.T) {
	cases := []struct {
		recordFile    string
		treatmentFile string
	}{
		{"testdata/pump-records.json", "testdata/pump-treatments.json"},
	}
	for _, c := range cases {
		records, err := readHistoryRecords(c.recordFile)
		if err != nil {
			t.Errorf("%v", err)
			continue
		}
		want, err := readTreatments(c.treatmentFile)
		if err != nil {
			t.Errorf("%v", err)
			continue
		}
		got := Treatments(records)
		for i, r1 := range want {
			r2 := got[i]
			if !reflect.DeepEqual(r1, r2) {
				t.Errorf("got %v, want %v", r2, r1)
			}
		}
		eq, msg := compareJSON(got, c.treatmentFile)
		if !eq {
			t.Errorf("JSON is different:\n%s\n", msg)
		}
	}
}

func readTreatments(file string) ([]nightscout.Treatment, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	d := json.NewDecoder(f)
	var records []nightscout.Treatment
	err = d.Decode(&records)
	f.Close() // nolint
	if err != nil {
		err = fmt.Errorf("%s: %v", file, err)
	}
	return records, err
}
