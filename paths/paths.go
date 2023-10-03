package paths

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

func Absolute(path string) (string, error) {
	// Get the current working directory.
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// If the path is not absolute, join it with the current working directory.
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}

	// Clean up the path, evaluating any . or .. elements.
	path = filepath.Clean(path)

	// Convert to an absolute path.
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// EnsureTrailingSlash ensures that a directory path ends with a trailing slash.
func EnsureTrailingSlash(dir string) string {
	if len(dir) == 0 {
		panic("Directory cannot be empty")
	}
	if dir[len(dir)-1] == filepath.Separator {
		return dir
	}

	// Append a trailing slash.
	return dir + string(filepath.Separator)
}

func Relative(root, path string) string {
	rfp, err := filepath.Rel(root, path)
	if err != nil {
		panic("Could not take relative path")
	}
	return rfp
}

func FileExists(file string) (exists bool, err error) {
	_, err = os.Stat(file)
	if errors.Is(err, fs.ErrNotExist) {
		err = nil
		exists = false
		goto end
	}
	if err != nil {
		goto end
	}
	exists = true
end:
	return exists, err
}
