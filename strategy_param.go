package hazana

import (
	"log"
	"strconv"
	"strings"
)

// method  key=value key2=value
type strategyParameters struct {
	line string
}

func (c strategyParameters) is(method string) bool {
	return strings.HasPrefix(c.line, method)
}

func (c strategyParameters) intParam(name string, absent int) int {
	tokens := strings.Split(strings.TrimSpace(c.line), " ")
	if len(tokens) <= 1 {
		return absent
	}
	for i := 1; i < len(tokens); i++ {
		kv := strings.Split(tokens[i], "=")
		if len(kv) == 0 {
			continue
		}
		if kv[0] == name {
			i, err := strconv.Atoi(kv[1])
			if err != nil {
				if *oDebug {
					log.Println("parameter fail", err)
				}
				return absent
			}
			return i
		}
	}
	return absent
}

func (c strategyParameters) floatParam(name string, absent float64) float64 {
	tokens := strings.Split(strings.TrimSpace(c.line), " ")
	if len(tokens) <= 1 {
		return absent
	}
	for i := 1; i < len(tokens); i++ {
		kv := strings.Split(tokens[i], "=")
		if len(kv) == 0 {
			continue
		}
		if kv[0] == name {
			f, err := strconv.ParseFloat(kv[1], 64)
			if err != nil {
				if *oDebug {
					log.Println("parameter fail", err)
				}
				return absent
			}
			return f
		}
	}
	return absent
}
