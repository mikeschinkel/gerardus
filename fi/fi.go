package fi

import (
	"context"

	"github.com/mikeschinkel/go-lib"
)

type ContextKey string

const Key ContextKey = "func-injection"

type FI struct {
	valid bool
}

func (fi FI) IsValid() bool {
	return fi.valid
}

func (fi FI) SetValid() {
	fi.valid = true
}

func GetFI[T any](ctx context.Context) T {
	fi := ctx.Value(Key)
	if fi == nil {
		lib.Panicf("Func Injector not yet set as a value for context.Context.")
	}
	return fi.(T)
}

func setValid(fi any) any {
	tmp := fi.(interface{ SetValid() })
	tmp.SetValid()
	return tmp
}

func WrapContextFI[T any](ctx context.Context, fi T) context.Context {
	return context.WithValue(ctx, Key, setValid(fi))
}

func UpdateContextFI[T any](ctx context.Context, f func(T) T) context.Context {
	return WrapContextFI(ctx, f(GetFI[T](ctx)))
}
