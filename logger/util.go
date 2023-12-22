package logger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

var StderrWriter io.Writer = os.Stderr

func panicf(msg string, args ...any) {
	panic(fmt.Sprintf(msg, args...))
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
		_, _ = fmt.Fprintln(StderrWriter, err.Error())
	}
}

func SourceRootDir() string {
	_, fp, _, ok := runtime.Caller(0)
	if !ok {
		panic("Unable to get full filepath of ./pkg/const.go")
	}
	// First Strip off /const.go then strip off /pkg
	return filepath.Dir(filepath.Dir(fp))
}
