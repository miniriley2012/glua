package glua

import (
	"bufio"
	"errors"
	"io"
	"unicode"

	"glua/5.5/internal/tables"
	"glua/5.5/token"
)

func New(rd io.Reader) *Scanner {
	if rd, satisfied := rd.(io.RuneScanner); satisfied {
		return &Scanner{rd: rd}
	}
	return &Scanner{rd: bufio.NewReader(rd)}
}

type Scanner struct {
	rd  io.RuneScanner
	p   int
	n   int
	err error

	Token token.Token
	Start int
	End   int
}

func (scn *Scanner) Err() error {
	return scn.err
}

func (scn *Scanner) Scan() bool {
	r, err := scn.nextRune()
	if err != nil {
		return scn.error(err)
	}

	for unicode.In(r, unicode.Pattern_White_Space) {
		r, err = scn.nextRune()
		if err != nil {
			return scn.error(err)
		}
	}

	err = scn.unreadRune()
	if err != nil {
		return scn.error(err)
	}

	scn.mark()

	if unicode.In(r, tables.XID_Start) {
		r, err = scn.nextRune()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return scn.error(err)
			}
		}

		var runes []rune
		for unicode.In(r, tables.XID_Continue) {
			runes = append(runes, r)
			r, err = scn.nextRune()
			if err != nil {
				if !errors.Is(err, io.EOF) {
					return scn.error(err)
				}
			}
		}

		if err == nil {
			err = scn.unreadRune()
			if err != nil {
				return scn.error(err)
			}
		}

		tok, found := keywords[string(runes)]
		if !found {
			tok = token.Name
		}

		return scn.push(tok)
	}

	r, err = scn.nextRune()
	if err != nil {
		return scn.error(err)
	}

	switch r {
	case '+':
		return scn.push(token.Plus)
	case '-':
		return scn.push(token.Minus)
	case '*':
		return scn.push(token.Asterisk)
	case '/':
		r, err = scn.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return scn.push(token.Solidus)
			} else {
				return scn.error(err)
			}
		}
		if r == '/' {
			return scn.push(token.DoubleSolidus)
		}
		return scn.push2(token.Solidus)
	case '%':
		return scn.push(token.Percent)
	case '^':
		return scn.push(token.Caret)
	case '#':
		return scn.push(token.Hash)
	case '&':
		return scn.push(token.BitAnd)
	case '~':
		r, err = scn.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return scn.push(token.Tilde)
			} else {
				return scn.error(err)
			}
		}
		if r == '=' {
			return scn.push(token.NotEqual)
		}
		return scn.push2(token.Tilde)
	case '|':
		return scn.push(token.BitOr)
	case '<':
		r, err = scn.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return scn.push(token.LessThan)
			} else {
				return scn.error(err)
			}
		}
		switch r {
		case '=':
			return scn.push(token.LessThanEqual)
		case '<':
			return scn.push(token.DoubleLessThan)
		}
		return scn.push2(token.LessThan)
	case '>':
		r, err = scn.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return scn.push(token.GreaterThan)
			} else {
				return scn.error(err)
			}
		}
		switch r {
		case '=':
			return scn.push(token.GreaterThanEqual)
		case '>':
			return scn.push(token.DoubleGreaterThan)
		}
		return scn.push2(token.GreaterThan)
	case '=':
		r, err = scn.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return scn.push(token.Equal)
			} else {
				return scn.error(err)
			}
		}
		if r == '=' {
			return scn.push(token.DoubleEqual)
		}
		return scn.push2(token.Equal)
	case '(':
		return scn.push(token.LeftParen)
	case ')':
		return scn.push(token.RightParen)
	case '{':
		return scn.push(token.LeftBrace)
	case '}':
		return scn.push(token.RightBrace)
	case '[':
		return scn.push(token.LeftBracket)
	case ']':
		return scn.push(token.RightBracket)
	case ':':
		r, err = scn.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return scn.push(token.Colon)
			} else {
				return scn.error(err)
			}
		}
		if r == ':' {
			return scn.push(token.DoubleColon)
		}
		return scn.push2(token.Colon)
	case ';':
		return scn.push(token.Semicolon)
	case ',':
		return scn.push(token.Comma)
	case '.':
		r, err = scn.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return scn.push(token.Dot)
			} else {
				return scn.error(err)
			}
		}
		if r == '.' {
			r, err = scn.nextRune()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return scn.push(token.DoubleDot)
				} else {
					return scn.error(err)
				}
			}
			if r == '.' {
				return scn.push(token.TripleDot)
			}
			return scn.push2(token.DoubleDot)
		}
		return scn.push2(token.Dot)
	}

	return scn.error(errors.ErrUnsupported)
}

func (scn *Scanner) error(err error) bool {
	scn.err = err
	return false
}

func (scn *Scanner) nextRune() (rune, error) {
	r, sz, err := scn.rd.ReadRune()
	if err != nil {
		return r, err
	}
	scn.p = scn.n
	scn.n += sz
	return r, nil
}

func (scn *Scanner) unreadRune() error {
	err := scn.rd.UnreadRune()
	if err != nil {
		return err
	}
	scn.n = scn.p
	return nil
}

func (scn *Scanner) mark() {
	scn.Start = scn.n
}

func (scn *Scanner) push(tok token.Token) bool {
	scn.End = scn.n
	scn.Token = tok
	return true
}

func (scn *Scanner) push2(tok token.Token) bool {
	scn.error(scn.unreadRune())
	return scn.push(tok)
}

var keywords = map[string]token.Token{
	"and":      token.And,
	"break":    token.Break,
	"do":       token.Do,
	"else":     token.Else,
	"elseif":   token.Elseif,
	"end":      token.End,
	"false":    token.False,
	"for":      token.For,
	"function": token.Function,
	"global":   token.Global,
	"goto":     token.Goto,
	"if":       token.If,
	"in":       token.In,
	"local":    token.Local,
	"nil":      token.Nil,
	"not":      token.Not,
	"or":       token.Or,
	"repeat":   token.Repeat,
	"return":   token.Return,
	"then":     token.Then,
	"true":     token.True,
	"until":    token.Until,
	"while":    token.While,
}
