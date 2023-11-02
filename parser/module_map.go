package parser

type ModuleMap map[string]*ModuleVersions

func (mm ModuleMap) ImportPaths() (ips []string) {
	ips = make([]string, 0, len(mm))
	for ip := range mm {
		ips = append(ips, ip)
	}
	return ips
}
