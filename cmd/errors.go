package main

import (
	"gerardus/serr"
)

var (
	errNoProjectSpecified        = serr.New("no project specified")
	errSourceURLAppearsInvalid   = serr.New("source URL appears invalid")
	errNoRepoURLSpecified        = serr.New("no repository URL specified")
	errInvalidGitHubRepoURL      = serr.New("not a valid GitHub repo URL")
	errInvalidGitHubRepoRootURL  = serr.New("repo URL does not begin with https://github.com")
	errURLCouldNotBeDereferenced = serr.New("URL could not be dereferenced")
	errNoVersionTagSpecified     = serr.New("no version tag specified")
	errCodebaseNotAdded          = serr.New("codebase not added")
	errProjectNotFound           = serr.New("project not found")
	errProjectNotAdded           = serr.New("project not added")
	errInvalidCodebaseSourceURL  = serr.New("invalid codebase source URL")
	errAddingCodebase            = serr.New("failed to add new codebase").ValidArgs()
	errReadingSourceDir          = serr.New("failed to read source directory.")
	errDirIsEmpty                = serr.New("directory is empty")
	errPathNotADir               = serr.New("path is not a directory")
)
