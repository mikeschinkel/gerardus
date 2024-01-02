package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

func StringDiff(s1, s2 string, pad int) (s string) {
	var start, prefix, suffix1, suffix2, end2 int

	lenS1 := len(s1)
	lenS2 := len(s2)

	switch {
	case lenS1 == 0 && lenS2 == 0:
		goto end
	case lenS1 == 0:
		s = fmt.Sprintf("<<2[%s]2>>", s2)
		goto end
	case lenS2 == 0:
		s = fmt.Sprintf("<<1[%s]1>>", s1)
		goto end
	}

	for s1[prefix] == s2[prefix] {
		prefix++
	}
	if pad > 0 {
		start = prefix - pad
	}

	suffix1 = lenS1
	suffix2 = lenS2
	for {
		if suffix1 < 1 {
			suffix1++
			break
		}
		if suffix2 < 1 {
			suffix2++
			break
		}
		if s1[suffix1-1] != s2[suffix2-1] {
			break
		}
		suffix1--
		suffix2--
	}
	end2 = lenS1
	if pad > 0 {
		try := suffix2 + pad
		if try < end2 {
			end2 = try
		}
	}

	switch {
	case prefix == 0 && suffix1 == len(s1) && suffix1 == suffix2:
		// S1 and S2 begin and end completely differently
		s = fmt.Sprintf("<<1[%s]1>><<2[%s]2>>", s1, s2)
	case prefix < suffix1:
		diff := StringDiff(s1[prefix:suffix1], s2[prefix:suffix2], pad)
		s = fmt.Sprintf("%s%s%s", s1[start:prefix], diff, s2[suffix2:end2])
	}
end:
	return s
}

func processMiddle(s1, s2 string) string {
	var result strings.Builder
	i, j := 0, 0

	for i < len(s1) && j < len(s2) {
		if s1[i] != s2[j] {
			result.WriteString(fmt.Sprintf("<<1<[%c]>1>>", s1[i]))
			i++
			result.WriteString(fmt.Sprintf("<<2<[%c]>2>>", s2[j]))
			j++
		} else {
			i++
			j++
		}
	}

	// Append remaining characters
	for i < len(s1) {
		result.WriteString(fmt.Sprintf("<<1<[%c]>1>>", s1[i]))
		i++
	}
	for j < len(s2) {
		result.WriteString(fmt.Sprintf("<<2<[%c]>2>>", s2[j]))
		j++
	}

	return result.String()
}
