package lexer

import (
	token "github.com/scottshotgg/express-token"
)

type (
	Lexer interface {
		// Tokens will return the collected tokens
		Tokens() []token.Token

		// Source returns the original source text fed to the lexer
		Source() string

		// Tokenize instructs the lexer to tokenize the source
		Tokenize() error

		// Collect() error
	}
)
