package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// StderrWriter is the writer used to direct output to StdErr, but can be set to
// a bytes.Buffer to capture output during tests.
var StderrWriter io.Writer = os.Stderr

func panicf(msg string, args ...any) {
	panic(fmt.Sprintf(msg, args...))
}

func StdErr(msg string, args ...any) {
	_, _ = fmt.Fprintf(StderrWriter, msg, args...)
}

func CheckURL(url string) (err error) {
	var resp *http.Response
	var status int

	resp, err = http.Get(url)
	if err != nil {
		goto end
	}
	defer Close(resp.Body, WarnOnError) // Make sure to close the response body.

	status = resp.StatusCode
	if status < 200 || status > 299 {
		err = fmt.Errorf("received HTTP status code %d from %s", status, url)
		goto end
	}
end:
	return err
}

func Close(c io.Closer, f func(err error)) {
	f(c.Close())
}

func WarnOnError(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}
}

func ExecutableFilepath(name string) string {
	dir, err := os.Getwd()
	if err != nil {
		panic("Cannot access current directory")
	}
	return filepath.Join(dir, "bin", name)
}
