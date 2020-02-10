package lexer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/pkg/errors"

	token "github.com/scottshotgg/express-token"
)

// Lexer holds all the needed variables to appropriately lex
type (
	SimpleLexer struct {
		source      []rune
		Accumulator string
		r           io.RuneScanner
		tokens      []token.Token
	}
)

// Lexemes are the specific symbols the lexer needs to recognize
var (
	Lexemes = []string{
		// "var",
		// "int",
		// "float",
		// "string",
		// "bool",
		// "char",
		// "object",

		":",
		"=",
		"+",
		"-",
		"*",
		"/",
		"(",
		")",
		"{",
		"}",
		"[",
		"]",
		"\"",
		"'",
		";",
		",",
		"#",
		"!",
		"<",
		">",
		"@",
		"\\",
		// "„",
		" ",
		"\n",
		"\t",

		// "select",
		// "SELECT",
		// "FROM",
		// "WHERE",
	}
)

func New(r io.RuneScanner) *SimpleLexer {
	return &SimpleLexer{
		r: r,
	}
}

// NewFromBytes is a simple helper function to provide a lexer from a list of bytes
func NewFromBytes(source []byte) *SimpleLexer {
	return New(bufio.NewReader(bytes.NewReader(source)))
}

// NewFromBytes is a simple helper function to provide a lexer from a string
func NewFromString(source string) *SimpleLexer {
	return New(bufio.NewReader(bytes.NewBufferString(source)))
}

func NewFromFile(f *os.File) *SimpleLexer {
	return New(bufio.NewReader(f))
}

func NewFromFolder(path string)

// NewFromPath is a simple helper function to provide a lexer from a provided path
func NewFromPath(path string, prefetchContents bool) (*SimpleLexer, error) {
	// TODO: this needs to determine if it is a folder or file

	var (
		l         *SimpleLexer
		stat, err = os.Stat(path)
	)

	if err != nil {
		return nil, err
	}

	if prefetchContents {
		if stat.IsDir() {
			return nil, errors.New("not implemented")
			// return NewFromFolder(path)
		}

		var data, err = ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		l = NewFromBytes(data)

	} else {
		if stat.IsDir() {
			return nil, errors.New("not implemented")
			// return NewFromFolder(path)
		}

		// I would imagine this would greatly reduce memory usage but I also think it will also significantly hamper performance
		var f, err = os.Open(path)
		if err != nil {
			return nil, err
		}

		l = NewFromFile(f)
	}

	return l, nil
}

// Source prints the source
func (l *SimpleLexer) Source() string {
	return string(l.source)
}

// Tokens return the collected tokens from the lexer
func (l *SimpleLexer) Tokens() []token.Token {
	return l.tokens
}

// Print displays all collected tokens
func (l *SimpleLexer) Print() {
	for _, t := range l.tokens {
		fmt.Println("token:", t)
	}
}

// TODO: this is a temporary method to plug the old lexer architecture into the new one
func (l *SimpleLexer) readAllRunes() error {
	for {
		var r, _, err = l.r.ReadRune()
		switch err {
		case nil:
			l.source = append(l.source, r)

			continue

		case io.EOF:
			break

		default:
			return err
		}
	}
}

