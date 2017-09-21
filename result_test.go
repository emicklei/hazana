package hazana

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRunReportWithError(t *testing.T) {
	f := filepath.Join(os.TempDir(), "report.json")
	r := NewErrorReport(errors.New("something broke"), Config{OutputFilename: f})
	PrintReport(r)
	data, _ := ioutil.ReadFile(f)
	b := RunReport{}
	if err := json.Unmarshal(data, &b); err != nil {
		t.Log(err)
	}
	t.Log(f)
}
