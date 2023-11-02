package parser_test

import (
	"os"
	"path/filepath"
	"testing"
)

func notEquals(t *testing.T, name string, got, want any) {
	t.Helper()
	if got == want {
		t.Errorf("%s: got %v, NOT want %v", name, got, want)
	}
}
func equals(t *testing.T, name string, got, want any) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v, want %v", name, got, want)
	}
}

func strPtrEquals(t *testing.T, name string, got *string, want string) {
	t.Helper()
	if got == nil {
		s := ""
		got = &s
	}
	equals(t, name, *got, want)
}

func strPtr(s string) *string {
	return &s
}

func rootPath(relPath string) string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, relPath)
}
