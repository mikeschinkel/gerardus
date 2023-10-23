package scanner

import (
	"path/filepath"
)

var _ File = (*file)(nil)

type Files []File
type FilesMap map[string]File

type File interface {
	Fullpath() string
	SourceDir() string
	RelPath() string
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
