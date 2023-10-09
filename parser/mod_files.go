package parser

import (
	"gerardus/scanner"
)

type ModFile struct {
	scanner.File
	Content []byte
}

func NewModFile(file scanner.File, content []byte) *ModFile {
	return &ModFile{
		File:    file,
		Content: content,
	}
}
