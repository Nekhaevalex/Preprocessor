package libpreproc

import "errors"

//ErrElseBranch - message that end of branch was met
var ErrElseBranch = errors.New("ElseBranch")

//ErrEndIf - message that end of if was reached
var ErrEndIf = errors.New("EndIf")

//ErrMacroEnd - message that end of Macro was reached
var ErrMacroEnd = errors.New("MacroEnd")
