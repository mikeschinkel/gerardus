package cli

import (
	"fmt"
	"os"
)

func panicf(msg string, args ...any) {
	panic(fmt.Sprintf(msg, args...))
}

func StdErr(msg string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, msg, args...)
}