// Tokenize is the primary function used to lex the source into tokens
func (l *SimpleLexer) Tokenize() error {
	// TODO: convert this to use a Reader instead
	// For now just read the entire thing into source to start
	var err = l.readAllRunes()
	if err != nil {
		return err
	}

	for index := 0; index < len(l.source); index++ {
		var char = string(l.source[index])

		// Else see if it's recognized lexeme
		var lexemeToken, ok = token.TokenMap[char]

		// // Only the operators are allowed to be without spaces after them; this may change, kinda hate no spaces between the symbols
		// // Also enclosers (rbrace, lbrace, etc) are allowed as well. End tokens (; and ,) as well
		// // Make something in the token library for this, a specific struct field

		// // If white space IS required after the token ...
		// if ok && lexemeToken.Type != token.Whitespace && !lexemeToken.WSNotRequired {
		// 	// If the current token is not allowed to not have whitespace after it, the next character has to be some sort of whitespace (space, newline, tab)
		// 	// next character is not a white space and we require it then there is an error
		// 	if index+1 < len(l.source) && !unicode.IsSpace(l.source[index+1]) {
		// 		// Not sure if using the unicode library is the right way to go ...
		// 		// return nil, errors.Errorf("Expected white space after token (%s), found: %s", string(l.source[index]), string(l.source[index+1]))
		// 		// It should not be a recognized token, add it to the accumulator and move on as if it was as normal char

		// 		// Test if the next character is a recognized token?
		// 		_, ok = token.TokenMap[string(l.source[index+1])]
		// 		if ok {
		// 			return nil, errors.Errorf("Expected white space after token (%s), found: %s", string(l.source[index]), string(l.source[index+1]))
		// 		}

		// 		l.Accumulator += char
		// 		continue
		// 	}
		// }

		// If it is not a recognized lexeme, add it to the accumulator and move on
		if !ok {
			l.Accumulator += char
			continue
		}

		// Filter out the comments
		switch lexemeToken.Value.Type {
		case "div":
			index++
			if index < len(l.source)-1 {
				switch l.source[index] {
				case '/':
					for {
						index++
						if index == len(l.source) || l.source[index] == '\n' {
							break
						}
					}

				case '*':
					for {
						index++
						if index == len(l.source) || (l.source[index] == '*' && l.source[index+1] == '/') {
							index++
							break
						}
					}

				default:
					l.tokens = append(l.tokens, token.TokenMap[char])
				}
			}

			continue

		// Use the lexer to parse strings
		case "squote":
			fallthrough

		case "dquote":
			// If the accumulator is not empty, check it before parsing the string
			if l.Accumulator != "" {
				ts, err := l.LexLiteral()
				if err != nil {
					return err
				}

				l.tokens = append(l.tokens, ts)
				l.Accumulator = ""
			}

			stringLiteral := ""

			index++
			for string(l.source[index]) != lexemeToken.Value.String {
				// If there is an escaping backslash in the string then just increment over
				// it so that the next accumulate and increment will pickup the next char naturally
				if string(l.source[index]) == "\\" {
					index++
				}

				stringLiteral += string(l.source[index])

				index++
			}

			// Don't allow strings to use single quotes like JS
			stringType := token.StringType
			if lexemeToken.Value.Type == "squote" {
				if len(stringLiteral) > 1 {
					return errors.Errorf("Too many values in character literal declaration: %s", stringLiteral)
				}

				stringType = token.CharType
			}

			l.tokens = append(l.tokens, token.Token{
				ID:   0,
				Type: token.Literal,
				Value: token.Value{
					Type:   stringType,
					True:   stringLiteral,
					String: stringLiteral,
				},
			})

			continue

		case "period":
			// For now just accumulate the period and evaluate it later during parsing
			l.Accumulator += char
			continue
		}

		// If the accumulator is not empty, check it
		if l.Accumulator != "" {
			ts, err := l.LexLiteral()
			if err != nil {
				return err
			}

			l.tokens = append(l.tokens, ts)
		}

		// Append the current token and reset the accumulator
		l.tokens = append(l.tokens, lexemeToken)
		l.Accumulator = ""
	}

	// If the accumulator is not empty, check it
	if l.Accumulator != "" {
		ts, err := l.LexLiteral()
		if err != nil {
			return err
		}

		l.tokens = append(l.tokens, ts)
	}

	return nil
}

// LexLiteral is used for determining whether something is a ident or literal
// If it is a literal is it a string, char, int, float, or bool
func (l *SimpleLexer) LexLiteral() (token.Token, error) {
	// Make a token and set the default value to bool; this is just because its the
	// first case in the switch and everything below sets it, so it makes the code a bit
	// cleaner
	// We COULD do this with tokens in the tokenMap for true and false
	var t = token.Token{
		Type: token.Literal,
		Value: token.Value{
			True:   false,
			Type:   token.BoolType,
			String: l.Accumulator,
		},
	}

	switch l.Accumulator {
	// Default value is false, we only need to catch the case to keep it out of the default
	case "false":

	// Check if its true
	case "true":
		t.Value.True = true

	// Else move on and figure out what kind of number it is (or an ident)
	default:
		// Figure out from the two starting characters
		var base = 10
		if len(l.Accumulator) > 2 {
			switch l.Accumulator[:2] {
			// Binary
			case "0b":
				base = 2

			// Octal
			case "0o":
				base = 8

			// Hex
			case "0x":
				base = 16
			}
		}

		// If the base is not 10 anymore, shave off the 0b, 0o, or 0x
		if base != 10 {
			l.Accumulator = l.Accumulator[2:]
		}

		// Attempt to parse an int from the accumulator
		var value, err = strconv.ParseInt(l.Accumulator, base, 64)

		// TODO: Convert the int64 to an int for now
		// I'll switch this when I'm ready to deal with different bit sizes
		t.Value.True = int(value)
		t.Value.Type = token.IntType

		// TODO: need to make something for scientific notation with carrots and e
		// If it errors, check to see if it is an float
		if err != nil {
			// Attempt to parse a float from the accumulator
			t.Value.True, err = strconv.ParseFloat(l.Accumulator, 64)
			t.Value.Type = token.FloatType
			if err != nil {
				// If it's not a float, check whether it is a keyword
				keyword, ok := token.TokenMap[l.Accumulator]
				if ok {
					t = keyword
				} else {
					// If it is not a keyword or a parse-able number, assume that it is an ident (for now)
					t.Type = token.Ident
					t.Value = token.Value{
						String: l.Accumulator,
					}
				}
			}
		}
	}

	return t, nil
}
