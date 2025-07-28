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

		for unicode.In(r, tables.XID_Continue) {
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

		return scn.push(token.Identifier)
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
