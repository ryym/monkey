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

func (l *Lexer) readIdentifier() string {
	from := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[from:l.position]
}

// Ignore negative numbers, floats, etc.
func (l *Lexer) readNumber() string {
	from := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[from:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) NextToken() tk.Token {
	tok := tk.Token{}

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(tk.ASSIGN, l.ch)
	case '+':
		tok = newToken(tk.PLUS, l.ch)
	case '(':
		tok = newToken(tk.LPAREN, l.ch)
	case ')':
		tok = newToken(tk.RPAREN, l.ch)
	case ',':
		tok = newToken(tk.COMMA, l.ch)
	case ';':
		tok = newToken(tk.SEMICOLON, l.ch)
	case '{':
		tok = newToken(tk.LBRACE, l.ch)
	case '}':
		tok = newToken(tk.RBRACE, l.ch)
	case 0:
		tok.Type = tk.EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = tk.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = tk.INT
			return tok
		} else {
			tok = newToken(tk.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func newToken(tp tk.TokenType, ch byte) tk.Token {
	return tk.Token{tp, string(ch)}
}

func isLetter(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
