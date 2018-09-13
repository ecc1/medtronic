package medtronic

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

const testDataDir = "testdata"

type testCase struct {
	testBase    string
	modelNumber int
	alternative int
}

func TestDecodeHistoryRecord(t *testing.T) {
	cases := []testCase{
		{"model", 512, 0},
		{"model", 515, 0},
		{"model", 522, 0},
		{"model", 523, 1},
		{"model", 523, 2},
		{"ps2", 522, 1},
		{"ps2", 522, 2},
		{"ps2", 523, 1},
		{"ps2", 523, 2},
		{"ps2", 523, 3},
		{"ps2", 523, 4},
		{"ps2", 523, 5},
		{"ps2", 523, 6},
		{"ps2", 551, 1},
		{"ps2", 551, 2},
		{"ps2", 551, 3},
		{"ps2", 551, 4},
		{"ps2", 554, 1},
		{"ps2", 554, 2},
		{"ps2", 554, 3},
		{"ps2", 554, 4},
		{"ps2", 554, 5},
		{"pump-records", 522, 0},
		{"records", 522, 0},
		{"records", 523, 0},
		{"records", 554, 0},
	}
	for _, c := range cases {
		testFile := testFileName(c)
		t.Run(testFile, func(t *testing.T) {
			jsonFile := testFile + ".json"
			family := testPumpFamily(c)
			records, err := decodeFromData(jsonFile, family)
			if err != nil {
				t.Error(err)
				return
			}
			checkHistory(t, records, jsonFile)
		})
	}
}

func testFileName(c testCase) string {
	s := fmt.Sprintf("%s/%s-%d", testDataDir, c.testBase, c.modelNumber)
	if c.alternative != 0 {
		s += fmt.Sprintf("-%d", c.alternative)
	}
	return s
}

func testPumpFamily(c testCase) Family {
	return Family(c.modelNumber % 100)
}

func decodeFromData(file string, family Family) (History, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
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
	cases := []testCase{
		{"model", 512, 0},
		{"model", 515, 0},
		{"model", 522, 0},
		{"model", 523, 1},
		{"model", 523, 2},
		{"ps2", 522, 1},
		{"ps2", 522, 2},
		{"ps2", 523, 1},
		{"ps2", 523, 2},
		{"ps2", 523, 3},
		{"ps2", 523, 4},
		{"ps2", 523, 5},
		{"ps2", 523, 6},
		{"ps2", 551, 1},
		{"ps2", 551, 2},
		{"ps2", 551, 3},
		{"ps2", 551, 4},
		{"ps2", 554, 1},
		{"ps2", 554, 2},
		{"ps2", 554, 3},
		{"ps2", 554, 4},
		{"ps2", 554, 5},
	}
	for _, c := range cases {
		testFile := testFileName(c)
		t.Run(testFile, func(t *testing.T) {
			family := testPumpFamily(c)
			f, err := os.Open(testFile + ".data")
			if err != nil {
				t.Error(err)
				return
			}
			data, err := readBytes(f)
			f.Close()
			if err != nil {
				t.Error(err)
				return
			}
			decoded, err := DecodeHistory(data, family)
			if err != nil {
				t.Errorf("DecodeHistory(% X, %d) returned %v", data, family, err)
				return
			}
			checkHistory(t, decoded, testFile+".json")
		})
	}
}

func checkHistory(t *testing.T, decoded History, jsonFile string) {
	eq, msg := compareDataToJSON(decoded, jsonFile)
	if !eq {
		t.Errorf("JSON is different:\n%s\n", msg)
	}
}

func (r HistoryRecord) String() string {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func TestTreatments(t *testing.T) {
	cases := []struct {
		records    testCase
		treatments testCase
	}{
		{testCase{"pump-records", 522, 0}, testCase{"pump-treatments", 522, 0}},
	}
	for _, c := range cases {
		testFile := testFileName(c.records)
		t.Run(testFile, func(t *testing.T) {
			recordFile := testFile + ".json"
			family := testPumpFamily(c.records)
			records, err := decodeFromData(recordFile, family)
			if err != nil {
				t.Error(err)
				return
			}
			treatments := Treatments(records)
			treatmentFile := testFileName(c.treatments) + ".json"
			eq, msg := compareDataToJSON(treatments, treatmentFile)
			if !eq {
				t.Errorf("JSON is different:\n%s\n", msg)
			}
		})
	}
}

const jqFilter = "cmd/pumphistory/openaps.jq"

func TestJQFilter(t *testing.T) {
	cases := []testCase{
		{"model", 512, 0},
		{"model", 515, 0},
		{"model", 522, 0},
		{"model", 523, 1},
		{"model", 523, 2},
		{"ps2", 522, 1},
		{"ps2", 522, 2},
		{"ps2", 523, 1},
		{"ps2", 523, 2},
		{"ps2", 523, 3},
		{"ps2", 523, 4},
		{"ps2", 523, 5},
		{"ps2", 523, 6},
		{"ps2", 551, 1},
		{"ps2", 551, 2},
		{"ps2", 551, 3},
		{"ps2", 551, 4},
		{"ps2", 554, 1},
		{"ps2", 554, 2},
		{"ps2", 554, 3},
		{"ps2", 554, 4},
		{"ps2", 554, 5},
		{"pump-records", 522, 0},
		{"records", 522, 0},
		{"records", 523, 0},
		{"records", 554, 0},
	}
	for _, c := range cases {
		testFile := testFileName(c)
		t.Run(testFile, func(t *testing.T) {
			jsonFile := testFile + ".json"
			openAPSFile := testFile + ".openaps"
			f, err := exec.Command("jq", "-f", jqFilter, jsonFile).Output()
			if err != nil {
				panic(err)
			}
			tmpfile, err := ioutil.TempFile("", "json")
			if err != nil {
				panic(err)
			}
			filterOutput := tmpfile.Name()
			defer os.Remove(filterOutput)
			tmpfile.Write(f)
			tmpfile.Close()
			eq, msg := diffJSON(filterOutput, openAPSFile)
			if !eq {
				t.Errorf("JSON is different:\n%s\n", msg)
			}
		})
	}
}
