package hazana

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// RunReport is a composition of configuration, measurements and custom output from a load run.
type RunReport struct {
	StartedAt     time.Time `json:"startedAt"`
	FinishedAt    time.Time `json:"finishedAt"`
	Configuration Config    `json:"configuration"`
	// RunError is set when a Run could not be called or executed.
	RunError string              `json:"runError"`
	Metrics  map[string]*Metrics `json:"metrics"`
	// Failed can be set by your load test program to indicate that the results are not acceptable.
	Failed bool `json:"failed"`
	// Output is used to publish any custom output in the report.
	Output map[string]interface{} `json:"output"`
}

// NewErrorReport returns a report when a Run could not be called or executed.
func NewErrorReport(err error, config Config) *RunReport {
	r := &RunReport{
		StartedAt:     time.Now(),
		FinishedAt:    time.Now(),
		Configuration: config,
		Failed:        err != nil, // clearly the run was not acceptable
		Output:        map[string]interface{}{},
	}
	if err != nil {
		r.RunError = err.Error()
	}
	return r
}

// PrintReport writes report
// - JSON report to a file
// - CSV report to a file
// - stdout
// depending on the configuration.
func PrintReport(r *RunReport) {
	// make secrets in Metadata unreadable
	for k := range r.Configuration.Metadata {
		if strings.HasSuffix(k, "*") {
			r.Configuration.Metadata[k] = "***---***---***"
		}
	}
	data, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		Printf("unable to create report [%v]", err)
		os.Exit(1)
	}
	if len(r.Configuration.OutputFilename) > 0 {
		if err := ioutil.WriteFile(r.Configuration.OutputFilename, data, os.ModePerm); err != nil {
			Printf("unable to write report [%v]", err)
			os.Exit(1)
		} else {
			Printf("JSON report written to [%s]\n", r.Configuration.OutputFilename)
		}
	}
	if len(r.Configuration.CSVOutputFilename) > 0 {
		PrintCSVReport(r, r.Configuration.CSVOutputFilename)
		Printf("CSV report written to [%s]\n", r.Configuration.CSVOutputFilename)
	}
	if len(r.Configuration.OutputFilename) == 0 &&
		len(r.Configuration.CSVOutputFilename) == 0 {
		// stdout
		fmt.Println(string(data))
	}
}

// PrintSummary logs a subset of the report for each metric label
func PrintSummary(r *RunReport) {
	for k, v := range r.Metrics {
		fmt.Println("---------")
		fmt.Println(k)
		fmt.Println("- - - - -")
		fmt.Println("requests:", v.Requests)
		fmt.Println("  errors:", v.Requests-v.success)
		fmt.Println("     rps:", v.Rate)
		fmt.Println("    mean:", v.Latencies.Mean)
		fmt.Println("    50th:", v.Latencies.P50)
		fmt.Println("    95th:", v.Latencies.P95)
		fmt.Println("    99th:", v.Latencies.P99)
		if v.Requests > 0 {
			fmt.Println("avg kB >:", v.BytesOut/v.Requests/1000)
			fmt.Println("avg kB <:", v.BytesIn/v.Requests/1000)
		}
		fmt.Println("     max:", v.Latencies.Max)
		fmt.Println(" success:", v.successLogEntry(), "%")
	}
}

// PrintCSVReport writes the metrics in CSV format
func PrintCSVReport(r *RunReport, filename string) {
	out, err := os.Create(filename)
	if err != nil {
		Printf("unable to create CSV files [%v]", err)
		return
	}
	defer out.Close()
	w := csv.NewWriter(out)
	header := []string{
		"labels",
		"latencies.total.ms",
		"latencies.mean.ms",
		"latencies.50th.ms",
		"latencies.95th.ms",
		"latencies.99th.ms",
		"latencies.max.ms",
		"wait.ms",
		"requests",
		"rate",
		"succes",
		"status_codes",
		"errors",
		"mean.bytes_out",
		"mean.bytes_in",
	}
	w.Write(header)
	fint64 := func(i int64) string {
		return strconv.FormatInt(i, 10)
	}
	for k, v := range r.Metrics {
		reqCount := int64(v.Requests)
		var meanBytesOut int64 = 0
		var meanBytesIn int64 = 0
		if reqCount > 0 {
			meanBytesOut = int64(v.BytesOut) / reqCount
			meanBytesIn = int64(v.BytesIn) / reqCount
		}
		w.Write([]string{
			k,
			fint64(v.Latencies.Total.Milliseconds()),
			fint64(v.Latencies.Mean.Milliseconds()),
			fint64(v.Latencies.P50.Milliseconds()),
			fint64(v.Latencies.P95.Milliseconds()),
			fint64(v.Latencies.P99.Milliseconds()),
			fint64(v.Latencies.Max.Milliseconds()),
			fint64(v.Wait.Milliseconds()),
			fint64(reqCount),
			fmt.Sprintf("%4.f", v.Rate),
			fint64(int64(v.successLogEntry())),
			fmt.Sprintf("%v", v.StatusCodes),
			fmt.Sprintf("%d", len(v.Errors)),
			fint64(meanBytesOut),
			fint64(meanBytesIn),
		})
	}
	w.Flush()
}
