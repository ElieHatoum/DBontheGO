package dbonthego

import (
	"errors"
	"strconv"
	"strings"
)

type Parser struct {
	buf string
	pos int
}

func NewParser(s string) Parser {
	return Parser{buf: s, pos: 0}
}

func isSpace(ch byte) bool {
	switch ch {
	case '\t', '\n', '\v', '\f', '\r', ' ':
		return true
	}
	return false
}
func isAlpha(ch byte) bool {
	return 'a' <= (ch|32) && (ch|32) <= 'z'
}
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
func isNameStart(ch byte) bool {
	return isAlpha(ch) || ch == '_'
}
func isNameContinue(ch byte) bool {
	return isAlpha(ch) || isDigit(ch) || ch == '_'
}
func isSeparator(ch byte) bool {
	return ch < 128 && !isNameContinue(ch)
}

func (p *Parser) skipSpaces() {
	for p.pos < len(p.buf) && isSpace(p.buf[p.pos]) {
		p.pos += 1
	}
}

func (p *Parser) tryKeyword(kw string) bool {
	p.skipSpaces()
	if !(p.pos + len(kw) <= len(p.buf) && strings.EqualFold(p.buf[p.pos:p.pos+len(kw)], kw)) {
		return false
	}

	if p.pos+len(kw) < len(p.buf) && !isSeparator(p.buf[p.pos+len(kw)]) {
		return false
	} 
	p.pos += len(kw)
	return true
}

func (p *Parser) tryName() (string, bool) {
	p.skipSpaces()
	start, curr := p.pos, p.pos
	if !(curr < len(p.buf) && isNameStart(p.buf[curr])) {
		return "", false
	}

	curr++
	for curr < len(p.buf) && isNameContinue(p.buf[curr]) {
		curr++
	}
	p.pos = curr
	return p.buf[start:curr], true
}

func (p *Parser) parseValue(out *Cell) error {
	p.skipSpaces()
	if p.pos >= len(p.buf) {
		return errors.New("expect value")
	}
	ch := p.buf[p.pos]
	if ch == '"' || ch == '\'' {
		return p.parseString(out)
	} else if isDigit(ch) || ch == '-' || ch == '+' {
		return p.parseInt(out)
	} else {
		return errors.New("expect value")
	}
}

func (p *Parser) parseString(out *Cell) error {
	quote := p.buf[p.pos]
	curr := p.pos + 1
	for curr < len(p.buf) {
		ch := p.buf[curr]
		switch ch {
		case '\\':
			curr++
			if curr < len(p.buf) && (p.buf[curr] == '\'' || p.buf[curr] == '"') {
				out.Str = append(out.Str, p.buf[curr])
				curr++
			} else {
				return errors.New("bad escape")
			}
		case quote:
			out.Type = TypeStr
			p.pos = curr + 1
			return nil
		default:
			out.Str = append(out.Str, p.buf[curr])
			curr++
		}
	}
	return errors.New("string is not terminated")
}

func (p *Parser) parseInt(out *Cell) (err error) {
	start, curr := p.pos, p.pos
	if p.buf[curr] == '-' || p.buf[curr] == '+' {
		curr++
	}

	for curr < len(p.buf) && isDigit(p.buf[curr]) {
		curr ++
	}

	if out.I64, err = strconv.ParseInt(p.buf[start:curr], 10, 64); err != nil {
		return err
	}
	out.Type = TypeI64
	return nil
}

func (p *Parser) isEnd() bool {
	p.skipSpaces()
	return p.pos >= len(p.buf)
}
