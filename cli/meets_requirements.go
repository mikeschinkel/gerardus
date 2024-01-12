package cli

import (
	"errors"
	"regexp"
	"strings"

	"github.com/mikeschinkel/go-serr"
)

type RequirementsMeeters []RequirementsMeeter
type RequirementsMeeter interface {
	EmptyStateSatisfied(Context, TokenType) error
	CheckExistence(Context) error
	ValidateByFunc(Context) error
}

func onMeetRequirementsPanic() {
	if r := recover(); r != nil {
		switch t := r.(type) {
		case string:
			switch {
			case strings.Contains(t, ": method is nil but DataStoreQueries."):
				// Example: DataStoreQueriesStub.LoadProjectByNameFunc: method is nil but DataStoreQueries.LoadProjectByName was just called
				f := regexp.MustCompile(`^DataStoreQueriesStub\.([^:]+):`).FindStringSubmatch(t)
				if len(f) == 2 {
					panicf("TODO: tt.queries.%s was not yet set in table test data.", f[1])
				}
			}
		}
		panicf("ERROR on cli.MeetRequirements(): %v", r)
	}
}

func MeetsRequirements[S ~[]RM, RM RequirementsMeeter](ctx Context, tt TokenType, rms S) (err error) {
	defer onMeetRequirementsPanic()

	// Check NotExist, MustExist or IgnoreExists requirements
	err = EmptyStateSatisfied(ctx, tt, rms)
	if err != nil {
		if errors.Is(err, ErrEmptyStateNotSatisfied) {
			err = err.(serr.SError).CloneUnwrap()
		}
		goto end
	}

	err = ValidateByFunc(ctx, rms)
	if err != nil {
		if errors.Is(err, ErrDoesNotValidate) {
			err = err.(serr.SError).CloneUnwrap()
		}
		goto end
	}

	err = CheckExistence(ctx, rms)
	if err != nil {
		if errors.Is(err, ErrDoesNotExist) {
			err = err.(serr.SError).CloneUnwrap()
		}
		goto end
	}
end:
	return err
}

func EmptyStateSatisfied[S ~[]RM, RM RequirementsMeeter](ctx Context, tt TokenType, rms S) (err error) {
	return forEachRequirement(ctx, rms, func(ctx Context, rm RequirementsMeeter) error {
		return rm.EmptyStateSatisfied(ctx, tt)
	})
}

func ValidateByFunc[S ~[]RM, RM RequirementsMeeter](ctx Context, rms S) (err error) {
	return forEachRequirement(ctx, rms, func(ctx Context, rm RequirementsMeeter) error {
		return rm.ValidateByFunc(ctx)
	})
}

func CheckExistence[S ~[]RM, RM RequirementsMeeter](ctx Context, rms S) (err error) {
	return forEachRequirement(ctx, rms, func(ctx Context, rm RequirementsMeeter) error {
		return rm.CheckExistence(ctx)
	})
}

func forEachRequirement[S ~[]RM, RM RequirementsMeeter](ctx Context, rms S, check func(Context, RequirementsMeeter) error) (err error) {
	for _, rm := range rms {
		err = check(ctx, rm)
		if err != nil {
			goto end
		}
	}
end:
	return err
}
