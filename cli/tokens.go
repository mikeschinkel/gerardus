package cli

import (
	"fmt"
	"strings"
)

type Token string
type Tokens []Token

func (t Tokens) Args() []string {
	var skip bool
	args := make([]string, len(t))
	index := 0
	for _, token := range t[1:] {
		if skip {
			skip = false
			continue
		}
		if token[0] != '-' {
			args[index] = string(token)
			index++
			continue
		}
		if !strings.Contains(string(token), "=") {
			skip = true
			continue
		}
	}
	return args[:index]
}

func (t Tokens) Options() []string {
	var option Token

	options := make([]string, len(t))
	index := 0
	for _, token := range t[1:] {
		if option != "" {
			option = ""
			options[index] = fmt.Sprintf("%s=%s", option, token)
			index++
			continue
		}
		if token[0] != '-' {
			continue
		}
		if strings.Contains(string(token), "=") {
			options[index] = string(token)
			index++
			continue
		}
		option = token
	}
	return options[:index]
}
