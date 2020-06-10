package hazana

import (
	"fmt"
	"time"
)

var start = time.Now()

// Printf prefixes the line with relative time indicator
func Printf(format string, args ...interface{}) {
	sub := (time.Now().Sub(start) / time.Second) * time.Second
	newargs := append([]interface{}{}, sub.String())
	newargs = append(newargs, args...)
	fmt.Printf("+%v - "+format, newargs...)
}
