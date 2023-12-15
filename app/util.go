package app

import (
	"os"
	"path/filepath"

	"github.com/mikeschinkel/gerardus/paths"
)

func makeAbs(path string) (string, error) {
	absDir, err := paths.Absolute(path)
	if err != nil {
		err = ErrFailedConvertingToAbsPath.Err(err, "path", path)
	}
	return absDir, err
}

func defaultSourceDir(opts Opts) string {
	dir := os.Getenv(opts.EnvPrefix() + "SOURCE_DIR")
	if len(dir) > 0 {
		goto end
	}
	dir = os.Getenv("GOROOT")
	if len(dir) > 0 {
		dir = filepath.Join(dir, "src")
		goto end
	}
	dir = "."
end:
	return dir
}
