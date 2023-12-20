package persister

import (
	"github.com/mikeschinkel/go-serr"
)

var (
	ErrFailedToInsertSpec     = serr.New("failed to insert Spec")
	ErrFailedWhilePersisting  = serr.New("failed while persisting")
	ErrInvalidGitHubRepoURL   = serr.New("invalid Github URL").ValidArgs("repo_url")
	ErrHTTPRequestFailed      = serr.New("failed HTTP request").ValidArgs("status_code", "request_url")
	ErrValueCannotBeEmpty     = serr.New("value cannot be empty").ValidArgs("which_value")
	ErrFailedConvertToAbsPath = serr.New("failed to convert to absolute path").ValidArgs("filepath")
	ErrFailedToInitDataStore  = serr.New("failed to initialize data store").ValidArgs("data_file")

	ErrFailedToReadHTTPResponseBody = serr.New("failed to read HTTP response body").ValidArgs("request_url")
	ErrFailedToUnmarshalJSON        = serr.New("failed to unmarshal the JSON").ValidArgs("source")
)
