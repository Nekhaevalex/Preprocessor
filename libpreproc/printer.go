package libpreproc

import "fmt"

var depth int = 0

func printWS() {
	i := 0
	for i < depth {
		fmt.Printf(".")
		i++
	}
}

//PrintNode - print tree
func PrintNode(line interface{}) {
	printWS()
	switch v := line.(type) {
	case *ImportStmt:
		fmt.Printf("Import: %s\n", v.ImportFile)
	case *DefStmt:
		fmt.Printf("Definition: %s <- %s\n", v.DefinitionName, v.ToDef)
	case *PextStmt:
		fmt.Printf("Pext: %s at %s\n", v.PextName, v.PextAddress)
	case *ErrorStmt:
		fmt.Printf("Error: %s\n", v.ErrorMsg)
	case *PragmaStmt:
		fmt.Printf("Pragma: %s\n", v.PragmaName)
	case *LineStmt:
		fmt.Printf("Line %s from %s\n", v.LineNumber, v.FileName)
	case *MsgStmt:
		fmt.Printf("Message: %s\n", v.Msg)
	case *IfStmt:
		fmt.Printf("%t if on %s\n", v.Negative, v.DefName)
		depth++
		for _, b1 := range v.Branch1body {
			PrintNode(b1)
		}
		depth--
		if v.Branch2body != nil {
			fmt.Printf("Else\n")
			depth++
			for _, b2 := range v.Branch2body {
				PrintNode(b2)
			}
			depth--
		}
		fmt.Printf("End if\n")
	case *SumStmt:
		fmt.Printf("%s += %s\n", v.X, v.Y)
	case *ResStmt:
		fmt.Printf("%s -= %s\n", v.X, v.Y)
	case *UndefStmt:
		fmt.Printf("Undef: %s\n", v.DefName)
	case *ReturnStmt:
		fmt.Printf("Return: %s\n", v.ReturnName)
	case *MacroStmt:
		fmt.Printf("Macro Decl: %s\nVars:\n", v.Name)
		for num, vars := range v.Vars {
			fmt.Printf("%d\t%s\n", num, vars)
		}
		depth++
		for _, b1 := range v.Body {
			PrintNode(b1)
		}
		depth--
		fmt.Printf("Macro end\n")
	}
}
