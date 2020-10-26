package libpreproc

import (
	"fmt"
	"io"
	"strconv"
)

//Parser represents a parser
type Parser struct {
	macroList []string
	s         *Scanner
	buf       struct {
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
	tok, lit := p.scan()
	var str string
	for tok != QUOTE {
		str = str + lit
		tok, lit = p.scan()
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
	case ERROR:
		stmt, er = p.ParseError()
	case UNDEF:
		stmt, er = p.ParseUndef()
	case IFDEF:
		stmt, er = p.ParseIfdef()
	case IFNDEF:
		stmt, er = p.ParseIfndef()
	case EOF:
		stmt = EOF
		er = nil
	case ELSE:
		p.unscan()
		stmt = EOF
		er = nil
	case ENDIF:
		stmt = EOF
		er = nil
	case RETURN:
		stmt, er = p.ParseReturn()
	case MACRO:
		stmt, er = p.ParseMacro()
	case ENDMACRO:
		stmt = EOF
		er = nil
	case IDENT:
		p.unscan()
		stmt, er = p.ParseIdent()
	}
	switch stmt.(type) {
	case *Variable:
		varVal := stmt.(*Variable).name
		return nil, fmt.Errorf("preprocessor directive expected, met variable: %q", varVal)
	}
	return stmt, er
}

//ParseDefine - #define
func (p *Parser) ParseDefine() (Stmt, error) {
	name, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	definition, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	return &Define{name: name, definition: definition}, nil
}

//ParseImport - #import
func (p *Parser) ParseImport() (Stmt, error) {
	name, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	return &Import{name: name}, nil
}

//ParseLine - #line
func (p *Parser) ParseLine() (Stmt, error) {
	name, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	lineNumber, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	return &Line{name: name, lineNumber: lineNumber}, nil
}

//ParseWarn - #warn
func (p *Parser) ParseWarn() (Stmt, error) {
	message, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	return &Warn{message: message}, nil
}

//ParseSumDef - #sumdef
func (p *Parser) ParseSumDef() (Stmt, error) {
	def1, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	def2, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	return &Sumdef{def1: def1, def2: def2}, nil
}

//ParseResDef - #resdef
func (p *Parser) ParseResDef() (Stmt, error) {
	def1, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	def2, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	return &Resdef{def1: def1, def2: def2}, nil
}

//ParsePext - #pext
func (p *Parser) ParsePext() (Stmt, error) {
	pextName, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	pextAddress, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	return &Pext{pextName: pextName, pextAddress: pextAddress}, nil
}

//ParseError - #error
func (p *Parser) ParseError() (Stmt, error) {
	message, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	return &Error{message: message}, nil
}

//ParseUndef - #undef
func (p *Parser) ParseUndef() (Stmt, error) {
	definition, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	return &Undef{definition: definition}, nil
}

//ParseIfdef - #ifdef
func (p *Parser) ParseIfdef() (Stmt, error) {
	definition, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	bodyTrue, err := p.ParseBlock()
	tok, _ := p.scanIgnoreWhitespace()
	var bodyFalse Block
	if tok == ELSE {
		bodyFalse, err = p.ParseBlock()
		if err != nil {
			return nil, err
		}
	} else {
		p.unscan()
	}
	return &Ifdef{definition: definition, bodyTrue: bodyTrue, bodyFalse: bodyFalse}, nil
}

//ParseIfndef - #ifdef
func (p *Parser) ParseIfndef() (Stmt, error) {
	definition, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	bodyTrue, err := p.ParseBlock()
	tok, _ := p.scanIgnoreWhitespace()
	var bodyFalse Block
	if tok == ELSE {
		bodyFalse, err = p.ParseBlock()
		if err != nil {
			return nil, err
		}
	} else {
		p.unscan()
	}
	return &Ifndef{definition: definition, bodyTrue: bodyTrue, bodyFalse: bodyFalse}, nil
}

//ParseReturn - #return
func (p *Parser) ParseReturn() (Stmt, error) {
	returnName, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	return &Return{returnValue: returnName}, nil
}

//ParseIdent - parses any unknown words: variables, strings, numbers,
//labels and macro calls
func (p *Parser) ParseIdent() (Stmt, error) {
	tok, ident := p.scanIgnoreWhitespace()
	switch tok {
	case IDENT: //macrocall, variable, number
		//test for number
		isNum, num := numberIdent(ident)
		if isNum {
			return &Number{value: num}, nil
		}
		tok, _ := p.scanIgnoreWhitespace()
		switch tok {
		//Test if it's a label: next token should be COLON
		case COLON:
			return &Label{name: &Variable{name: ident}}, nil
		default:
			p.unscan()
			_, foundMacro := find(p.macroList, ident)
			if foundMacro {
				return p.ParseMacroCall()
			}
			return &Variable{name: ident}, nil
		}
	case QUOTE: //SimpleString
		str := p.getStringValue()
		return &SimpleString{value: str}, nil
	default: //Else (???)
		return nil, fmt.Errorf("forbidden symbol %q in context", ident)
	}
}

//ParseMacro - #macro
func (p *Parser) ParseMacro() (Stmt, error) {
	tok, macroName := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("macro name expected, met %q", macroName)
	}
	var args []string
	tok, arg := p.scan()
	for {
		if tok == IDENT {
			args = append(args, arg)
		} else if tok == WS {
			if hasNewLine(arg) == true {
				break
			}
		}
		tok, arg = p.scan()
	}
	body, err := p.ParseBlock()
	if err != nil {
		return nil, err
	}
	p.rememberMacro(macroName)
	return &Macro{macroName: macroName, args: args, body: body}, nil
}

func hasNewLine(str string) bool {
	s := []rune(str)
	for _, r := range s {
		if r == 10 {
			return true
		}
	}
	return false
}

//ParseMacroCall - parses any macro call
func (p *Parser) ParseMacroCall() (Stmt, error) {
	return nil, nil
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func numberIdent(num string) (bool, int) {
	runeRepr := []rune(num)
	if runeRepr[0] == '0' {
		switch runeRepr[1] {
		case 'x':
			pnum, err := strconv.ParseInt(string(runeRepr[2:]), 16, 64)
			if err != nil {
				return false, 0
			}
			return true, int(pnum)
		case 'o':
			pnum, err := strconv.ParseInt(string(runeRepr[2:]), 8, 64)
			if err != nil {
				return false, 0
			}
			return true, int(pnum)
		case 'b':
			pnum, err := strconv.ParseInt(string(runeRepr[2:]), 2, 64)
			if err != nil {
				return false, 0
			}
			return true, int(pnum)
		}
	} else {
		pnum, err := strconv.Atoi(num)
		if err != nil {
			return false, 0
		}
		return true, pnum
	}
	return false, 0
}

func (p *Parser) rememberMacro(macroName string) {
	p.macroList = append(p.macroList, macroName)
}
