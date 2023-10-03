package collector

import (
	"fmt"
)

func panicf(msg string, args ...any) {
	panic(fmt.Sprintf(msg, args...))
}

func debugBreakpointHere(...any) {
	// just a function for debugging
}
