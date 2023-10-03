package options

import (
	"fmt"
	"os"
)

func StdErr(msg string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, msg, args...)
}

func fail(msg string, args ...any) {
	StdErr(msg, args...)
	os.Exit(2)
}
