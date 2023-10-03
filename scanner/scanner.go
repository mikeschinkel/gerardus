package scanner

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gerardus/options"
	"gerardus/parser"
	"gerardus/paths"
)

type DoScanFunc func(path string) bool

type Scanner struct {
	sourceDir  string
	Files      parser.Files
	DoScanFunc DoScanFunc
}

func NewScanner(srcDir string) *Scanner {
	return NewScannerWithFunc(srcDir, func(string) bool {
		return true
	})
}

func NewScannerWithFunc(srcDir string, f DoScanFunc) *Scanner {
	return &Scanner{
		sourceDir:  srcDir,
		Files:      make(parser.Files, 0),
		DoScanFunc: f,
	}
}

func (s *Scanner) Scan() (_ parser.Files, err error) {
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
	return s.Files, err
}

func (s *Scanner) AddFile(file parser.File) {
	s.Files = append(s.Files, file)
}

func (s *Scanner) scanFile(path string, info os.FileInfo, err error) error {

	if err != nil {
		goto end
	}

	path = paths.Relative(s.sourceDir, path)

	if info.IsDir() {
		goto end
	}

	// Skip files that aren't Go source files or go.mod
	if !slices.Contains(options.IncludeFilesByExtensions, filepath.Ext(path)) {
		goto end
	}

	// Skip files that contain excluded path segments or fragments
	for _, pc := range options.ExcludeFilesByPathContains {
		if strings.Contains(path, pc) {
			goto end
		}
	}

	if !s.DoScanFunc(path) {
		goto end
	}

	s.AddFile(parser.NewFile(path, &s.sourceDir))

end:
	return nil
}
