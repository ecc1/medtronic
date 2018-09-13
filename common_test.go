package medtronic

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Force timezone to match test data.
func init() {
	os.Setenv("TZ", "America/New_York")
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

func parseBytes(hex string) []byte {
	var data []byte
	r := strings.NewReader(hex)
	data, err := readBytes(r)
	if err != nil {
		panic(err)
	}
	return data
}

var layouts = []string{
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04",
}

func parseTime(s string) time.Time {
	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.ParseInLocation(layout, s, time.Local)
		if err == nil {
			return t
		}
	}
	panic(err)
}

func parseTD(s string) TimeOfDay {
	t, err := ParseTimeOfDay(s)
	if err != nil {
		panic(err)
	}
	return t
}

func compareDataToJSON(data interface{}, jsonFile string) (bool, string) {
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
	return diffJSON(jsonFile, tmpfile.Name())
}

func diffJSON(file1, file2 string) (bool, string) {
	// Write JSON in canonical form for comparison.
	canon1 := canonicalJSON(file1)
	defer os.Remove(canon1)
	canon2 := canonicalJSON(file2)
	defer os.Remove(canon2)
	// Find differences.
	cmd := exec.Command("diff", "-u", "--label", file1, "--label", file2, canon1, canon2)
	diffs, err := cmd.Output()
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
		panic(err)
	}
	tmpfile, err := ioutil.TempFile("", "json")
	if err != nil {
		panic(err)
	}
	tmpfile.Write(canon)
	tmpfile.Close()
	return tmpfile.Name()
}
