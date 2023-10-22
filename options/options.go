package options

import (
	"os"
)

type Opts interface {
	EnvPrefix() string
}

var opts Opts

func Initialize(o Opts) error {
	opts = o
	return nil
}

var options = struct {
	includeFilesByExtensions   []string
	excludeFilesByPathContains []string
	sourceDir                  string
	dataFile                   string
	projectName                string
	versionTag                 string
	sourceURL                  string
	repoURL                    string
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
	envDB := os.Getenv(opts.EnvPrefix() + "DB")
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

func SetSourceURL(url string) {
	options.sourceURL = url
}

func SourceURL() string {
	return options.sourceURL
}

func SetProjectName(name string) {
	options.projectName = name
}

func ProjectName() string {
	return options.projectName
}

func SetRepoURL(url string) {
	options.repoURL = url
}

func RepoURL() string {
	return options.repoURL
}

func SetVersionTag(name string) {
	options.versionTag = name
}

func VersionTag() string {
	return options.versionTag
}
