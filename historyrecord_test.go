package medtronic

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"
)

func TestDecodeHistoryRecord(t *testing.T) {
	cases := []struct {
		jsonFile string
		family   Family
	}{
		{"testdata/model512.json", 12},
		{"testdata/model515.json", 15},
		{"testdata/model522.json", 22},
		{"testdata/model523-1.json", 23},
		{"testdata/model523-2.json", 23},
		{"testdata/ps2-522-1.json", 22},
		{"testdata/ps2-522-2.json", 22},
		{"testdata/ps2-523-1.json", 23},
		{"testdata/ps2-523-2.json", 23},
		{"testdata/ps2-523-3.json", 23},
		{"testdata/ps2-523-4.json", 23},
		{"testdata/ps2-523-5.json", 23},
		{"testdata/ps2-523-6.json", 23},
		{"testdata/ps2-551-1.json", 51},
		{"testdata/ps2-551-2.json", 51},
		{"testdata/ps2-551-3.json", 51},
		{"testdata/ps2-551-4.json", 51},
		{"testdata/ps2-554-1.json", 54},
		{"testdata/ps2-554-2.json", 54},
		{"testdata/ps2-554-3.json", 54},
		{"testdata/ps2-554-4.json", 54},
		{"testdata/ps2-554-5.json", 54},
		{"testdata/pump-records-522.json", 22},
		{"testdata/records-522.json", 22},
		{"testdata/records-523.json", 23},
		{"testdata/records-554.json", 54},
	}
	for _, c := range cases {
		t.Run(c.jsonFile, func(t *testing.T) {
			records, err := decodeFromData(c.jsonFile, c.family)
			if err != nil {
				t.Error(err)
				return
			}
			checkHistory(t, records, c.jsonFile)
		})
	}
}

func decodeFromData(file string, family Family) (History, error) {
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
		r, err := DecodeHistoryRecord(data, family)
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
		family   Family
	}{
		{"testdata/model512.data", "testdata/model512.json", 12},
		{"testdata/model515.data", "testdata/model515.json", 15},
		{"testdata/model522.data", "testdata/model522.json", 22},
		{"testdata/model523-1.data", "testdata/model523-1.json", 23},
		{"testdata/model523-2.data", "testdata/model523-2.json", 23},
		{"testdata/ps2-522-1.data", "testdata/ps2-522-1.json", 22},
		{"testdata/ps2-522-2.data", "testdata/ps2-522-2.json", 22},
		{"testdata/ps2-523-1.data", "testdata/ps2-523-1.json", 23},
		{"testdata/ps2-523-2.data", "testdata/ps2-523-2.json", 23},
		{"testdata/ps2-523-3.data", "testdata/ps2-523-3.json", 23},
		{"testdata/ps2-523-4.data", "testdata/ps2-523-4.json", 23},
		{"testdata/ps2-523-5.data", "testdata/ps2-523-5.json", 23},
		{"testdata/ps2-523-6.data", "testdata/ps2-523-6.json", 23},
		{"testdata/ps2-551-1.data", "testdata/ps2-551-1.json", 51},
		{"testdata/ps2-551-2.data", "testdata/ps2-551-2.json", 51},
		{"testdata/ps2-551-3.data", "testdata/ps2-551-3.json", 51},
		{"testdata/ps2-551-4.data", "testdata/ps2-551-4.json", 51},
		{"testdata/ps2-554-1.data", "testdata/ps2-554-1.json", 54},
		{"testdata/ps2-554-2.data", "testdata/ps2-554-2.json", 54},
		{"testdata/ps2-554-3.data", "testdata/ps2-554-3.json", 54},
		{"testdata/ps2-554-4.data", "testdata/ps2-554-4.json", 54},
		{"testdata/ps2-554-5.data", "testdata/ps2-554-5.json", 54},
	}
	for _, c := range cases {
		t.Run(c.pageFile, func(t *testing.T) {
			f, err := os.Open(c.pageFile)
			if err != nil {
				t.Error(err)
				return
			}
			data, err := readBytes(f)
			_ = f.Close()
			if err != nil {
				t.Error(err)
				return
			}
			decoded, err := DecodeHistory(data, c.family)
			if err != nil {
				t.Errorf("DecodeHistory(% X, %d) returned %v", data, c.family, err)
				return
			}
			checkHistory(t, decoded, c.jsonFile)
		})
	}
}

func checkHistory(t *testing.T, decoded History, jsonFile string) {
	eq, msg := compareJSON(decoded, jsonFile)
	if !eq {
		t.Errorf("JSON is different:\n%s\n", msg)
	}
}

func readBytes(r io.Reader) ([]byte, error) {
	var data []byte
	for {
		var b byte
		n, err := fmt.Fscanf(r, "%02x", &b)
		if n == 0 {
			break
		}
		if err != nil {
			return data, err
		}
		data = append(data, b)
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
		family        Family
	}{
		{"testdata/pump-records-522.json", "testdata/pump-treatments-522.json", 22},
	}
	for _, c := range cases {
		t.Run(c.recordFile, func(t *testing.T) {
			records, err := decodeFromData(c.recordFile, c.family)
			if err != nil {
				t.Error(err)
				return
			}
			treatments := Treatments(records)
			eq, msg := compareJSON(treatments, c.treatmentFile)
			if !eq {
				t.Errorf("JSON is different:\n%s\n", msg)
			}
		})
	}
}
