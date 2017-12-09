package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ryym/monkey/01/lexer"
	"github.com/ryym/monkey/01/token"
)

const PROMPT = ">> "

// XXX: `out` is unused.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		if !scanner.Scan() {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
