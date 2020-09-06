package preprocessor

//PreprocStmt - preprocessor statement
type PreprocStmt interface{}

//ImportStmt - import statement
type ImportStmt struct {
	importFile string
}

//DefStmt - define statement
type DefStmt struct {
	definitionName string
	toDef          PreprocStmt
}

//PextStmt - pext declaration statement
type PextStmt struct {
	pextName    string
	pextAddress string
}

//ErrorStmt - error message statement
type ErrorStmt struct {
	errorMsg string
}

//PragmaStmt - pragma stmt
type PragmaStmt struct {
	pragmaName string
}

//LineStmt - insert line stmt
type LineStmt struct {
	lineNumber string
	fileName   string
}

//MsgStmt - message stmt
type MsgStmt struct {
	msg string
}

//IfStmt - if stmt
type IfStmt struct {
	defName     string
	negative    bool
	branch1body []PreprocStmt
	branch2body []PreprocStmt
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
	defName string
}

//ReturnStmt - return
type ReturnStmt struct {
	returnName string
}

//MacroStmt - macro
type MacroStmt struct {
	name string
	vars []string
	body []PreprocStmt
}
