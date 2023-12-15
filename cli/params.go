package cli

type Params struct {
	AppName   string
	EnvPrefix string
	OSArgs    []string
	tokens    Tokens
}

func (p Params) Tokens() Tokens {
	if p.tokens != nil {
		goto end
	}
	p.tokens = make(Tokens, len(p.OSArgs))
	for i, a := range p.OSArgs {
		p.tokens[i] = Token(a)
	}
end:
	return p.tokens
}

func (p Params) Args() []string {
	return p.Tokens().Args()
}

func (p Params) Options() []string {
	return p.Tokens().Options()
}
