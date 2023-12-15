package app

import (
	"github.com/mikeschinkel/gerardus/cli"
	"github.com/mikeschinkel/gerardus/persister"
)

const (
	AppName   = "gerardus"
	EnvPrefix = "GERARDUS_"
)

const (
	ProjectArg    cli.ArgName = "project"
	VersionTagArg cli.ArgName = "version_tag"
	SourceURLArg  cli.ArgName = "source_url"
	WebsiteArg    cli.ArgName = "website"
	AboutArg      cli.ArgName = "about"

	RepoURLArg cli.ArgName = persister.RepoURLArg
)
