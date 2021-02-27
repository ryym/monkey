package lexer

import tk "github.com/ryym/monkey/token"

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

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0 // EOF
	} else {
		return l.input[l.readPosition]
	}
}

// To keep things simple, only support ASCIIs.
func (l *Lexer) readChar() {
	l.ch = l.peekChar()
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

func (l *Lexer) readTwoChars() string {
	ch := l.ch
	l.readChar()
	return string(ch) + string(l.ch)
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
		if l.peekChar() == '=' {
			lit := l.readTwoChars()
			tok.Type = tk.EQ
			tok.Literal = lit
		} else {
			tok = newToken(tk.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(tk.PLUS, l.ch)
	case '-':
		tok = newToken(tk.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			lit := l.readTwoChars()
			tok.Type = tk.NOT_EQ
			tok.Literal = lit
		} else {
			tok = newToken(tk.BANG, l.ch)
		}
	case '*':
		tok = newToken(tk.ASTERISK, l.ch)
	case '/':
		tok = newToken(tk.SLASH, l.ch)
	case '<':
		tok = newToken(tk.LT, l.ch)
	case '>':
		tok = newToken(tk.GT, l.ch)
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
	return tk.Token{Type: tp, Literal: string(ch)}
}

func isLetter(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
