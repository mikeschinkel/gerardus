//go:build debug

package parser

func init() {
	ModFile{}.DebugString()
	GoFile{}.DebugString()
	Module{}.DebugString()
	Package{}.DebugString()
}

func (mf ModFile) DebugString() string {
	return mf.debugString
}
func (gf GoFile) DebugString() string {
	return gf.debugString
}

func (m Module) DebugString() string {
	return m.debugString
}

func (p Package) DebugString() string {
	return p.debugString
}
