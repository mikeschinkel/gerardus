package fi

//
//import (
//	"fmt"
//	"reflect"
//	"runtime"
//	"strings"
//)
//
//// FuncMapKey returns a function signature in slightly-modified form from that
//// which would be syntactically correct in Go. It drops "func ", replaces parens
//// with vertical bars so there is no need for a trailing parens for return
//// values, replaces three dots (...) for variadic with a single dot, and prefixes
//// with the funcs package, e.g.
////
////	func Foo(int,bool)(string,error) => main.Foo|int,bool|string,error
////	func (b Bar) String() string     => main.Bar.String||string
////	func (b *Bar) PtrMethod() string => main.(*Bar).PtrMethod|*main.Bar|string
////	func VariadicExample(...int)     => main.VariadicExample|.int|
////	func Baz(f *Bar, n int)          => main.Baz|*main.Bar,int|
////
//// Developed with the help of ChatGPT.
////
////	See: https://chat.openai.com/share/ed10e21b-1e3c-4aef-9d25-c1aaf8a1961a
//func FuncMapKey(f any) (s string) {
//	var sb strings.Builder
//
//	rt := reflect.TypeOf(f)
//	if rt.Kind() != reflect.Func {
//		goto end
//	}
//	sb.WriteString(FullFuncnameOf(f))
//	sb.WriteByte('|')
//	for i := 0; i < rt.NumIn(); i++ {
//		if i > 0 {
//			sb.WriteString(",")
//		}
//		pt := rt.In(i)
//		if rt.IsVariadic() && i == rt.NumIn()-1 {
//			sb.WriteString("." + pt.Elem().Name())
//		} else {
//			sb.WriteString(fmt.Sprintf("%v", pt))
//		}
//	}
//	sb.WriteByte('|')
//
//	for i := 0; i < rt.NumOut(); i++ {
//		if i > 0 {
//			sb.WriteString(",")
//		}
//		ot := rt.Out(i)
//		sb.WriteString(ot.Name())
//	}
//
//	s = sb.String()
//end:
//	return s
//}
//
//// FullFuncnameOf returns a function name with its package prefixed given a
//// func value e.g. if we have a func `Foo()` in `bar` package then the func
//// reference would be `bar.Foo`, unless the current package was `bar` and then it
//// would just be `Foo`. If we have a method `bar.Baz()` where `bar` in an
//// instance of the type and then the func reference would be just `bar.Baz`
//func FullFuncnameOf(v any) (name string) {
//	var p uintptr
//	var fn *runtime.Func
//	rv := reflect.ValueOf(v)
//	if rv.Kind() != reflect.Func {
//		goto end
//	}
//	p = rv.Pointer()
//	if p == 0 {
//		name = "<anonymous>"
//		goto end
//	}
//	fn = runtime.FuncForPC(p)
//	if fn == nil {
//		name = "<unknown>"
//	}
//	// Go uses "-fm" for internally generated funcs for a closure with receiver
//	// capture, per ChatGPT. See:
//	// https://chat.openai.com/share/ed10e21b-1e3c-4aef-9d25-c1aaf8a1961a
//	// Search for "-fm"
//	name = strings.TrimSuffix(fn.Name(), "-fm")
//end:
//	return name
//}
//
////func main() {
////	fmt.Println(FuncMapKey(Foo))
////	fmt.Println(FuncMapKey(Bar{}.String))
////	fmt.Println(FuncMapKey((*Bar).PointerMethod))
////	fmt.Println(FuncMapKey(VariadicExample))
////	fmt.Println(FuncMapKey(Baz))
////}
////
////func Foo(n int, b bool) (s string, err error) {
////	return s, err
////}
////
////type Bar struct{}
////
////func (b Bar) String() string {
////	return ""
////}
////
////func (b *Bar) PointerMethod() string {
////	return ""
////}
////
////func VariadicExample(args ...int) {
////}
////
////func Baz(f *Bar, n int) {}
////
