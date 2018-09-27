package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ryym/monkey/lexer"
	"github.com/ryym/monkey/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		if !scanner.Scan() {
			return
		}
		line := scanner.Text()

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}

}

func printParseErrors(out io.Writer, errs []string) {
	io.WriteString(out, "ERROR\n")
	for _, msg := range errs {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
