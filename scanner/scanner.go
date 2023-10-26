package scanner

import (
	"context"
	"io/fs"
	"log/slog"
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
	err = WalkFiles(dir, func(path string, d fs.DirEntry) (err error) {
		path = filepath.Join(path, d.Name())
		return s.scanFile(ctx, path, d)
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

func (s *Scanner) scanFile(ctx context.Context, path string, d fs.DirEntry) (err error) {
	var f File

	path = paths.Relative(s.sourceDir, path)
	slog.Info("Scanning file", "filepath", path)

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
