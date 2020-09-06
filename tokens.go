package preprocessor

//Token represents a lexical token
type Token int

const (
	//Special tokens

	//ILLEGAL - illegal character
	ILLEGAL Token = iota
	//EOF - End of file
	EOF
	//WS - White space
	WS

	//IDENT - Literals
	IDENT

	//Misc characters

	//COMMA - ,
	COMMA

	//Keywords

	//IMPORT - #import
	IMPORT
	//DEFINE - #define
	DEFINE
	//PEXT - #pext
	PEXT
	//ERROR - #error
	ERROR
	//PRAGMA - #pragma
	PRAGMA
	//LINE - #line
	LINE
	//MESSAGE - #message
	MESSAGE
	//IFDEF - #ifdef
	IFDEF
	//IFNDEF - #ifndef
	IFNDEF
	//ENDIF - #endif
	ENDIF
	//ELSE - #else
	ELSE
	//SUMDEF - #sumdef
	SUMDEF
	//RESDEF - #resdef
	RESDEF
	//UNDEF - #undef
	UNDEF
	//RETURN - #return
	RETURN
	//MACRO - #macro
	MACRO
	//ENDMACRO - #endmacro
	ENDMACRO
)
