package app

import (
	"github.com/mikeschinkel/gerardus/logger"
)

var _ logger.Opts = (*_opts)(nil)

type Opts interface {
	AppName() string
	EnvPrefix() string
	Args() []string
}

var opts = _opts{
	appName:   AppName,
	envPrefix: EnvPrefix,
}

type _opts struct {
	appName   string
	envPrefix string
	args      []string
}

func (o _opts) AppName() string {
	return o.appName
}
func (o _opts) EnvPrefix() string {
	return o.envPrefix
}
func (o _opts) Args() []string {
	return o.args
}
