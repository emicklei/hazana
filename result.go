package hazana

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"
)

type result struct {
	begin, end time.Time
	elapsed    time.Duration
	doResult   DoResult
}

// DoResult is the return value of a Do call on an Attack.
type DoResult struct {
	// Label identifying the request that was send which is only used for reporting the metrics.
	RequestLabel string
	// The error that happened when sending the request or receiving the response.
	Error error
	// The HTTP status code.
	StatusCode int
	// Number of bytes transferred when sending the request.
	BytesIn int64
	// Number of bytes transferred when receiving the response.
	BytesOut int64
}

// RunReport is an composition of configuration and measurements from a load run.
type RunReport struct {
	StartedAt     time.Time           `json:"startedAt"`
	FinishedAt    time.Time           `json:"finishedAt"`
	Configuration Config              `json:"configuration"`
	RunError      error               `json:"error"`
	Metrics       map[string]*Metrics `json:"metrics"`
}

// NewErrorReport returns a report when a Run could not be called or executed.
func NewErrorReport(err error, config Config) RunReport {
	return RunReport{
		StartedAt:     time.Now(),
		FinishedAt:    time.Now(),
		RunError:      err,
		Configuration: config,
	}
}

// PrintReport writes the JSON report to a file or stdout, depending on the configuration.
func PrintReport(r RunReport) {
	var out io.Writer
	if len(r.Configuration.OutputFilename) > 0 {
		file, err := os.Create(r.Configuration.OutputFilename)
		if err != nil {
			log.Fatal("unable to create output file", err)
		}
		defer file.Close()
		out = file
	} else {
		out = os.Stdout
	}
	data, _ := json.MarshalIndent(r, "", "\t")
	out.Write(data)
}
