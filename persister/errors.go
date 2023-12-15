package persister

import (
	"github.com/mikeschinkel/go-serr"
)

var (
	ErrFailedToInsertSpec    = serr.New("failed to insert Spec")
	ErrFailedWhilePersisting = serr.New("failed while persisting")
	ErrInvalidGitHubRepoURL  = serr.New("invalid Github URL").ValidArgs("repo_url")
	ErrHTTPRequestFailed     = serr.New("failed HTTP request").ValidArgs("status_code", "request_url")
	ErrValueCannotBeEmpty    = serr.New("value cannot be empty").ValidArgs("which_value")
)
