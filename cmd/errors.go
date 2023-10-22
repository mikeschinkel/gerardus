package main

import (
	"gerardus/serr"
)

//goland:noinspection GoUnusedGlobalVariable
var (
	errSourceURLAppearsInvalid = serr.New("source URL appears invalid")
)

var (
	errNoProjectSpecified        = serr.New("no project specified")
	errNoRepoURLSpecified        = serr.New("no repository URL specified")
	errInvalidGitHubRepoURL      = serr.New("not a valid GitHub repo URL")
	errInvalidGitHubRepoRootURL  = serr.New("repo URL does not begin with https://github.com")
	errURLCouldNotBeDereferenced = serr.New("URL could not be dereferenced")
	errNoVersionTagSpecified     = serr.New("no version tag specified")
	errFailedToAddCodebase       = serr.New("failed to add codebase")
	errProjectNotFound           = serr.New("project not found")
	errFailedToAddProject        = serr.New("failed to add project").ValidArgs("project", "repo_url")
	errInvalidCodebaseSourceURL  = serr.New("invalid codebase source URL")
	errAddingCodebase            = serr.New("failed to add new codebase").ValidArgs()
	errReadingSourceDir          = serr.New("failed to read source directory.")
	errDirIsEmpty                = serr.New("directory is empty")
	errPathNotADir               = serr.New("path is not a directory")
	errMapCommandFailed          = serr.New("`map` command failed").ValidArgs("source_dir")
	errFailedConvertingToAbsPath = serr.New("failed to convert directory to absolute path").ValidArgs("path")
)
