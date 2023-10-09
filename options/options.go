package options

import (
	"os"
)

const EnvDBVarName = "GERARDUS_DB"

var options = struct {
	includeFilesByExtensions   []string
	excludeFilesByPathContains []string
	sourceDir                  string
	dataFile                   string
}{
	includeFilesByExtensions:   []string{".go", ".mod"},
	excludeFilesByPathContains: []string{"/internal/", "test"},
}

func IncludeFilesByExtensions() []string {
	return options.includeFilesByExtensions
}
func ExcludeFilesByPathContains() []string {
	return options.excludeFilesByPathContains
}

func SetSourceDir(dir string) {
	options.sourceDir = dir
}

func SourceDir() string {
	return options.sourceDir
}

func SetDataFile(f string) {
	envDB := os.Getenv(EnvDBVarName)
	if envDB != "" {
		options.dataFile = envDB
		goto end
	}
	options.dataFile = f
end:
}

func DataFile() string {
	return options.dataFile
}
