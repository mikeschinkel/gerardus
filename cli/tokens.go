package cli

import (
	"fmt"
	"strings"
)

type Token string
type Tokens []Token

// Join joins a slice of Tokens into a strings.
// It is a copy of strings.Join() with small changes for tokens.
func (tt Tokens) Join(sep string) string {
	const maxInt = int(^uint(0) >> 1)

	switch len(tt) {
	case 0:
		return ""
	case 1:
		return string(tt[0])
	}

	var n int
	if len(sep) > 0 {
		if len(sep) >= maxInt/(len(tt)-1) {
			panic("strings: Join output length overflow")
		}
		n += len(sep) * (len(tt) - 1)
	}
	for _, elem := range tt {
		if len(elem) > maxInt-n {
			panic("strings: Join output length overflow")
		}
		n += len(elem)
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(string(tt[0]))
	for _, s := range tt[1:] {
		b.WriteString(sep)
		b.WriteString(string(s))
	}
	return b.String()
}

// Args returns the tokens from os.Args that are not a `-flag` ot its value
func (tt Tokens) Args() (_ Tokens, err error) {
	var skipNext bool
	args := make(Tokens, 0)
	index := 0
	for _, token := range tt[1:] {
		if skipNext {
			// Skip the value for the prior `-flag`
			skipNext = false
			continue
		}
		if token[0] != '-' {
			// It's not a flag, collect it
			args = append(args, token)
			index++
			continue
		}
		if len(args) > 0 {
			// Flags specified after args, return error
			err = ErrOptionAfterArgs.Args(
				"option", token,
				"args", args.Join(" "),
			)
			goto end
		}
		if !strings.Contains(string(token), "=") {
			// It's in format of `-flag value`, not `-flag=value` so ignore the next token as
			// it is a value and not an arg (aka part of the command.)
			skipNext = true
			continue
		}
	}
	tt = args[:index]
end:
	return tt, err
}

func (tt Tokens) Count() int {
	return len(tt) - 1
}

func (tt Tokens) StringSlice() (ss []string) {
	ss = make([]string, len(tt))
	for i, t := range tt {
		ss[i] = string(t)
	}
	return ss
}

func (tt Tokens) Options() Tokens {
	var option Token

	options := make(Tokens, len(tt))
	index := 0
	for _, token := range tt[1:] {
		if option != "" {
			option = ""
			options[index] = Token(fmt.Sprintf("%s=%s", option, token))
			index++
			continue
		}
		if token[0] != '-' {
			continue
		}
		if strings.Contains(string(token), "=") {
			options[index] = token
			index++
			continue
		}
		option = token
	}
	return options[:index]
}
