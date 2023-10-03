package _archive

//
//import (
//	"errors"
//	"fmt"
//	"os"
//
//	"gerardus"
//	"gerardus/paths"
//	"gerardus/surveyor"
//)
//
//type CodeGenerator struct {
//	dir    string
//	survey *surveyor.CodeSurveyor
//}
//
//func NewCodeGenerator(dir string, survey *surveyor.CodeSurveyor) *CodeGenerator {
//	return &CodeGenerator{
//		dir:    paths.EnsureTrailingSlash(dir),
//		survey: survey,
//	}
//}
//
//func (g *CodeGenerator) Generate() (err error) {
//	dir, err := paths.Absolute(g.dir)
//	if err != nil {
//		goto end
//	}
//	if len(dir) == 0 {
//		err = errors.New("cannot write to an empty directory")
//		goto end
//	}
//	if len(dir) <= 3 {
//		err = fmt.Errorf("directory too short: %s", dir)
//		goto end
//	}
//	g.dir = dir
//
//	// Delete the destination directory
//	err = os.RemoveAll(g.dir)
//	if err != nil {
//		goto end
//	}
//	for _, file := range g.survey.Files {
//		fg := NewFileGenerator(file, g.survey)
//		if fg.ShouldGenerate() {
//			fw := gerardus.NewFileWriter(file, g.dir)
//			err = fw.WriteFile(fg.GenerateContent())
//		}
//		if err != nil {
//			goto end
//		}
//	}
//end:
//	return err
//}
