package libpreproc

import (
	"fmt"
	"io"
)

//Parser represents a parser
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  //last read token
		lit string //last read literal
		n   int    //buffer size (max = 1)
	}
}

//NewParser returns a new instance of Parser
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	//Otherwise read the next token from the scanner
	tok, lit = p.s.Scan()

	//Save it to the buffer in case we unscan later
	p.buf.tok, p.buf.lit = tok, lit
	return
}

//unscan pushes the prev read token back to buffer
func (p *Parser) unscan() {
	p.buf.n = 1
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return tok, lit
}

func (p *Parser) getStringValue() string {
	tok, lit := p.scanIgnoreWhitespace()
	var str string
	for tok != QUOTE {
		str = str + lit
		tok, lit = p.scanIgnoreWhitespace()
	}
	return str
}

//ParseFile parses the whole file
func (p *Parser) ParseFile() (Program, error) {
	var prog Program
	var er error
	for {
		tok, lit := p.scanIgnoreWhitespace()
		if tok == EOF {
			break
		}
		var section Section
		if tok == SECTION {
			tok, lit = p.scanIgnoreWhitespace()
			if tok != IDENT {
				return prog, fmt.Errorf("found %q, expected section name", lit)
			}
			section.sectionName = lit
			section.sectionContent, er = p.ParseBlock()
			if er != nil {
				return prog, er
			}
			prog.sections = append(prog.sections, section)
		}
	}
	return prog, nil
}

//ParseBlock parses one section
func (p *Parser) ParseBlock() (Block, error) {
	var block Block
	for {
		stmt, err := p.Parse()
		if stmt == EOF {
			break
		}
		if err != nil {
			return block, err
		}
		block.elements = append(block.elements, stmt)
	}
	return block, nil
}

//Parse parses all keywords and calls their handlers
func (p *Parser) Parse() (Stmt, error) {
	tok, _ := p.scanIgnoreWhitespace()
	var stmt Stmt
	var er error
	switch tok {
	case SECTION:
		p.unscan()
		stmt, er = EOF, nil
	case DEFINE:
		stmt, er = p.ParseDefine()
	case IMPORT:
		stmt, er = p.ParseImport()
	case LINE:
		stmt, er = p.ParseLine()
	case WARN:
		stmt, er = p.ParseWarn()
	case SUMDEF:
		stmt, er = p.ParseSumDef()
	case RESDEF:
		stmt, er = p.ParseResDef()
	case PEXT:
		stmt, er = p.ParsePext()
	case EOF:
		stmt = EOF
		er = nil
	}
	return stmt, er
}

//ParseDefine - #define
func (p *Parser) ParseDefine() (Stmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected definition identifier", lit)
	}
	name := lit
	tok, lit = p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected definition identifier", lit)
	}
	definition := lit
	return &Define{name: name, definition: definition}, nil
}

//ParseImport - #import
func (p *Parser) ParseImport() (Stmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	var name string
	switch tok {
	case IDENT:
		name = lit
	case QUOTE:
		name = p.getStringValue()
	default:
		return nil, fmt.Errorf("found %q, expected import identifier", lit)
	}
	return &Import{name: name}, nil
}

//ParseLine - #line
func (p *Parser) ParseLine() (Stmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected line source indetifier", lit)
	}
	name := lit
	tok, lit = p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected line number identifier", lit)
	}
	lineNumber := lit
	return &Line{name: name, lineNumber: lineNumber}, nil
}

//ParseWarn - #warn
func (p *Parser) ParseWarn() (Stmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	var message string
	switch tok {
	case IDENT:
		message = lit
	case QUOTE:
		message = p.getStringValue()
	default:
		return nil, fmt.Errorf("found %q, expected import identifier", lit)
	}
	return &Warn{message: message}, nil
}

//ParseSumDef - #sumdef
func (p *Parser) ParseSumDef() (Stmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected definition identifier", lit)
	}
	def1 := lit
	tok, lit = p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q expected definition identifier", lit)
	}
	def2 := lit
	return &Sumdef{def1: def1, def2: def2}, nil
}

//ParseResDef - #resdef
func (p *Parser) ParseResDef() (Stmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected definition identifier", lit)
	}
	def1 := lit
	tok, lit = p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected definition identifier", lit)
	}
	def2 := lit
	return &Resdef{def1: def1, def2: def2}, nil
}

//ParsePext - #pext
func (p *Parser) ParsePext() (Stmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected pExt name identifier", lit)
	}
	pextName := lit
	tok, lit = p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected pExt mount point identifier", lit)
	}
	pextAddress := lit
	return &Pext{pextName: pextName, pextAddress: pextAddress}, nil
}

//ParseError - #error
