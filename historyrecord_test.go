package medtronic

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"testing"
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
		{"testdata/pump-records.json", false},
		{"testdata/model512.json", false},
	}
	for _, c := range cases {
		records, err := decodeFromData(c.jsonFile, c.newer)
		if err != nil {
			t.Errorf("%v", err)
			continue
		}
		checkHistory(t, records, c.jsonFile)
	}
}

func decodeFromData(file string, newer bool) (History, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	d := json.NewDecoder(f)
	var maps []interface{}
	err = d.Decode(&maps)
	if err != nil {
		return nil, err
	}
	var records History
	for _, v := range maps {
		m := v.(map[string]interface{})
		base64data, ok := m["Data"].(string)
		if !ok {
			return records, fmt.Errorf("no data in %+v", v)
		}
		data, err := base64.StdEncoding.DecodeString(base64data)
		if err != nil {
			return records, err
		}
		r, err := DecodeHistoryRecord(data, newer)
		if err != nil {
			return records, err
		}
		records = append(records, r)
	}
	return records, nil
}

func TestDecodeHistory(t *testing.T) {
	cases := []struct {
		pageFile string
		jsonFile string
		newer    bool
	}{
		{"testdata/new1.data", "testdata/new1.json", true},
		{"testdata/new2.data", "testdata/new2.json", true},
		{"testdata/old1.data", "testdata/old1.json", false},
		{"testdata/model512.data", "testdata/model512.json", false},
	}
	for _, c := range cases {
		data, err := readBytes(c.pageFile)
		if err != nil {
			t.Errorf("%v", err)
			continue
		}
		decoded, err := DecodeHistory(data, c.newer)
		if err != nil {
			t.Errorf("DecodeHistory(% X, %v) returned %v", data, c.newer, err)
			continue
		}
		checkHistory(t, decoded, c.jsonFile)
	}
}

func checkHistory(t *testing.T, decoded History, jsonFile string) {
	eq, msg := compareJSON(decoded, jsonFile)
	if !eq {
		t.Errorf("JSON is different:\n%s\n", msg)
	}
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

func compareJSON(data interface{}, jsonFile string) (bool, string) {
	// Write data in JSON format to temporary file.
	tmpfile, err := ioutil.TempFile("", "json")
	if err != nil {
		return false, err.Error()
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()
	e := json.NewEncoder(tmpfile)
	e.SetIndent("", "  ")
	err = e.Encode(data)
	_ = tmpfile.Close()
	if err != nil {
		return false, err.Error()
	}
	// Write JSON in canonical form for comparison.
	canon1 := canonicalJSON(jsonFile)
	canon2 := canonicalJSON(tmpfile.Name())
	// Find differences.
	cmd := exec.Command("diff", "-u", "--label", jsonFile, "--label", "decoded", canon1, canon2)
	diffs, err := cmd.Output()
	_ = os.Remove(canon1)
	_ = os.Remove(canon2)
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
	_, _ = tmpfile.Write(canon)
	_ = tmpfile.Close()
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
		records, err := decodeFromData(c.recordFile, false)
		if err != nil {
			t.Errorf("%v", err)
			continue
		}
		treatments := Treatments(records)
		eq, msg := compareJSON(treatments, c.treatmentFile)
		if !eq {
			t.Errorf("JSON is different:\n%s\n", msg)
		}
	}
}
