package main

import (
	"bufio"
	"fmt"
	"os"
	p "preprocessor/libpreproc"
)

func main() {
	filename := os.Args[1]
	file, _ := os.Open(filename)
	fmt.Printf("Parsing %s...\n\n", filename)
	stmt, _ := p.NewParser(bufio.NewReader(file)).ParseFile()
	body := stmt.Body
	fmt.Printf("Printing AST:\n")
	for _, a := range body {
		p.PrintNode(a)
	}
}
