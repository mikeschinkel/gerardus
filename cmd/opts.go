package main

import (
	"gerardus/logger"
	"gerardus/options"
)

var _ logger.Opts = (*_opts)(nil)
var _ options.Opts = (*_opts)(nil)

type Opts interface {
	AppName() string
	EnvPrefix() string
}

var opts = _opts{
	appName:   AppName,
	envPrefix: EnvPrefix,
}

type _opts struct {
	appName   string
	envPrefix string
}

func (o _opts) AppName() string {
	return o.appName
}
func (o _opts) EnvPrefix() string {
	return o.envPrefix
}
