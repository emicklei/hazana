package hazana

import (
	"fmt"
	"strings"
	"time"
)

var start = time.Now()

// Printf prefixes the line with relative time indicator
func Printf(format string, args ...interface{}) {
	sub := (time.Now().Sub(start) / time.Second) * time.Second
	dur := "+" + sub.String()
	fmt.Printf(rightpad(dur, 6)+" - "+format, args...)
}

func rightpad(s string, size int) string {
	return strings.Repeat(" ", size-len(s)) + s
}
