package _archive

//import (
//	"fmt"
//	"strings"
//
//	"gerardus"
//	"gerardus/parser"
//)
//
//type Dumper struct {
//	Files parser.Files
//}
//
//func NewDumper(files parser.Files) *Dumper {
//	return &Dumper{
//		Files: files,
//	}
//}
//
//func (d *Dumper) Dump() {
//	for _, f := range d.Files {
//		d.DumpFile(f)
//	}
//}
//func (d *Dumper) DumpFile(f parser.File) (s string) {
//	rp := f.RelPath()
//	line := strings.Repeat("-", len(rp)+1)
//	fmt.Printf("\n%s\n", line)
//	fmt.Printf("%s:\n", rp)
//	fmt.Printf("%s\n", line)
//
//	switch t := f.(type) {
//	case *parser.GoFile:
//		formatter := NewASTFile(t.ast, t.Package)
//		s = formatter.String()
//	case *parser.ModFile:
//		s = string(t.Content)
//	default:
//		gerardus.panicf("Unhandled File type '%T'", f)
//	}
//	return s
//}
