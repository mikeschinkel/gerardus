package cli

type StringMap map[string]string
type ExecFunc func(*CommandInvoker) error
