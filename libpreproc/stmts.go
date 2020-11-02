package libpreproc

//Reg - register type
type Reg int

const (
	nr Reg = -1
	a  Reg = 0
	b  Reg = 1
	pc Reg = 2
)

//Ident - identifier
type Ident interface {
}

//SimpleString - specified identifier containing string
// e.g. "Hello, World"
type SimpleString struct {
	value string
}

//Number - specified identifier contatinig number
// e.g. 42
type Number struct {
	value int
}

//Variable - specified identifier contatining preprocessor
// value e.g. A
type Variable struct {
	name string
}

//Program - program object
type Program struct {
	sections []Section
}

//Stmt - program statement (Section/Directive/Opcode/Label)
type Stmt interface {
}

//Section - program section
type Section struct {
	sectionName    string
	sectionContent Block
}

//Block - program block
type Block struct {
	elements []Stmt
}

//Directives

//Define - #define
type Define struct {
	name       Ident
	definition Ident
}

//Import - #import
type Import struct {
	name Ident
}

// //Line - #line
// type Line struct {
// 	name       Ident
// 	lineNumber Ident
// }

//Warn - #warn
type Warn struct {
	message Ident
}

//Sumdef - #sumdef {
type Sumdef struct {
	def1 Ident
	def2 Ident
}

//Resdef - #resdef
type Resdef struct {
	def1 Ident
	def2 Ident
}

//Pext - #pext
type Pext struct {
	pextName    Ident
	pextAddress Ident
}

//Error - #error
type Error struct {
	message Ident
}

//Undef - #undef
type Undef struct {
	definition Ident
}

//Ifdef - #ifdef
type Ifdef struct {
	definition Ident
	bodyTrue   Block
	bodyFalse  Block
}

//Ifndef - #ifdef
type Ifndef struct {
	definition Ident
	bodyTrue   Block
	bodyFalse  Block
}

//Macro - #macro
type Macro struct {
	macroName string
	args      []string
	body      Block
}

//Return - #return
type Return struct {
	returnValue Ident
}

//Opcodes

//Opcode - opcode interface
type Opcode interface {
}

//Add - add
type Add struct {
	reg   Reg
	value Ident
}

//Mov - mov
type Mov struct {
	reg1 Reg
	reg2 Reg
	fa   Ident
}

//In - in
type In struct {
	reg Reg
}

//Out - out
type Out struct {
	reg Reg
	fa  Ident
}

//Cmp - cmp
type Cmp struct {
	regA      Reg
	regB      Reg
	operation Ident
}

//Jmp - jmp
type Jmp struct {
	regB Reg
	addr Ident
}

//Jnc - jnc
type Jnc struct {
	regB Reg
	addr Ident
}

//MacroCall - macro call
type MacroCall struct {
	macroName string
	args      []Ident
}

//Label - label
type Label struct {
	name Ident
}
