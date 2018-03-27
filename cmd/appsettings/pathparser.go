package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

func parsePath(s string) ([]string, error) {
	r := bufio.NewReader(strings.NewReader(s))

	x := Scanner{false, r}

	out := []string{}

	prev := TEof
	for {
		lex := x.Scan()

		if lex.Type == TEof {
			break
		} else if lex.Type == TIllegal {
			return out, fmt.Errorf("invalid path at '%s'", lex.Value)
		}

		if prev == TEof || prev == TSep && lex.Type == TString {
			out = append(out, lex.Value)
			prev = lex.Type
			continue
		}

		if prev == TEof || prev == TString && lex.Type == TSep {
			prev = lex.Type
			continue
		}

		return out, fmt.Errorf("invalid path at '%s':%d", lex.Value, lex.Type)
	}

	return out, nil
}

type TokenType int

const (
	TIllegal TokenType = iota
	TEof

	TString
	TSep
)

const eof = rune(0)

type Lexeme struct {
	Type  TokenType
	Value string
}

type Scanner struct {
	init bool
	r    *bufio.Reader
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}

	return ch
}

func (s *Scanner) peek() rune {
	r := s.read()
	s.unread()

	return r
}

func (s *Scanner) unread() {
	err := s.r.UnreadRune()
	if err != nil {
		panic(err)
	}
}

func (s *Scanner) Scan() Lexeme {
	// Read the next rune.
	ch := s.read()

	if isLetter(ch) || isDigit(ch) {
		s.unread()
		return s.scanIdent()
	} else if ch == '"' {
		s.unread()
		return s.scanString()
	}

	switch ch {
	case eof:
		return Lexeme{TEof, ""}
	case '.':
		return Lexeme{TSep, string(ch)}
	}

	return Lexeme{TIllegal, string(ch)}
}

func (s *Scanner) scanString() Lexeme {
	var buf bytes.Buffer

	mark := s.read()
	prev := eof

	for {
		if ch := s.read(); ch == eof {
			return Lexeme{TIllegal, string(ch)}
		} else if ch == mark && prev != '\\' {
			break
		} else {
			buf.WriteRune(ch)
			prev = ch
		}
	}

	out := strings.Replace(buf.String(), `\"`, `"`, -1)

	return Lexeme{TString, out}
}

func (s *Scanner) scanIdent() Lexeme {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	// Otherwise return as a regular identifier.
	return Lexeme{TString, buf.String()}
}

// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool {
	// unicode class n includes junk we don't want
	return (ch >= '0' && ch <= '9')
}
