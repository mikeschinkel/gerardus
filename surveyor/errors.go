package surveyor

import (
	"github.com/mikeschinkel/go-serr"
)

var (
	errFailedToReadFile  = serr.New("failed to read file").ValidArgs("filename")
	errFailedToParseFile = serr.New("failed to parse file").ValidArgs("filename")
)
