package parser

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

func (file *file) RelPath() string {
	return file.relPath
}

func (file *file) SourceDir() string {
	return *file.sourceDir
}

func (file *file) Fullpath() string {
	return filepath.Join(*file.sourceDir, file.relPath)
}

func (file *file) Filename() string {
	return filepath.Base(file.relPath)
}
