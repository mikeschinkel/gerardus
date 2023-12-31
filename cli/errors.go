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

	ErrAlreadyExists   = serr.New("already exists")
	ErrDoesNotExist    = serr.New("does not exist")
	ErrDoesNotValidate = serr.New("does not validate")

	// ErrTokenValueCannotBeEmpty omits the word "token" to display a message "value
	// cannot be empty" and omits in order to differentiate from a potential generic
	// `ErrValueCannotBeEmpty` error.
	ErrTokenValueCannotBeEmpty = serr.New("value cannot be empty").ValidArgs(string(ArgType), string(FlagType))
	// ErrTokenValueMustBeEmpty omits the word "token" to display a message "value
	// must be empty" and omits in order to differentiate from a potential generic
	// `ErrValueMustBeEmpty` error.
	ErrTokenValueMustBeEmpty = serr.New("value must be empty").ValidArgs(string(ArgType), string(FlagType))
)
