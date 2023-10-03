package parser

import (
	"sort"
)

type Imports []*Import
type ImportsMap map[string]*Import
type Import struct {
	Package *Package
	Alias   string
}

func NewImport(pkg *Package, alias string) *Import {
	return &Import{
		Package: pkg,
		Alias:   alias,
	}
}

func (ii ImportsMap) Sorted() Imports {
	var sorted = make(Imports, len(ii))
	var n int
	if len(ii) == 0 {
		goto end
	}
	for _, i := range ii {
		sorted[n] = i
		n++
	}
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Package.Name < sorted[j].Package.Name
	})
end:
	return sorted

}
