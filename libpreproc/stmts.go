package libpreproc

//Reg - register type
type Reg int

const (
	a Reg = 0
	b Reg = 1
)

//PreprocIdent - preprocessor identifier
type PreprocIdent interface {
}

//PreprocStmt - preprocessor statement
type PreprocStmt interface {
}

//PreprocProg - instance of parsed code
type PreprocProg struct {
	Body []PreprocStmt
}

//ImportStmt - import statement
type ImportStmt struct {
	ImportFile string
}

//DefStmt - define statement
type DefStmt struct {
	DefinitionName string
	ToDef          PreprocStmt
}

//PextStmt - pext declaration statement
type PextStmt struct {
	PextName    string
	PextAddress string
}

//ErrorStmt - error message statement
type ErrorStmt struct {
	ErrorMsg string
}

//PragmaStmt - pragma stmt
type PragmaStmt struct {
	PragmaName string
}

//LineStmt - insert line stmt
type LineStmt struct {
	LineNumber string
	FileName   string
}

//MsgStmt - message stmt
type MsgStmt struct {
	Msg string
}

//IfStmt - if stmt
type IfStmt struct {
	DefName     string
	Negative    bool
	Branch1body []PreprocStmt
	Branch2body []PreprocStmt
}

//SumStmt - sumdef
type SumStmt struct {
	X string
	Y string
}

//ResStmt - resdef
type ResStmt struct {
	X string
	Y string
}

//UndefStmt - undef
type UndefStmt struct {
	DefName string
}

//ReturnStmt - return
type ReturnStmt struct {
	ReturnName string
}

//MacroStmt - macro
type MacroStmt struct {
	Name string
	Vars []string
	Body []PreprocStmt
}

//AsmStmt - assembler stmt
type AsmStmt interface {
}

//AddStmt - add
type AddStmt struct {
	arg1 Reg
	fa   int
}

//MovStmt - mov
type MovStmt struct {
	arg1 Reg
	arg2 Reg
	fa   PreprocIdent
}

//InStmt - in
type InStmt struct {
	arg Reg
}

//OutStmt - out
type OutStmt struct {
	arg PreprocIdent
}

//CmpStmt - cmp
type CmpStmt struct {
	arg1 Reg
	arg2 Reg
	fa   PreprocIdent
}

//JmpStmt - jmp
type JmpStmt struct {
	arg PreprocIdent
}

//JncStmt - jnc
type JncStmt struct {
	arg PreprocIdent
}
