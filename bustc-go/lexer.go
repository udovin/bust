package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Position struct {
	// Line contains line number.
	Line int
	// Column contains column number.
	Column int
}

type Lexer struct {
	reader   *bufio.Reader
	token    Token
	position Position
	err      error
}

func (l *Lexer) Next() bool {
	l.token.Position = l.position
	r, ok := l.readRune()
	for unicode.IsSpace(r) && ok {
		l.token.Position = l.position
		r, ok = l.readRune()
	}
	if !ok {
		if l.err == io.EOF {
			l.err = nil
		}
		return false
	}
	b := strings.Builder{}
	b.WriteRune(r)
	switch {
	case r == '\'':
		return l.nextRuneToken(&b)
	case r == '"':
		return l.nextStringToken(&b)
	case r == '{':
		l.token.Kind = OpenBraceToken
	case r == '}':
		l.token.Kind = CloseBraceToken
	case r == '(':
		l.token.Kind = OpenParenthesisToken
	case r == ')':
		l.token.Kind = CloseParenthesisToken
	case r == '[':
		l.token.Kind = OpenBracketToken
	case r == ']':
		l.token.Kind = CloseBracketToken
	case r == ',':
		l.token.Kind = CommaToken
	case r == '+':
		l.token.Kind = AddToken
	case r == '-':
		l.token.Kind = SubToken
	case r == '*':
		l.token.Kind = MulToken
	case r == '<':
		l.token.Kind = LessToken
	case r == '>':
		l.token.Kind = GreaterToken
	case r == '/':
		r2, unread, ok := l.tryReadRune()
		if ok {
			if r2 == '/' {
				b.WriteRune(r2)
				return l.nextLineCommentToken(&b)
			}
			unread()
		}
		l.token.Kind = DivToken
	case r == ':':
		r2, ok := l.readRune()
		if ok {
			switch r2 {
			case '=':
				b.WriteRune(r2)
				l.token.Kind = InitToken
				l.token.Text = b.String()
				return true
			case ':':
				b.WriteRune(r2)
				l.token.Kind = ScopeToken
				l.token.Text = b.String()
				return true
			}
		}
		l.err = fmt.Errorf("unknown token: %c", r)
		l.token.Kind = InvalidToken
		l.token.Text = b.String()
		return false
	case r == '&':
		r2, unread, ok := l.tryReadRune()
		if ok {
			if r2 == '&' {
				b.WriteRune(r2)
				l.token.Kind = BoolAndToken
			} else {
				unread()
				l.token.Kind = AndToken
			}
		} else {
			l.token.Kind = AndToken
		}
	case r == '=':
		r2, unread, ok := l.tryReadRune()
		if ok {
			if r2 == '=' {
				b.WriteRune(r2)
				l.token.Kind = EqualToken
			} else {
				unread()
				l.token.Kind = AssignToken
			}
		} else {
			l.token.Kind = AssignToken
		}
	case isDigit(r):
		return l.nextNumberToken(&b)
	case isLetter(r):
		return l.nextIdentifierToken(&b)
	default:
		l.err = fmt.Errorf("unknown token: %c", r)
		l.token.Kind = InvalidToken
		l.token.Text = b.String()
		return false
	}
	l.token.Text = b.String()
	return true
}

func (l *Lexer) Err() error {
	return l.err
}

func (l *Lexer) Token() Token {
	return l.token
}

func (l *Lexer) Position() Position {
	return l.position
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		reader: bufio.NewReader(reader),
	}
}

func (l *Lexer) readRune() (rune, bool) {
	r, _, err := l.reader.ReadRune()
	if err != nil {
		l.err = err
		return 0, false
	}
	l.position.Column++
	if isNextLine(r) {
		l.position.Line++
		l.position.Column = 0
	}
	return r, true
}

func (l *Lexer) tryReadRune() (rune, func(), bool) {
	p := l.position
	r, ok := l.readRune()
	if !ok {
		return 0, nil, false
	}
	cancel := func() {
		if err := l.reader.UnreadRune(); err != nil {
			panic(err)
		}
		l.position = p
	}
	return r, cancel, true
}

func (l *Lexer) nextIdentifierToken(b *strings.Builder) bool {
	r, unread, ok := l.tryReadRune()
	for (isLetter(r) || isDigit(r)) && ok {
		b.WriteRune(r)
		r, unread, ok = l.tryReadRune()
	}
	if !ok {
		if l.err != io.EOF {
			return false
		}
	} else {
		unread()
	}
	l.token.Text = b.String()
	if kind, ok := keywords[l.token.Text]; ok {
		l.token.Kind = kind
	} else {
		l.token.Kind = IdentifierToken
	}
	return true
}

func (l *Lexer) nextNumberToken(b *strings.Builder) bool {
	r, unread, ok := l.tryReadRune()
	for isDigit(r) && ok {
		b.WriteRune(r)
		r, unread, ok = l.tryReadRune()
	}
	if !ok {
		if l.err != io.EOF {
			return false
		}
	} else {
		unread()
	}
	l.token.Kind = NumberToken
	l.token.Text = b.String()
	return true
}

func (l *Lexer) readEscapedRune(b *strings.Builder) (rune, bool, bool) {
	r, ok := l.readRune()
	if !ok {
		return r, false, ok
	}
	if r == '\\' {
		b.WriteRune(r)
		r, ok = l.readRune()
		return r, true, ok
	}
	return r, false, ok
}

func (l *Lexer) nextRuneToken(b *strings.Builder) bool {
	if r, escaped, ok := l.readEscapedRune(b); !ok {
		return false
	} else {
		if !escaped && r == '\'' {
			l.err = fmt.Errorf("unexpected character in rune: %c", r)
			return false
		}
		b.WriteRune(r)
	}
	if r, ok := l.readRune(); !ok {
		return false
	} else {
		if r != '\'' {
			l.err = fmt.Errorf("rune should ends with: %c", '\'')
			return false
		}
		b.WriteRune(r)
	}
	l.token.Kind = RuneToken
	l.token.Text = b.String()
	return true
}

func (l *Lexer) nextStringToken(b *strings.Builder) bool {
	r, escaped, ok := l.readEscapedRune(b)
	for (escaped || r != '"') && ok {
		b.WriteRune(r)
		r, escaped, ok = l.readEscapedRune(b)
	}
	if !ok {
		return false
	}
	b.WriteRune(r)
	l.token.Kind = StringToken
	l.token.Text = b.String()
	return true
}

func (l *Lexer) nextLineCommentToken(b *strings.Builder) bool {
	r, ok := l.readRune()
	for !isNextLine(r) && ok {
		b.WriteRune(r)
		r, ok = l.readRune()
	}
	l.token.Kind = LineCommentToken
	l.token.Text = b.String()
	return true
}

func isNextLine(r rune) bool {
	return r == '\n'
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
