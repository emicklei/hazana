package hazana

import "flag"

var oRPS = flag.Int("rps", 0, "target number of requests per second, must be greater than zero")
var oAttackTime = flag.Int("attack", 10, "duration of the attack in seconds")
var oRampupTime = flag.Int("ramp", 5, "ramp up time in seconds")
var oMaxAttackers = flag.Int("max", 100, "maximum concurrent attackers")
var oVerbose = flag.Bool("v", false, "verbose logging")

// Config holds settings for a Runner.
type Config struct {
	RPS           int
	AttackTimeSec int
	RampupTimeSec int
	MaxAttackers  int
	Verbose       bool
}

// Validate checks all settings and returns a list of strings with problems.
func (c Config) Validate() (list []string) {
	if c.RPS <= 0 {
		list = append(list, "please set the RPS to a positive number of seconds")
	}
	if c.AttackTimeSec <= 0 {
		list = append(list, "please set the attack time to a positive number of seconds")
	}
	if c.RampupTimeSec <= 0 {
		list = append(list, "please set the attack time to a positive number of seconds")
	}
	if c.MaxAttackers <= 0 {
		list = append(list, "please set a positive maximum number of attackers")
	}
	return
}

// ConfigFromFlags creates and validates a Config for use in a runner.
func ConfigFromFlags() Config {
	flag.Parse()
	return Config{
		RPS:           *oRPS,
		AttackTimeSec: *oAttackTime,
		RampupTimeSec: *oRampupTime,
		Verbose:       *oVerbose,
		MaxAttackers:  *oMaxAttackers,
	}
}
