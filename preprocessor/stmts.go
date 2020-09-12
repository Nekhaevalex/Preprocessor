package preprocessor

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
