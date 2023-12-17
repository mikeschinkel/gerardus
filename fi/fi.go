package fi

import (
	"context"

	"github.com/mikeschinkel/go-lib"
)

type ContextKey string

const Key ContextKey = "func-injection"

type FI struct{}

func GetFI(ctx context.Context) any {
	fi := ctx.Value(Key)
	if fi == nil {
		lib.Panicf("Func Injector not yet set as a value for context.Context.")
	}
	return fi
}

func WrapContext(ctx context.Context, fi any) context.Context {
	return context.WithValue(ctx, Key, fi)
}

//func (fi *FI) Call(f any, args ...any) (results []any) {
//	rv := reflect.ValueOf(f)
//	in := make([]reflect.Value, len(args))
//	for i, arg := range args {
//		in[i] = reflect.ValueOf(arg)
//	}
//	out := rv.Call(in)
//	results = make([]any, len(out))
//	for i, result := range out {
//		results[i] = result.Interface()
//	}
//	return results
//}
