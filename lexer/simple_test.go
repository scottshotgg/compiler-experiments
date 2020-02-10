package lexer_test

import (
	"fmt"
	"testing"

	"github.com/scottshotgg/compiler-experiments/lexer"
)

func TestTokenize(t *testing.T) {

	var (
		test = "hi my name is scott"

		l = lexer.NewFromString(test)
	)

	fmt.Println("hi")
	err := l.Tokenize()

	if err != nil {
		t.Fatal(err)
	}

	l.Print()
}
