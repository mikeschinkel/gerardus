package gerardus

import (
	"fmt"
)

//goland:noinspection GoUnusedFunction
func panicf(msg string, args ...any) {
	panic(fmt.Sprintf(msg, args...))
}

//goland:noinspection GoUnusedFunction, GoUnusedParameter
func debugBreakpointHere(...any) {
	// just a function for debugging
}
