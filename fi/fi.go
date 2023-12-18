package fi

import (
	"context"

	"github.com/mikeschinkel/go-lib"
)

type ContextKey string

const Key ContextKey = "func-injection"

type FI struct{}

func GetFI[T any](ctx context.Context) T {
	fi := ctx.Value(Key)
	if fi == nil {
		lib.Panicf("Func Injector not yet set as a value for context.Context.")
	}
	return fi.(T)
}

func WrapContextFI[T any](ctx context.Context, fi T) context.Context {
	return context.WithValue(ctx, Key, fi)
}

func UpdateContextFI[T any](ctx context.Context, f func(T) T) context.Context {
	return WrapContextFI(ctx, f(GetFI[T](ctx)))
}
