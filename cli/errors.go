package cli

import (
	"github.com/mikeschinkel/go-serr"
)

var (
	ErrNoCommandSpecified = serr.New("no command specified")
	ErrNoExecFuncFound    = serr.New("no exec func found")

	ErrNoCLIArgsProvided  = serr.New("no command line args provided")
	ErrHelpNeeded         = serr.New("help needed")
	ErrFailedToRunCommand = serr.New("failed to run command").ValidArgs("command")

	ErrTooFewArgsPassed  = serr.New("too few arguments passed").ValidArgs("expected", "got")
	ErrTooManyArgsPassed = serr.New("too many arguments passed").ValidArgs("expected", "got")
	ErrArgCannotBeEmpty  = serr.New("argument cannot be empty").ValidArgs("arg_name")
	ErrArgMustBeEmpty    = serr.New("argument must be empty").ValidArgs("arg_name")

	ErrAlreadyExists   = serr.New("already exists")
	ErrDoesNotExist    = serr.New("does not exist")
	ErrDoesNotValidate = serr.New("does not validate")
)
