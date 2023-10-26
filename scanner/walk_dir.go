package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

var omitRegex = regexp.MustCompile(`^\.(git|idea)$`)

// WalkFiles WalkDir is similar to filepath.WalkDir(), but processes files before
// directories at each directory level.
func WalkFiles(root string, f func(string, fs.DirEntry) error) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	fileCnt, dirCnt := 0, 0
	for _, entry := range entries {
		if omitRegex.MatchString(entry.Name()) {
			continue
		}
		if entry.IsDir() {
			dirCnt++
			continue
		}
		fileCnt++
	}

	dirs := make([]fs.DirEntry, 0, dirCnt)
	files := make([]fs.DirEntry, 0, fileCnt)

	for _, entry := range entries {
		switch {
		case omitRegex.MatchString(entry.Name()):
			continue
		case entry.IsDir():
			dirs = append(dirs, entry)
		default:
			files = append(files, entry)
		}
	}

	// Process files first
	for _, fe := range files {
		err = f(root, fe)
		if err != nil {
			goto end
		}
	}

	// Then process directories
	for _, de := range dirs {
		p := filepath.Join(root, de.Name())
		err = WalkFiles(p, f)
		if err != nil {
			goto end
		}
	}
end:
	return err
}
