//go:build debug

package scanner

func init() {
	var s string
	f := &file{sourceDir: &s}
	f.DebugString()
}
func (file *file) DebugString() string {
	return file.relPath
}
