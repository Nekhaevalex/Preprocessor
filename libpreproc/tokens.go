package libpreproc

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
	//LOC - current location in memory
	LOC

	//IDENT - Literals
	IDENT
	//QUOTE - "
	QUOTE

	//Misc characters

	//COMMA - ,
	COMMA
	//SEMI - ;
	SEMI
	//COLON - :
	COLON

	//Keywords

	/* Preprocessor keywords */

	//SECTION - section
	SECTION
	//EOS - end of section
	EOS
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
	// //LINE - #line
	// LINE

	//WARN - #warn
	WARN
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

	/*Assembler keywords*/

	//ADD - e.g. add a, b
	ADD
	//MOV - e.g. mov a, b
	MOV
	//IN - e.g. in a
	IN
	//OUT - e.g. out a
	OUT
	//CMP - e.g. cmp a, b, 1
	CMP
	//JMP - e.g. jmp label
	JMP
	//JNC - e.g. jnc label
	JNC

	/*Reg keywords*/

	//A - a reg
	A
	//B - b reg
	B
	//PC - pc reg
	PC
	//NR - no reg
	NR
)
