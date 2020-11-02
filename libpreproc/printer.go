package libpreproc

import "fmt"

var depth int

//PrintProg - prints program AST
func PrintProg(prg Program) {
	fmt.Printf("program\n")
	depth = 0
	for _, v := range prg.sections {
		printSection(v)
	}
}

func printSection(sec Section) {
	depth++
	fmt.Printf("\tsection\n")
	fmt.Printf("\t\tsection_name: %s\n", sec.sectionName)
	fmt.Printf("\t\tsection_data:\n")
	printBlock(sec.sectionContent)
	depth--
}

func printBlock(blk Block) {
	for _, v := range blk.elements {
		for i := 0; i < depth; i++ {
			fmt.Printf("\t\t\t")
		}
		estimateStmt(v)
	}
}

func estimateStmt(stmt interface{}) {
	switch stmt.(type) {
	case Define:
		v, _ := stmt.(Define)
		printDefine(v)
	case Import:
		v, _ := stmt.(Import)
		printImport(v)
	// case Line:
	// 	v, _ := stmt.(Line)
	// 	printLine(v)
	case Warn:
		v, _ := stmt.(Warn)
		printWarn(v)
	case Sumdef:
		v, _ := stmt.(Sumdef)
		printSumdef(v)
	case Resdef:
		v, _ := stmt.(Resdef)
		printResdef(v)
	case Pext:
		v, _ := stmt.(Pext)
		printPext(v)
	case Error:
		v, _ := stmt.(Error)
		printError(v)
	case Undef:
		v, _ := stmt.(Undef)
		printUndef(v)
	case Ifdef:
		v, _ := stmt.(Ifdef)
		printIfdef(v)
	case Ifndef:
		v, _ := stmt.(Ifndef)
		printIfndef(v)
	case Label:
		v, _ := stmt.(Label)
		printLabel(v)
	case Macro:
		v, _ := stmt.(Macro)
		printMacro(v)
	case MacroCall:
		v, _ := stmt.(MacroCall)
		printMacroCall(v)
	case Add:
		v, _ := stmt.(Add)
		printAdd(v)
	case Mov:
		v, _ := stmt.(Mov)
		printMov(v)
	case Out:
		v, _ := stmt.(Out)
		printOut(v)
	case Jmp:
		v, _ := stmt.(Jmp)
		printJmp(v)
	case Jnc:
		v, _ := stmt.(Jnc)
		printJnc(v)
	}

}

func printDefine(define Define) {
	fmt.Printf("define_directive: (%s -> %s)\n", define.definition, define.name)
}

func printImport(imprt Import) {
	fmt.Printf("import_directive: (%s)\n", imprt.name)
}

// func printLine(line Line) {
// 	fmt.Printf("line_directive: from %s paste line %s\n", line.name, line.lineNumber)
// }

func printWarn(warn Warn) {
	fmt.Printf("warn_directive: %s\n", warn.message)
}

func printSumdef(sumdef Sumdef) {
	fmt.Printf("sumdef: %s = %s + %s\n", sumdef.def1, sumdef.def1, sumdef.def2)
}

func printResdef(resdef Resdef) {
	fmt.Printf("resdef: %s = %s - %s\n", resdef.def1, resdef.def1, resdef.def2)
}

func printPext(pext Pext) {
	fmt.Printf("pext: %s connects to %s\n", pext.pextName, pext.pextAddress)
}

func printError(err Error) {
	fmt.Printf("error: %s\n", err.message)
}

func printUndef(undef Undef) {
	fmt.Printf("undefined: %s\n", undef.definition)
}

func printIfdef(ifdef Ifdef) {
	fmt.Printf("if %s defined:\n", ifdef.definition)
	depth++
	printBlock(ifdef.bodyTrue)
	depth--
	if len(ifdef.bodyFalse.elements) != 0 {
		for i := 0; i < depth; i++ {
			fmt.Printf("\t\t\t")
		}
		fmt.Printf("else:\n")
		depth++
		printBlock(ifdef.bodyFalse)
		depth--
	}
}

func printIfndef(ifndef Ifndef) {
	fmt.Printf("if %s not defined:\n", ifndef.definition)
	depth++
	printBlock(ifndef.bodyTrue)
	depth--
	if len(ifndef.bodyFalse.elements) != 0 {
		for i := 0; i < depth; i++ {
			fmt.Printf("\t\t\t")
		}
		fmt.Printf("else:\n")
		depth++
		printBlock(ifndef.bodyFalse)
		depth--
	}
}
func printLabel(label Label) {
	fmt.Printf("Label: %s\n", label.name)
}

func printMacro(macro Macro) {
	fmt.Printf("macro %s: %s {\n", macro.macroName, macro.args)
	depth++
	printBlock(macro.body)
	depth--
	for i := 0; i < depth; i++ {
		fmt.Printf("\t\t\t")
	}
	fmt.Printf("}\n")
}

func printMacroCall(macrocall MacroCall) {
	fmt.Printf("call: %s %s\n", macrocall.macroName, macrocall.args)
}

func printAdd(add Add) {
	fmt.Printf("add %d, %s\n", add.reg, add.value)
}

func printMov(mov Mov) {
	fmt.Printf("mov %d, %d +%d\n", mov.reg1, mov.reg2, mov.fa)
}

func printOut(out Out) {
	fmt.Printf("out %d %s\n", out.reg, out.fa)
}

func printCmp(cmp Cmp) {
	fmt.Printf("cmp %d, %d, %d\n", cmp.regA, cmp.regB, cmp.operation)
}

func printJmp(jmp Jmp) {
	fmt.Printf("jmp %d|%d\n", jmp.regB, jmp.addr)
}

func printJnc(jnc Jnc) {
	fmt.Printf("jnc %d|%d\n", jnc.regB, jnc.addr)
}
