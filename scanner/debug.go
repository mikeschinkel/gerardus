//go:build debug

package scanner

func init() {
	var s string
	f := &file{sourceDir: &s}
	f.DebugString()
}
func (f *file) DebugString() string {
	return f.relPath
}
