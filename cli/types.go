package cli

import (
	"context"
)

type Context = context.Context
type StringMap map[string]string
type ExecFunc func(Context, *CommandInvoker) error
