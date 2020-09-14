package libpreproc

import (
	"bufio"
	"bytes"
	"io"
)

func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '#'
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

var eof = rune(0)

//Scanner - represents a lexical scanner/
type Scanner struct {
	r *bufio.Reader
}

//NewScanner - returns a new instance of Scanner
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

//Returns the rune(0) if error occurs (or io.EOF is returned)
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

//unread places the previously read rune back on the reader
func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

//Scan returns the next token and literal value
func (s *Scanner) Scan() (tok Token, lit string) {
	//Read the next rune.
	ch := s.read()

	if isWhiteSpace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	} else if isDigit(ch) {
		s.unread()
		return s.scanIdent()
	}

	switch ch {
	case eof:
		return EOF, ""
	case ',':
		return COMMA, string(ch)
	}
	return ILLEGAL, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	//Create a buffer and read the current character into it
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	//Read every subsequent whitespaces character into the buffer
	//Non-whitespace characters and EOF will cause the loop to exit
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhiteSpace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}
	return WS, buf.String()
}

// scanIdent consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanIdent() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	switch buf.String() {
	case "#import":
		return IMPORT, buf.String()
	case "#define":
		return DEFINE, buf.String()
	case "#pext":
		return PEXT, buf.String()
	case "#error":
		return ERROR, buf.String()
	case "#pragma":
		return PRAGMA, buf.String()
	case "#line":
		return LINE, buf.String()
	case "#warn":
		return MESSAGE, buf.String()
	case "#ifdef":
		return IFDEF, buf.String()
	case "#ifndef":
		return IFNDEF, buf.String()
	case "#endif":
		return ENDIF, buf.String()
	case "#else":
		return ELSE, buf.String()
	case "#sumdef":
		return SUMDEF, buf.String()
	case "#resdef":
		return RESDEF, buf.String()
	case "#undef":
		return UNDEF, buf.String()
	case "#return":
		return RETURN, buf.String()
	case "#macro":
		return MACRO, buf.String()
	case "#endmacro":
		return ENDMACRO, buf.String()
	case "add":
		return ADD, buf.String()
	case "mov":
		return MOV, buf.String()
	case "in":
		return IN, buf.String()
	case "out":
		return OUT, buf.String()
	case "cmp":
		return CMP, buf.String()
	case "jmp":
		return JMP, buf.String()
	case "jnc":
		return JNC, buf.String()
	case "a":
		return A, buf.String()
	case "b":
		return B, buf.String()
	}
	return IDENT, buf.String()
}
