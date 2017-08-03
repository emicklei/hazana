package hazana

import "flag"

var oRPS = flag.Int("rps", 20, "target number of requests per second")
var oAttackTime = flag.Int("attack", 10, "duration of the attack in seconds")
var oRampupTime = flag.Int("ramp", 5, "ramp up time in seconds")
var oMaxAttackers = flag.Int("max", 100, "maximum concurrent attackers")
var oVerbose = flag.Bool("v", false, "verbose logging")

type Config struct {
	RPS           int
	AttackTimeSec int
	RampupTimeSec int
	MaxAttackers  int
	Verbose       bool
}

func ConfigFromFlags() Config {
	flag.Parse()
	return Config{
		RPS:           *oRPS,
		AttackTimeSec: *oAttackTime,
		RampupTimeSec: *oRampupTime,
		Verbose:       *oVerbose,
	}
}
