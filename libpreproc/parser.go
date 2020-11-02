package libpreproc

import (
	"fmt"
	"io"
	"strconv"
)

//Module represents program with all its imports
type Module struct {
	macroList []string
	labelList []string
	main      Program
	imports   []Program
}

//Parser represents a parser
type Parser struct {
	macroList []string
	labelList []string
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
		if stmt == EOF || stmt == ENDIF || stmt == EOS || stmt == ELSE {
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
		stmt, er = EOS, nil
		p.unscan()
	case DEFINE:
		stmt, er = p.ParseDefine()
	case IMPORT:
		stmt, er = p.ParseImport()
	// case LINE:
	// 	stmt, er = p.ParseLine()
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
		stmt = ELSE
		er = nil
	case ENDIF:
		p.unscan()
		stmt = ENDIF
		er = nil
	case RETURN:
		stmt, er = p.ParseReturn()
	case MACRO:
		stmt, er = p.ParseMacro()
	case ENDMACRO:
		stmt = ENDMACRO
		er = nil
	case ADD:
		stmt, er = p.ParseAdd()
	case MOV:
		stmt, er = p.ParseMov()
	case IN:
		stmt, er = p.ParseIn()
	case OUT:
		stmt, er = p.ParseOut()
	case CMP:
		stmt, er = p.ParseCmp()
	case JMP:
		stmt, er = p.ParseJmp()
	case JNC:
		stmt, er = p.ParseJnc()
	case IDENT:
		p.unscan()
		stmt, er = p.ParseIdent()
	case LOC:
		stmt = LOC
		er = nil
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
	return Define{name: name, definition: definition}, nil
}

//ParseImport - #import
func (p *Parser) ParseImport() (Stmt, error) {
	name, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	return Import{name: name}, nil
}

// //ParseLine - #line
// func (p *Parser) ParseLine() (Stmt, error) {
// 	name, err := p.ParseIdent()
// 	if err != nil {
// 		return nil, err
// 	}
// 	lineNumber, err := p.ParseIdent()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return Line{name: name, lineNumber: lineNumber}, nil
// }

//ParseWarn - #warn
func (p *Parser) ParseWarn() (Stmt, error) {
	message, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	return Warn{message: message}, nil
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
	return Sumdef{def1: def1, def2: def2}, nil
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
	return Resdef{def1: def1, def2: def2}, nil
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
	return Pext{pextName: pextName, pextAddress: pextAddress}, nil
}

//ParseError - #error
func (p *Parser) ParseError() (Stmt, error) {
	message, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	return Error{message: message}, nil
}

//ParseUndef - #undef
func (p *Parser) ParseUndef() (Stmt, error) {
	definition, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	return Undef{definition: definition}, nil
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
		_, _ = p.scanIgnoreWhitespace()
	} else {
		if tok != ENDIF {
			p.unscan()
		}
	}
	return Ifdef{definition: definition, bodyTrue: bodyTrue, bodyFalse: bodyFalse}, nil
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
	return Ifndef{definition: definition, bodyTrue: bodyTrue, bodyFalse: bodyFalse}, nil
}

//ParseReturn - #return
func (p *Parser) ParseReturn() (Stmt, error) {
	returnName, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	return Return{returnValue: returnName}, nil
}

//ParseIdent - parses any unknown words: variables, strings, numbers,
//labels and macro calls
func (p *Parser) ParseIdent() (Ident, error) {
	tok, ident := p.scanIgnoreWhitespace()
	switch tok {
	case IDENT: //macrocall, variable, number
		//test for number
		isNum, num := numberIdent(ident)
		if isNum {
			return Number{value: num}, nil
		}
		tok, _ := p.scanIgnoreWhitespace()
		switch tok {
		//Test if it's a label: next token should be COLON
		case COLON:
			p.labelList = append(p.labelList, ident)
			colErr := p.checkLabelMacroCollision(ident)
			if colErr != nil {
				return nil, colErr
			}
			return Label{name: Variable{name: ident}}, nil
		default:
			p.unscan()
			_, foundMacro := find(p.macroList, ident)
			if foundMacro {
				return p.ParseMacroCall(ident)
			}
			label, foundLabel := find(p.labelList, ident)
			if foundLabel {
				return Label{name: Variable{name: p.labelList[label]}}, nil
			}
			return Variable{name: ident}, nil
		}
	case QUOTE: //SimpleString
		str := p.getStringValue()
		return SimpleString{value: str}, nil
	default: //Else (???)
		return nil, fmt.Errorf("forbidden symbol %q in context", ident)
	}
}

