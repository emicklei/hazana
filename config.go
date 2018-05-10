package hazana

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	fRPS            = "rps"
	fAttackTime     = "attack"
	fRampupTime     = "ramp"
	fMaxAttackers   = "max"
	fOutput         = "o"
	fVerbose        = "verbose"
	fSample         = "t"
	fRampupStrategy = "s"
	fDoTimeout      = "timeout"
)

var (
	oRPS            = flag.Int(fRPS, 1, "target number of requests per second, must be greater than zero")
	oAttackTime     = flag.Int(fAttackTime, 60, "duration of the attack in seconds")
	oRampupTime     = flag.Int(fRampupTime, 10, "ramp up time in seconds")
	oMaxAttackers   = flag.Int(fMaxAttackers, 10, "maximum concurrent attackers")
	oOutput         = flag.String(fOutput, "", "output file to write the metrics per sample request index (use stdout if empty)")
	oVerbose        = flag.Bool(fVerbose, false, "produce more verbose logging")
	oSample         = flag.Int(fSample, 0, "test your attack implementation with a number of sample calls. Your program exits after this")
	oRampupStrategy = flag.String(fRampupStrategy, defaultRampupStrategy, "set the rampup strategy, possible values are {linear,exp2}")
	oDoTimeout      = flag.Int(fDoTimeout, 5, "timeout in seconds for each attack call")
)

var fullAttackStartedAt time.Time

// Config holds settings for a Runner.
type Config struct {
	RPS            int               `json:"rps"`
	AttackTimeSec  int               `json:"attackTimeSec"`
	RampupTimeSec  int               `json:"rampupTimeSec"`
	RampupStrategy string            `json:"rampupStrategy"`
	MaxAttackers   int               `json:"maxAttackers"`
	OutputFilename string            `json:"outputFilename,omitempty"`
	Verbose        bool              `json:"verbose"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	DoTimeoutSec   int               `json:"doTimeoutSec"`
}

// Validate checks all settings and returns a list of strings with problems.
func (c Config) Validate() (list []string) {
	if c.RPS <= 0 {
		list = append(list, "please set the RPS to a positive number of seconds")
	}
	if c.AttackTimeSec < 2 {
		list = append(list, "please set the attack time to a positive number of seconds > 1")
	}
	if c.RampupTimeSec < 1 {
		list = append(list, "please set the attack time to a positive number of seconds > 0")
	}
	if c.MaxAttackers <= 0 {
		list = append(list, "please set a positive maximum number of attackers")
	}
	if c.DoTimeoutSec <= 0 {
		list = append(list, "please set the Do() timeout to a positive maximum number of seconds")
	}
	return
}

// timeout is in seconds
func (c Config) timeout() time.Duration {
	return time.Duration(c.DoTimeoutSec) * time.Second
}

func (c Config) rampupStrategy() string {
	if len(c.RampupStrategy) == 0 {
		return defaultRampupStrategy
	}
	return c.RampupStrategy
}

// ConfigFromFlags creates a Config for use in a runner.
func ConfigFromFlags() Config {
	flag.Parse()
	return Config{
		RPS:            *oRPS,
		AttackTimeSec:  *oAttackTime,
		RampupTimeSec:  *oRampupTime,
		RampupStrategy: *oRampupStrategy,
		Verbose:        *oVerbose,
		MaxAttackers:   *oMaxAttackers,
		OutputFilename: *oOutput,
		Metadata:       map[string]string{},
		DoTimeoutSec:   *oDoTimeout,
	}
}

// ConfigFromFile loads a Config for use in a runner.
func ConfigFromFile(named string) Config {
	c := ConfigFromFlags() // always parse flags
	f, err := os.Open(named)
	if err != nil {
		log.Fatal("unable to read configuration", err)
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		log.Fatal("unable to decode configuration", err)
	}
	applyFlagOverrides(&c)
	return c
}

// override with any flag set
func applyFlagOverrides(c *Config) {
	flag.Visit(func(each *flag.Flag) {
		switch each.Name {
		case fRPS:
			c.RPS = *oRPS
		case fAttackTime:
			c.AttackTimeSec = *oAttackTime
		case fRampupTime:
			c.RampupTimeSec = *oRampupTime
		case fVerbose:
			c.Verbose = *oVerbose
		case fMaxAttackers:
			c.MaxAttackers = *oMaxAttackers
		case fOutput:
			c.OutputFilename = *oOutput
		case fDoTimeout:
			c.DoTimeoutSec = *oDoTimeout
		}
	})
}

// GetEnv returns the environment variable value or absentValue if it is missing
func GetEnv(key, absentValue string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		if *oVerbose {
			log.Printf("hazana - environment variable [%s] not set, returning [%s...](%d)\n", key, absentValue[:1], len(absentValue))
		}
		return absentValue
	}
	return v
}

// ReadFile returns the text contents of a file or absentValue if it errored
func ReadFile(name, absentValue string) string {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		if *oVerbose {
			log.Printf("hazana - error reading file [%s], returning [%s...](%d)\n", name, absentValue[:1], len(absentValue))
		}
		return absentValue
	}
	return string(data)
}
