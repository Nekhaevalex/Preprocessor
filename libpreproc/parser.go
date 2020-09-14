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

//ParseFile parses the whole file
func (p *Parser) ParseFile() (PreprocProg, error) {
	var prog PreprocProg
	for {
		stmt, err := p.Parse()
		if stmt == EOF {
			break
		}
		if err != nil {
			return prog, err
		}
		prog.Body = append(prog.Body, PreprocStmt(stmt))
	}
	return prog, nil
}

//Parse parses the file
func (p *Parser) Parse() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	var stmt PreprocStmt
	var er error
	switch tok {
	case IMPORT:
		stmt, er = p.parseImport()
	case DEFINE:
		stmt, er = p.parseDefine()
	case PEXT:
		stmt, er = p.parsePext()
	case ERROR:
		stmt, er = p.parseError()
	case PRAGMA:
		stmt, er = p.parsePragma()
	case LINE:
		stmt, er = p.parseLine()
	case MESSAGE:
		stmt, er = p.parseMsg()
	case IFDEF:
		stmt, er = p.parseIf(false)
	case IFNDEF:
		stmt, er = p.parseIf(true)
	case ENDIF:
		return nil, ErrEndIf
	case ELSE:
		return nil, ErrElseBranch
	case SUMDEF:
		stmt, er = p.parseSum()
	case RESDEF:
		stmt, er = p.parseRes()
	case UNDEF:
		stmt, er = p.parseUndef()
	case RETURN:
		stmt, er = p.parseReturn()
	case MACRO:
		stmt, er = p.parseMacro()
	case ENDMACRO:
		return nil, ErrMacroEnd
	case EOF:
		return EOF, nil
	case IDENT:
		return IDENT, nil
	default:
		fmt.Printf("Unknown token: %s\n", lit)
	}
	return stmt, er
}

func (p *Parser) parseImport() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}

	return &ImportStmt{ImportFile: lit}, nil
}

func (p *Parser) parseDefine() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	defName := lit
	toDef, _ := p.Parse()
	return &DefStmt{DefinitionName: defName, ToDef: toDef}, nil
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
	return &PextStmt{PextName: pextName, PextAddress: pextAdr}, nil
}

func (p *Parser) parseError() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	errorMsg := lit
	return &ErrorStmt{ErrorMsg: errorMsg}, nil
}

func (p *Parser) parsePragma() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	pragma := lit
	return &PragmaStmt{PragmaName: pragma}, nil
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
	return &LineStmt{LineNumber: lineNum, FileName: fileName}, nil
}

func (p *Parser) parseMsg() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	return &MsgStmt{Msg: lit}, nil
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
	return &IfStmt{DefName: defName, Negative: neg, Branch1body: branchIf, Branch2body: branchElse}, nil
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
	return &UndefStmt{DefName: lit}, nil
}

func (p *Parser) parseReturn() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	return &ReturnStmt{ReturnName: lit}, nil
}

func (p *Parser) parseMacro() (PreprocStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected ident", lit)
	}
	name := lit
	var vars []string
	for {
		tok, lit = p.scanIgnoreWhitespace()
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected ident", lit)
		}
		vars = append(vars, lit)
		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}
	}
	var body []PreprocStmt
	body = p.parseBody()
	return &MacroStmt{Name: name, Vars: vars, Body: body}, nil
}

func (p *Parser) parseBody() []PreprocStmt {
	var body []PreprocStmt
	for {
		stmt, err := p.Parse()
		if err == ErrMacroEnd {
			break
		}
		body = append(body, stmt)
	}
	return body
}

func (p *Parser) parseAdd() (AsmStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	var arg1 Reg
	var arg2 Reg
	switch tok {
	case A:
		arg1 = a
	case B:
		arg1 = b
	default:
		fmt.Printf("Unknown identifier: %s\n", lit)
	}
	tok, lit = p.scanIgnoreWhitespace()
	if tok == COMMA {
		tok, lit = p.scanIgnoreWhitespace()
	}
	switch tok {
	case A:
		arg2 = a
	case B:
		arg2 = b
	default:
		fmt.Printf("Unknown identifier: %s\n", lit)
	}

}
