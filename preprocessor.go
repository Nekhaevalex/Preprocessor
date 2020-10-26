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
	stmt, err := p.NewParser(bufio.NewReader(file)).ParseFile()
	if err != nil {
		fmt.Printf(err.Error())
	} else {
		p.PrintProg(stmt)
	}
}
