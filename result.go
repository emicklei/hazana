package hazana

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
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
	RunError      string              `json:"error"`
	Metrics       map[string]*Metrics `json:"metrics"`
}

// NewErrorReport returns a report when a Run could not be called or executed.
func NewErrorReport(err error, config Config) RunReport {
	return RunReport{
		StartedAt:     time.Now(),
		FinishedAt:    time.Now(),
		RunError:      err.Error(),
		Configuration: config,
	}
}

// PrintReport writes the JSON report to a file or stdout, depending on the configuration.
func PrintReport(r RunReport) {
	// make secrets in Metadata unreadable
	for k := range r.Configuration.Metadata {
		if strings.HasSuffix(k, "*") {
			r.Configuration.Metadata[k] = "***---***---***"
		}
	}
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
	// if verbose and filename is given
	if len(r.Configuration.OutputFilename) > 0 && r.Configuration.Verbose {
		os.Stdout.Write(data)
	}
}

// PrintSummary logs a subset of the report for each metric label
func PrintSummary(r RunReport) {
	for k, v := range r.Metrics {
		log.Println("---------")
		log.Println(k)
		log.Println("- - - - -")
		log.Println("requests:", v.Requests)
		log.Println("     rps:", v.Rate)
		log.Println("    mean:", v.Latencies.Mean)
		log.Println("    95th:", v.Latencies.P95)
		log.Println("     max:", v.Latencies.Max)
		log.Println("  errors:", len(v.Errors))
	}
}
