package libpreproc

import "fmt"

//PrintProg - prints program AST
func PrintProg(prg Program) {
	fmt.Printf("program\n")
	for _, v := range prg.sections {
		printSection(v)
	}
}

func printSection(sec Section) {
	fmt.Printf("\tsection\n")
	fmt.Printf("\t\tsection_name: %s\n", sec.sectionName)
	fmt.Printf("\t\tsection_data:\n")
	printBlock(sec.sectionContent)
}

func printBlock(blk Block) {
	for _, v := range blk.elements {
		fmt.Printf("\t\t\t")
		estimateStmt(v)
	}
}

func estimateStmt(stmt interface{}) {
	switch stmt.(type) {
	case *Define:
		v, _ := stmt.(*Define)
		printDefine(v)
	case *Import:
		v, _ := stmt.(*Import)
		printImport(v)
	case *Line:
		v, _ := stmt.(*Line)
		printLine(v)
	case *Warn:
		v, _ := stmt.(*Warn)
		printWarn(v)
	}
}

func printDefine(define *Define) {
	fmt.Printf("define_directive: (%s -> %s)\n", define.definition, define.name)
}

func printImport(imprt *Import) {
	fmt.Printf("import_directive: (%s)\n", imprt.name)
}

func printLine(line *Line) {
	fmt.Printf("line_directive: from %s paste line %s\n", line.name, line.lineNumber)
}

func printWarn(warn *Warn) {
	fmt.Printf("warn_directive: %s\n", warn.message)
}
