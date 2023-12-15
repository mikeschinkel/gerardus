package options

import (
	"os"

	"github.com/mikeschinkel/gerardus/cli"
)

type Options struct {
	includeFilesByExtensions   []string
	excludeFilesByPathContains []string
	sourceDir                  string
	dataFile                   string
	projectName                string
	versionTag                 string
	sourceURL                  string
	repoURL                    string
}

var envPrefix string
var options Options

type Params struct {
	EnvPrefix string
}

func Initialize(params Params) error {
	envPrefix = params.EnvPrefix
	options = Options{
		includeFilesByExtensions:   []string{".go", ".mod"},
		excludeFilesByPathContains: []string{"/internal/", "test"},
	}
	return nil
}

func IncludeFilesByExtensions() []string {
	return options.includeFilesByExtensions
}
func ExcludeFilesByPathContains() []string {
	return options.excludeFilesByPathContains
}

func SetSourceDir(dir *cli.Value) {
	options.sourceDir = dir.String()
}

func SourceDir() string {
	return options.sourceDir
}

func SetDataFile(f *cli.Value) {
	envDB := os.Getenv(envPrefix + "DB")
	if envDB != "" {
		options.dataFile = envDB
		goto end
	}
	options.dataFile = f.String()
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

func SetProjectName(name *cli.Value) {
	options.projectName = name.String()
}

func ProjectName() string {
	return options.projectName
}

func SetRepoURL(url *cli.Value) {
	options.repoURL = url.String()
}

func RepoURL() string {
	return options.repoURL
}

func SetVersionTag(name *cli.Value) {
	options.versionTag = name.String()
}

func VersionTag() string {
	return options.versionTag
}
