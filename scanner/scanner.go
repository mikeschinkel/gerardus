package scanner

import (
	"log/slog"
	"os"
	"path/filepath"
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

func (s *Scanner) Scan() (_ Files, err error) {
	s.ScanMode = scanToSlice
	err = s.scan()
	return s.files, err
}
func (s *Scanner) ScanChan(ch chan<- File) (err error) {
	s.filesChan = ch
	s.ScanMode = scanToChan
	err = s.scan()
	return err
}

func (s *Scanner) scan() (err error) {
	var dir string

	dir, err = paths.Absolute(s.sourceDir)
	if err != nil {
		goto end
	}
	s.sourceDir = paths.EnsureTrailingSlash(dir)

	err = filepath.Walk(s.sourceDir, s.scanFile)
	if err != nil {
		goto end
	}

end:
	return err
}

func (s *Scanner) AddFile(file File) {
	s.files = append(s.files, file)
}

func (s *Scanner) scanFile(path string, info os.FileInfo, err error) error {
	var f File

	if err != nil {
		goto end
	}

	path = paths.Relative(s.sourceDir, path)
	slog.Info("Scanning file", "filepath", path)

	if info.IsDir() {
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
		s.filesChan <- f
	}

end:
	return nil
}
