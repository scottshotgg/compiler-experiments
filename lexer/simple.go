package lexer

import (
	"fmt"
	"io/ioutil"
)

// Maybe I should make implementations that build on eachother, like StringLexer, FileLexer, PackageLexer, etc

type SimpleLexer struct {
	source  []byte
	lexemes []*Lexeme
}

// Later on there will be a NewSimpleFromFolder, NewSimpleFromPackage, NewSimpleFromJSON, NewSimpleFromSnapshot, etc

func NewSimpleFromFile(sourceFile string) (*SimpleLexer, error) {
	var contents, err = ioutil.ReadFile(sourceFile)
	if err != nil {
		return nil, err
	}

	return &SimpleLexer{
		source: contents,
	}, nil
}

func NewSimpleFromString(source string) *SimpleLexer {
	return &SimpleLexer{
		source: []byte(source),
	}
}

// TODO: this function needs to be fixed; it does not properly retain ownership of the lexemes as it returns pointers
func (l *SimpleLexer) Tokens() []*Lexeme {
	return l.lexemes
}

func (l *SimpleLexer) Source() string {
	return string(l.source)
}

// TODO: implement json marshaler

// TODO: for now this will be a simple synchronous function fully tokenizing a single file
// TODO: this should really implement scanner

func (l *SimpleLexer) Tokenize() error {
	// TODO: reduce allocations later and just store two indexes and a pointer to the data from source
	// var accumulator []byte
	var (
		accumulator string

		// TODO: start needs to be reset later
		// start int
		letter byte
	)

	// TODO: looking ahead and/or using anonymous functions may drastically simplify the logic here
	for i := 0; i < len(l.source); i++ {
		letter = l.source[i]

		// TODO: can probably use a switch or something more concise later on
		if isWhitespace(letter) {
			// Cut the token, don't track these for now
			l.lexemes = append(l.lexemes, &Lexeme{
				value: accumulator,
			})

			accumulator = ""

		} else if isEnding(letter) {
			// Cut the token and collect the ending

			l.lexemes = append(l.lexemes, &Lexeme{
				value: accumulator,
			})

			accumulator = ""

			l.lexemes = append(l.lexemes, &Lexeme{
				value: string(l.source[i]),
			})

		} else {
			// Accumulate
			// accumulator = append(accumulator, letter)
			accumulator += string(letter)
		}
	}

	if len(accumulator) > 0 {
		l.lexemes = append(l.lexemes, &Lexeme{
			value: accumulator,
		})

		accumulator = ""
	}

	return nil
}

func (l *SimpleLexer) PrintLexemes() {
	for _, t := range l.lexemes {
		fmt.Println(t.value)
	}
}

func isWhitespace(letter byte) bool {
	// TODO: just transform it to string, change this later to reduce needless allocations
	var _, ok = whitespaceLexemes[string(letter)]

	return ok
}

func isEnding(letter byte) bool {
	// TODO: make this more extensible latter
	return letter == ';'
}
