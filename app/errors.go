package app

import (
	"github.com/mikeschinkel/go-serr"
)

//goland:noinspection GoUnusedGlobalVariable
var (
	ErrSourceURLAppearsInvalid = serr.New("source URL appears invalid")
)

var (
	ErrNoProjectSpecified        = serr.New("no project specified")
	ErrNoRepoURLSpecified        = serr.New("no repository URL specified")
	ErrInvalidGitHubRepoURL      = serr.New("not a valid GitHub repo URL")
	ErrInvalidGitHubRepoRootURL  = serr.New("repo URL does not begin with https://github.com")
	ErrURLCouldNotBeDereferenced = serr.New("URL could not be dereferenced")
	ErrNoVersionTagSpecified     = serr.New("no version tag specified")
	ErrVersionAlreadyExists      = serr.New("version already exists").ValidArgs("project", "version_tag")
	ErrFailedToAddCodebase       = serr.New("failed to add codebase")
	ErrProjectNotFound           = serr.New("project not found")
	ErrFailedToAddProject        = serr.New("failed to add project").ValidArgs("project", "repo_url")
	ErrInvalidCodebaseSourceURL  = serr.New("invalid codebase source URL")
	ErrAddingCodebase            = serr.New("failed to add new codebase").ValidArgs()
	ErrReadingSourceDir          = serr.New("failed to read source directory.")
	ErrDirIsEmpty                = serr.New("directory is empty")
	ErrPathNotADir               = serr.New("path is not a directory")
	ErrMapCommandFailed          = serr.New("`map` command failed").ValidArgs("source_dir")
	ErrFailedConvertingToAbsPath = serr.New("failed to convert directory to absolute path").ValidArgs("path")
	ErrVersionIsNotValid         = serr.New("version is not valid").ValidArgs("version_tag")
)
