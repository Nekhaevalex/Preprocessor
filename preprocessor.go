package main

import (
	"bufio"
	"fmt"
	"os"
	p "preprocessor/preprocessor"
)

func printNode(line interface{}) {
	switch v := line.(type) {
	case *p.ImportStmt:
		fmt.Printf("Import: %s\n", v.ImportFile)
	case *p.DefStmt:
		fmt.Printf("Definition: %s <- %s\n", v.DefinitionName, v.ToDef)
	case *p.PextStmt:
		fmt.Printf("Pext: %s at %s\n", v.PextName, v.PextAddress)
	case *p.ErrorStmt:
		fmt.Printf("Error: %s\n", v.ErrorMsg)
	case *p.PragmaStmt:
		fmt.Printf("Pragma: %s\n", v.PragmaName)
	case *p.LineStmt:
		fmt.Printf("Line %s from %s\n", v.LineNumber, v.FileName)
	case *p.MsgStmt:
		fmt.Printf("Message: %s\n", v.Msg)
	case *p.IfStmt:
		fmt.Printf("%t if on %s\n", v.Negative, v.DefName)
		for _, b1 := range v.Branch1body {
			printNode(b1)
		}
		if v.Branch2body != nil {
			fmt.Printf("Else\n")
			for _, b2 := range v.Branch2body {
				printNode(b2)
			}
		}
		fmt.Printf("End if\n")
	case *p.SumStmt:
		fmt.Printf("%s += %s\n", v.X, v.Y)
	case *p.ResStmt:
		fmt.Printf("%s -= %s\n", v.X, v.Y)
	case *p.UndefStmt:
		fmt.Printf("Undef: %s\n", v.DefName)
	case *p.ReturnStmt:
		fmt.Printf("Return: %s\n", v.ReturnName)
	case *p.MacroStmt:
		fmt.Printf("Macro Decl: %s\nVars:\n", v.Name)
		for num, vars := range v.Vars {
			fmt.Printf("%d\t%s\n", num, vars)
		}
		for _, b1 := range v.Body {
			printNode(b1)
		}
		fmt.Printf("Macro end\n")
	}
}

func main() {
	filename := os.Args[1]
	file, _ := os.Open(filename)
	fmt.Printf("Parsing %s ...\n", filename)
	stmt, _ := p.NewParser(bufio.NewReader(file)).ParseFile()
	body := stmt.Body
	for _, a := range body {
		printNode(a)
	}
}
