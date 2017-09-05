package hazana

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"
)

var (
	oRPS                = flag.Int("rps", 1, "target number of requests per second, must be greater than zero")
	oAttackTime         = flag.Int("attack", 60, "duration of the attack in seconds")
	oRampupTime         = flag.Int("ramp", 10, "ramp up time in seconds")
	oMaxAttackers       = flag.Int("max", 10, "maximum concurrent attackers")
	oOutput             = flag.String("o", "", "output file to write the metrics per sample request index (use stdout if empty)")
	oVerbose            = flag.Bool("v", false, "verbose logging")
	oSample             = flag.Int("t", 0, "test your attack implementation with a number of sample calls. Your program exits after this.")
	oRampupStrategy     = flag.String("s", defaultRampupStrategy, "possible values are {linear,exp2}")
	fullAttackStartedAt time.Time
)

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
	return
}

func (c Config) rampupStrategy() string {
	if len(c.RampupStrategy) == 0 {
		return defaultRampupStrategy
	}
	return c.RampupStrategy
}

// ConfigFromFlags creates a Config for use in a runner.
func ConfigFromFlags() Config {
	flag.Parse() // always parse flags
	return Config{
		RPS:            *oRPS,
		AttackTimeSec:  *oAttackTime,
		RampupTimeSec:  *oRampupTime,
		RampupStrategy: *oRampupStrategy,
		Verbose:        *oVerbose,
		MaxAttackers:   *oMaxAttackers,
		OutputFilename: *oOutput,
		Metadata:       map[string]string{},
	}
}

// ConfigFromFile loads a Config for use in a runner.
func ConfigFromFile(named string) (c Config) {
	flag.Parse() // always parse flags
	f, err := os.Open(named)
	if err != nil {
		log.Fatal("unable to read configuration", err)
	}
	defer f.Close()
	json.NewDecoder(f).Decode(&c)
	return
}
