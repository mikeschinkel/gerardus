package scanner

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"gerardus/channels"
	"gerardus/options"
	"gerardus/paths"
)

type DoScanFunc func(path string) bool

type ScanMode int

const (
	scanModeInvalid ScanMode = iota
	scanToSlice
	scanToChan
)

type Scanner struct {
	sourceDir  string
	files      Files
	filesChan  chan<- File
	DoScanFunc DoScanFunc
	ScanMode   ScanMode
	match      *regexp.Regexp
}

func NewScanner(srcDir string) *Scanner {
	return NewScannerWithFunc(srcDir, func(string) bool {
		return true
	})
}

func NewScannerWithFunc(srcDir string, f DoScanFunc) *Scanner {
	return &Scanner{
		sourceDir:  srcDir,
		files:      make(Files, 0),
		DoScanFunc: f,
	}
}

func (s *Scanner) Scan(ctx context.Context, match *regexp.Regexp) (_ Files, err error) {
	s.ScanMode = scanToSlice
	s.match = match
	err = s.scan(ctx)
	return s.files, err
}
func (s *Scanner) ScanChan(ctx context.Context, match *regexp.Regexp, ch chan<- File) (err error) {
	s.filesChan = ch
	defer close(s.filesChan)
	s.match = match
	s.ScanMode = scanToChan
	err = s.scan(ctx)
	return err
}

func (s *Scanner) scan(ctx context.Context) (err error) {
	var dir string

	dir, err = paths.Absolute(s.sourceDir)
	if err != nil {
		goto end
	}
	dir = paths.EnsureTrailingSlash(dir)
	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		return s.scanFile(ctx, path, info, err)
	})
	if err != nil {
		goto end
	}

end:
	return err
}

func (s *Scanner) AddFile(file File) {
	s.files = append(s.files, file)
}

func (s *Scanner) scanFile(ctx context.Context, path string, info os.FileInfo, err error) error {
	var f File

	if err != nil {
		goto end
	}

	path = paths.Relative(s.sourceDir, path)
	slog.Info("Scanning file", "filepath", path)

	if info.IsDir() {
		goto end
	}

	if s.match != nil && !s.match.MatchString(path) {
		goto end
	}

	// Skip files that aren't Go source files or go.mod
	if !slices.Contains(options.IncludeFilesByExtensions(), filepath.Ext(path)) {
		goto end
	}

	// Skip files that contain excluded path segments or fragments
	for _, pc := range options.ExcludeFilesByPathContains() {
		if strings.Contains(path, pc) {
			goto end
		}
	}

	if !s.DoScanFunc(path) {
		goto end
	}

	f = NewFile(path, &s.sourceDir)
	switch s.ScanMode {
	case scanToSlice:
		s.AddFile(f)
	case scanToChan:
		err = channels.WriteTo(ctx, s.filesChan, f)
		if err != nil {
			goto end
		}
	}

end:
	return err
}
