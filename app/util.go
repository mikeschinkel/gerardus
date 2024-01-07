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

func DefaultSourceDir(envPrefix string) (dir string) {
	dir = os.Getenv(envPrefix + "SOURCE_DIR")
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

// normalizeVersionTag ensures that a version tag does not begin with 'go' but
// does begin with 'v'.
func normalizeVersionTag(verTag string) string {
	if verTag[:2] == "go" && len(verTag) > 2 {
		// Special case: strip "go" off beginning
		verTag = verTag[2:]
	}
	if len(verTag) > 0 && verTag[1] != 'v' {
		// Ensure version starts with 'v' for Semver
		verTag = "v" + verTag
	}
	return verTag
}
