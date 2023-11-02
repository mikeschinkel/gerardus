package parser

type ModuleVersions struct {
	VersionMap VersionMap
	Module     *Module
	ModPath    string
	GoMod      bool
}

func NewModuleVersions() *ModuleVersions {
	return &ModuleVersions{
		VersionMap: make(VersionMap),
	}
}

//// PathForVersion returns path for version given a VersionMap.
//func (mv *ModuleVersions) PathForVersion(ver string) (path string, found bool) {
//	for _, m := range mv.VersionMap {
//		if m.VersionName() != ver {
//			continue
//		}
//		if m.VersionName() == "." {
//			path = m.Name()
//			found = true
//			goto end
//		}
//		path = m.String()
//		found = true
//		goto end
//	}
//end:
//	return path, found
//}
