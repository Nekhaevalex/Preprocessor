package preprocessor

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
	return
}

//Parse parses the file
func (p *Parser) Parse() (PreprocStmt, error) {
	tok, _ := p.scanIgnoreWhitespace()
	var stmt PreprocStmt
	var er error
	switch tok {
	case IMPORT:
		stmt, er = p.parseImport()
	case DEFINE:
		stmt, er = p.parseDefine()
	}
	return stmt, er
}

func (p *Parser) parseImport() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}

	return &ImportStmt{importFile: lit}, nil
}

func (p *Parser) parseDefine() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	defName := lit
	toDef, _ := p.Parse()
	return &DefStmt{definitionName: defName, toDef: toDef}, nil
}

func (p *Parser) parsePext() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	pextName := lit
	tok, lit = p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	pextAdr := lit
	return &PextStmt{pextName: pextName, pextAddress: pextAdr}, nil
}

func (p *Parser) parseError() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	errorMsg := lit
	return &ErrorStmt{errorMsg: errorMsg}, nil
}

func (p *Parser) parsePragma() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	pragma := lit
	return &PragmaStmt{pragmaName: pragma}, nil
}

func (p *Parser) parseLine() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	lineNum := lit
	tok, lit = p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	fileName := lit
	return &LineStmt{lineNumber: lineNum, fileName: fileName}, nil
}

func (p *Parser) parseMsg() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	return &MsgStmt{msg: lit}, nil
}

func (p *Parser) parseIf(negative bool) (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	defName := lit
	neg := negative
	var branchIf []PreprocStmt
	var branchElse []PreprocStmt
	var errRes error
	var stmt PreprocStmt
	var hasSecondBranch bool
	for {
		stmt, errRes = p.Parse()
		if errRes == ErrElseBranch {
			hasSecondBranch = true
			break
		}
		if errRes == ErrEndIf {
			hasSecondBranch = false
			break
		}
		branchIf = append(branchIf, stmt)
	}
	if hasSecondBranch {
		for {
			stmt, errRes = p.Parse()
			if errRes == ErrEndIf {
				break
			}
			branchElse = append(branchElse, stmt)
		}
	} else {
		branchElse = nil
	}
	return &IfStmt{defName: defName, negative: neg, branch1body: branchIf, branch2body: branchElse}, nil
}

func (p *Parser) parseSum() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	arg1 := lit
	tok, lit = p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	arg2 := lit
	return &SumStmt{X: arg1, Y: arg2}, nil
}

func (p *Parser) parseRes() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	arg1 := lit
	tok, lit = p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	arg2 := lit
	return &ResStmt{X: arg1, Y: arg2}, nil
}

func (p *Parser) parseUndef() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	return &UndefStmt{defName: lit}, nil
}

func (p *Parser) parseReturn() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	return &ReturnStmt{returnName: lit}, nil
}

func (p *Parser) parseMacro() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	name := lit
	var vars []string
	tok, vars = p.parseVars()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	var body []PreprocStmt
	tok, body = p.parseBody()
	return &MacroStmt{name: name, vars: vars, body: body}, nil
}

func (p *Parser) parseVars() (PreprocStmt, []string) {

}
