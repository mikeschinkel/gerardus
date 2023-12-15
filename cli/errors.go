package cli

import (
	"github.com/mikeschinkel/go-serr"
)

var ErrNoCommandSpecified = serr.New("no command specified")
var ErrNoExecFuncFound = serr.New("no exec func found")
var ErrNoCLIArgsProvided = serr.New("no command line args provided")
var ErrHelpNeeded = serr.New("help needed")
var ErrFailedToRunCommand = serr.New("failed to run command").ValidArgs("command")

var ErrTooFewArgsPassed = serr.New("too few arguments passed").ValidArgs("expected", "got")
var ErrTooManyArgsPassed = serr.New("too many arguments passed").ValidArgs("expected", "got")

var ErrArgCannotBeEmpty = serr.New("argument cannot be empty").ValidArgs("arg_name")
var ErrArgMustBeEmpty = serr.New("argument must empty").ValidArgs("arg_name")
