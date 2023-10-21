package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var LogFilename string

var options = Options{
	logDir:  ".",
	showLog: false,
}

type ReplacerFunc func(groups []string, a slog.Attr) slog.Attr

var logger *slog.Logger
var replacerFuncs = []ReplacerFunc{}

func AddReplacer(r ReplacerFunc) {
	replacerFuncs = append(replacerFuncs, r)
}

var _ slog.Handler = (*NullSLogHandler)(nil)

func Initialize(appName string) (err error) {
	var h slog.Handler
	var w io.Writer

	LogFilename = strings.ToLower(appName + ".log")

	if options.LogLevel() == LogLevelNone {
		slog.SetDefault(slog.New(NullSLogHandler{}))
		goto end
	}

	w, err = os.OpenFile(options.LogFilepath(), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		goto end
	}
	h = slog.NewJSONHandler(w, &slog.HandlerOptions{
		AddSource: true,
		Level:     options.SLogLevel(),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			for _, r := range replacerFuncs {
				a = r(groups, a)
			}
			return a
		},
	})
	if options.ShowLog() {
		h2 := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     options.SLogLevel(),
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				for _, r := range replacerFuncs {
					a = r(groups, a)
				}
				return a
			},
		})
		h = NewSLogTeeHandler(h, h2)
	}
	logger = slog.New(h)
	slog.SetDefault(logger)

	AddReplacer(func(groups []string, a slog.Attr) slog.Attr {
		var s string
		switch a.Key {
		case "time":
			s = a.Value.String()[:19]

		case "source":
			s = FilepathForLog(a)

		default:
			goto end
		}
		a = slog.String(a.Key, s)
	end:
		return a
	})
	slog.Info(strings.Repeat("=", 50), "status", "starting")
end:
	return err
}

var noRootDirRE = regexp.MustCompile(`^` + SourceRootDir())
var noPrefixRE = regexp.MustCompile(`^github.com/newclarity/`)

func FilepathForLog(a slog.Attr) string {
	s, ok := a.Value.Any().(*slog.Source)
	if !ok {
		return a.Value.String()
	}
	return fmt.Sprintf("%s(): %s:%d",
		noPrefixRE.ReplaceAllString(s.Function, ""),
		noRootDirRE.ReplaceAllString(s.File, "."),
		s.Line,
	)
}

type LogLevelName string

const (
	LogLevelDebug LogLevelName = "debug"
	LogLevelInfo  LogLevelName = "info"
	LogLevelWarn  LogLevelName = "warn"
	LogLevelError LogLevelName = "error"
	LogLevelNone  LogLevelName = "none"
)

// validLogLevels provides a slice of log level in name (string) vs number form.
func validLogLevels() []LogLevelName {
	return []LogLevelName{
		LogLevelInfo,
		LogLevelWarn,
		LogLevelError,
		LogLevelDebug,
		LogLevelNone,
	}
}

// ValidLogLevelsString provides a human-readable string showing log level
// options for CLI help.
func ValidLogLevelsString() string {
	sb := strings.Builder{}
	lls := validLogLevels()
	eol := len(lls) - 1
	for i, ll := range lls {
		sb.WriteByte('\'')
		sb.Write([]byte(ll))
		sb.WriteByte('\'')
		switch {
		case i == eol-1:
			sb.WriteString(" or ")
		case i < eol:
			sb.WriteString(", ")
		}
	}
	s := sb.String()
	return s
}

func DefaultLogFilepath() (fp string) {
	b := []byte(filepath.Join(options.LogDir(), LogFilename))
	b[0] = '.'
	return string(b)
}

type NullSLogHandler struct{}

func (n NullSLogHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
func (n NullSLogHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}
func (n NullSLogHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return n
}
func (n NullSLogHandler) WithGroup(_ string) slog.Handler {
	return n
}

type Options struct {
	logDir    string
	logLevel  string
	sLogLevel slog.Leveler
	showLog   bool
	// Additional fields as required
}

func (o *Options) LogDir() string {
	return o.logDir
}

func (o *Options) LogLevel() LogLevelName {
	level := os.Getenv("GERARDUS_LOG_LEVEL")
	if len(level) > 0 {
		o.logLevel = level
	}
	if o.logLevel == "" {
		o.logLevel = string(LogLevelInfo)
	}
	return LogLevelName(o.logLevel)
}

func (o *Options) SLogLevel() (level slog.Level) {
	err := level.UnmarshalText([]byte(o.logLevel))
	if err != nil {
		slog.Warn("Invalid value for loglevel", "loglevel", o.LogLevel, "error", err.Error())
		level = slog.LevelInfo
	}
	return level
}
func (o *Options) LogFilepath() string {
	return filepath.Join(o.logDir, LogFilename)
}

func (o *Options) ShowLog() bool {
	return o.showLog
}
