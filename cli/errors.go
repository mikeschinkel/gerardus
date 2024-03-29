package cli

import (
	"github.com/mikeschinkel/go-serr"
)

var (
	ErrOptionAfterArgs    = serr.New("option found after args").ValidArgs("options", "args")
	ErrNoCommandSpecified = serr.New("no command specified")
	ErrCommandNotValid    = serr.New("command is not valid").ValidArgs("command")
	ErrNoExecFuncFound    = serr.New("no exec func found")

	//ErrNoCLIArgsProvided  = serr.New("no command line args provided")
	//ErrHelpNeeded         = serr.New("help needed")
	//ErrFailedToRunCommand = serr.New("failed to run command").ValidArgs("command")

	ErrTooFewArgsPassed  = serr.New("too few arguments passed").ValidArgs("expected", "got")
	ErrTooManyArgsPassed = serr.New("too many arguments passed").ValidArgs("expected", "got")

	ErrAlreadyExists   = serr.New("already exists").ValidArgs("arg_name", "value")
	ErrDoesNotExist    = serr.New("does not exist").ValidArgs("arg_name", "value")
	ErrDoesNotValidate = serr.New("does not validate").ValidArgs("arg_name", "value")

	ErrEmptyStateNotSatisfied = serr.New("not satisfied").ValidArgs("arg_name", "value")

	// ErrTokenValueCannotBeEmpty omits the word "token" to display a message "value
	// cannot be empty" and omits in order to differentiate from a potential generic
	// `ErrValueCannotBeEmpty` error.
	ErrTokenValueCannotBeEmpty = serr.New("value cannot be empty").ValidArgs(string(ArgType), string(FlagType))
	// ErrTokenValueMustBeEmpty omits the word "token" to display a message "value
	// must be empty" and omits in order to differentiate from a potential generic
	// `ErrValueMustBeEmpty` error.
	ErrTokenValueMustBeEmpty = serr.New("value must be empty").ValidArgs(string(ArgType), string(FlagType))

	// ErrHelpSentinel is used to trigger help.Usage() in main() by just checking for err==nil when returned from an app.Main().
	ErrHelpSentinel = serr.New(string(HelpArg))
)
