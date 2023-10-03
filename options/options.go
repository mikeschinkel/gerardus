package options

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gerardus/paths"
)

const SqliteDB = "/Volumes/Tech/SQLiteDBs/go-code-surveys.db"

var SourceDir string
var OutputDir string
var DataFile string
var IncludeFilesByExtensions []string
var ExcludeFilesByPathContains []string

func GetSourceDir() string {
	return SourceDir
}
func GetDataFile() string {
	return DataFile
}

func InitOptions() (err error) {

	if len(os.Args) < 2 || os.Args[len(os.Args)-1] != "run" {
		err = fmt.Errorf("command 'run' not specified: '%s' provided instead",
			strings.Join(os.Args[1:], " "))
		goto end
	}

	//var srcDir, outDir string
	// Parse command line arguments
	flag.StringVar(&SourceDir, "src", defaultSourceDir(), "Source directory")
	flag.StringVar(&OutputDir, "out", "./output", "Output directory")
	flag.StringVar(&DataFile, "data", SqliteDB, "Data file (sqlite3)")
	flag.Parse()

	IncludeFilesByExtensions = []string{".go", ".mod"}
	ExcludeFilesByPathContains = []string{"/internal/", "test"}

	SourceDir, err = checkDir(SourceDir)
	if err != nil {
		goto end
	}
	OutputDir, err = ensureDir(OutputDir)
	if err != nil {
		goto end
	}
end:
	return err
}

func defaultSourceDir() string {
	dir := os.Getenv("GOROOT")
	if len(dir) > 0 {
		dir = filepath.Join(dir, "src")
	}
	return dir
}

// checkDir validates source directory
func checkDir(dir string) (string, error) {
	var info os.FileInfo

	absDir, err := makeAbs(dir)
	if err != nil {
		goto end
	}
	info, err = os.Stat(absDir)
	if err != nil {
		err = fmt.Errorf("error reading source dir: %s; %w", absDir, err)
		goto end
	}
	if !info.IsDir() {
		err = fmt.Errorf("provided source dir is not a directory: %s", absDir)
		goto end
	}
end:
	return absDir, err
}

// ensureDir validates or creates destination directory
func ensureDir(dir string) (string, error) {
	var info os.FileInfo

	absDir, err := makeAbs(dir)
	if err != nil {
		goto end
	}
	info, err = os.Stat(absDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			err = fmt.Errorf("error creating destination directory %s; %w", dir, err)
		}
		goto end
	}
	if err != nil {
		err = fmt.Errorf("error reading destination directory %s: %w", dir, err)
	}
	if !info.IsDir() {
		err = fmt.Errorf("provided output path is not a directory: %s", dir)
	}
end:
	return absDir, err
}

func makeAbs(path string) (string, error) {
	absDir, err := paths.Absolute(path)
	if err != nil {
		err = fmt.Errorf("error converting to absolute path: %s; %w",
			path, err)
	}
	return absDir, err
}
