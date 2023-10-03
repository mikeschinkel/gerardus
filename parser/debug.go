//go:build debug

package parser

func init() {
	var s string
	f1 := &file{sourceDir: &s}
	f1.DebugString()
	f2 := &GoFile{
		File: f1,
	}
	f2.DebugString()
}
func (file *file) DebugString() string {
	return file.relPath
}

func (gf *GoFile) DebugString() string {
	return gf.Path
}

//func (t *_archive.Type) DebugString() string {
//	// Any complicated expression simply won't display in debugger
//	return t.debugString
//}
