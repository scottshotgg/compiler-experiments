package lexer

// Lexeme may need an interface as well

type (
	Lexeme struct {
		// TODO: use bytes here to reduce allocations
		value string
		// start int
		// end   int

		// TODO: this guy should have helper functions like IsWhitespace, etc
	}
)
