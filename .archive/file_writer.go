package _archive

//import (
//	"os"
//	"path/filepath"
//
//	"gerardus/parser"
//)
//
//type FileWriter struct {
//	file      parser.File
//	outputDir string
//}
//
//func NewFileWriter(file parser.File, dir string) *FileWriter {
//	return &FileWriter{
//		file:      file,
//		outputDir: dir,
//	}
//}
//
//func (w *FileWriter) WriteFile(content string) (err error) {
//	fp := filepath.Join(w.outputDir, w.file.RelPath())
//	err = os.MkdirAll(filepath.Dir(fp), os.ModePerm)
//	if err != nil {
//		goto end
//	}
//	err = os.WriteFile(fp, []byte(content), os.ModePerm)
//	if err != nil {
//		goto end
//	}
//end:
//	return err
//}
