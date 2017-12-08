package lexer

import tk "github.com/ryym/monkey/01/token"

type Lexer struct {
	input        string
	position     int // current position
	readPosition int // curreint reading position (next char)
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// To keep things simple, only support ASCIIs.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() tk.Token {
	tok := tk.Token{}

	switch l.ch {
	case '=':
		tok = newToken(tk.ASSIGN, l.ch)
	case ';':
		tok = newToken(tk.SEMICOLON, l.ch)
	case '(':
		tok = newToken(tk.LPAREN, l.ch)
	case ')':
		tok = newToken(tk.RPAREN, l.ch)
	case ',':
		tok = newToken(tk.COMMA, l.ch)
	case '+':
		tok = newToken(tk.PLUS, l.ch)
	case '{':
		tok = newToken(tk.LBRACE, l.ch)
	case '}':
		tok = newToken(tk.RBRACE, l.ch)
	case 0:
		tok.Type = tk.EOF
		tok.Literal = ""
	}

	l.readChar()
	return tok
}

func newToken(tp tk.TokenType, ch byte) tk.Token {
	return tk.Token{tp, string(ch)}
}
