package scanner

import (
	"path/filepath"
	"strings"
)

var _ File = (*file)(nil)

type Files []File
type FilesMap map[string]File

type File interface {
	Fullpath() string
	SourceDir() string
	RelPath() string
	RelDir(string) (string, bool)
	Filename() string
}

type file struct {
	relPath   string
	sourceDir *string
}

func NewFile(path string, sourceDir *string) File {
	return &file{
		relPath:   path,
		sourceDir: sourceDir,
	}
}

func (f *file) RelDir(rootDir string) (string, bool) {
	_, rd, found := strings.Cut(
		filepath.Dir(f.Fullpath()),
		rootDir+string(filepath.Separator),
	)
	return rd, found
}

func (f *file) RelPath() string {
	return f.relPath
}

func (f *file) SourceDir() string {
	return *f.sourceDir
}

func (f *file) Fullpath() string {
	return filepath.Join(*f.sourceDir, f.relPath)
}

func (f *file) Filename() string {
	return filepath.Base(f.relPath)
}