//ParseMacro - #macro
func (p *Parser) ParseMacro() (Stmt, error) {
	tok, macroName := p.scanIgnoreWhitespace()
	colErr := p.checkLabelMacroCollision(macroName)
	if colErr != nil {
		return nil, colErr
	}
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
	colErr = p.checkLabelMacroCollision(macroName)
	if colErr != nil {
		return nil, colErr
	}
	return Macro{macroName: macroName, args: args, body: body}, nil
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
func (p *Parser) ParseMacroCall(macroName string) (Stmt, error) {
	var args []Ident
	tok, arg := p.scan()
	for {
		if tok == IDENT {
			p.unscan()
			stmt, er := p.ParseIdent()
			if er != nil {
				return nil, er
			}
			args = append(args, stmt)
		} else if tok == WS {
			if hasNewLine(arg) == true {
				break
			}
		}
		tok, arg = p.scan()
	}
	return MacroCall{macroName: macroName, args: args}, nil
}

//ParseAdd - add
func (p *Parser) ParseAdd() (Opcode, error) {
	reg, err := p.ParseReg()
	if err != nil {
		return nil, err
	}
	tok, _ := p.scanIgnoreWhitespace()
	if tok != COMMA {
		p.unscan()
	}
	value, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	return Add{reg: reg, value: value}, nil
}

//ParseMov - mov
func (p *Parser) ParseMov() (Opcode, error) {
	reg1, err := p.ParseReg()
	if err != nil {
		return nil, err
	}
	tok, _ := p.scanIgnoreWhitespace()
	if tok != COMMA {
		p.unscan()
	}
	tok, lit := p.scanIgnoreWhitespace()
	p.unscan()
	reg2 := nr
	var val Ident
	if tok == A || tok == B || tok == PC {
		reg2, err = p.ParseReg()
		if err != nil {
			return nil, err
		}
		tok, _ = p.scanIgnoreWhitespace()
		p.unscan()
		if tok == IDENT {
			val, err = p.ParseIdent()
			if err != nil {
				return nil, err
			}
		}
	} else if tok == IDENT {
		val, err = p.ParseIdent()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("expected reg or ident, met %q", lit)
	}
	return Mov{reg1: reg1, reg2: reg2, fa: val}, nil
}

//ParseIn - in
func (p *Parser) ParseIn() (Opcode, error) {
	reg, er := p.ParseReg()
	if er != nil {
		return nil, er
	}
	return In{reg: reg}, nil
}

//ParseOut - out
func (p *Parser) ParseOut() (Opcode, error) {
	tok, lit := p.scanIgnoreWhitespace()
	p.unscan()
	reg := nr
	var val Ident
	var err error
	if tok == A || tok == B {
		reg, err = p.ParseReg()
		if err != nil {
			return nil, err
		}
		tok, _ = p.scanIgnoreWhitespace()
		p.unscan()
		if tok == IDENT {
			val, err = p.ParseIdent()
			if err != nil {
				return nil, err
			}
		}
	} else if tok == IDENT {
		val, err = p.ParseIdent()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("expected reg a, reg b or ident, met %q", lit)
	}
	return Out{reg: reg, fa: val}, nil
}

//ParseCmp - cmp
func (p *Parser) ParseCmp() (Opcode, error) {
	regA, err := p.ParseReg()
	if err != nil {
		return nil, err
	}
	tok, _ := p.scanIgnoreWhitespace()
	if tok != COMMA {
		p.unscan()
	}
	regB, err := p.ParseReg()
	if err != nil {
		return nil, err
	}
	tok, _ = p.scanIgnoreWhitespace()
	if tok != COMMA {
		p.unscan()
	}
	op, err := p.ParseIdent()
	return Cmp{regA: regA, regB: regB, operation: op}, nil
}

//ParseJmp - jmp
func (p *Parser) ParseJmp() (Opcode, error) {
	tok, _ := p.scanIgnoreWhitespace()
	if tok == IDENT {
		p.unscan()
		addr, err := p.ParseIdent()
		if err != nil {
			return nil, err
		}
		return Jmp{regB: nr, addr: addr}, nil
	}
	p.unscan()
	reg, err := p.ParseReg()
	if err != nil {
		return nil, err
	}
	return Jmp{regB: reg, addr: nil}, nil
}

//ParseJnc - jnc
func (p *Parser) ParseJnc() (Opcode, error) {
	tok, _ := p.scanIgnoreWhitespace()
	if tok == IDENT {
		p.unscan()
		addr, err := p.ParseIdent()
		if err != nil {
			return nil, err
		}
		return Jnc{regB: nr, addr: addr}, nil
	}
	p.unscan()
	reg, err := p.ParseReg()
	if err != nil {
		return nil, err
	}
	return Jnc{regB: reg, addr: nil}, nil
}

//ParseReg - parses any register
func (p *Parser) ParseReg() (Reg, error) {
	tok, reg := p.scanIgnoreWhitespace()
	switch tok {
	case A:
		return a, nil
	case B:
		return b, nil
	case PC:
		return pc, nil
	default:
		return nr, fmt.Errorf("expected register, met %q", reg)
	}
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

func (p *Parser) checkLabelMacroCollision(ident string) error {
	_, isLabel := find(p.labelList, ident)
	_, isMacro := find(p.macroList, ident)
	if isMacro && isLabel {
		return fmt.Errorf("ident %q already exists", ident)
	}
	return nil
}
