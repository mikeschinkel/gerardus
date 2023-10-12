package cli

type StringMap map[string]string
type ExecFunc func(StringMap) error
