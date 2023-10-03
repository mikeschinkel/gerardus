package parser

type ModFile struct {
	File
	Content []byte
}

func NewModFile(file File, content []byte) *ModFile {
	return &ModFile{
		File:    file,
		Content: content,
	}
}
