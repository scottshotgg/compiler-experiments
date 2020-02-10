package lexer

type Lexer interface {
	// Tokens will return the collected tokens
	Tokens() []*Lexeme

	// Source returns the original source text fed to the lexer
	Source() string

	// Tokenize instructs the lexer to tokenize the source
	Tokenize() error

	// Collect() error
}
