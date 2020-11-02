package main

import (
	"bufio"
	"fmt"
	"os"
	p "preprocessor/libpreproc"
)

func main() {
	filename := os.Args[1]
	fmt.Printf("Parsing %s...\n\n", filename)
	file, _ := os.Open(filename)
	stmt, err := p.NewParser(bufio.NewReader(file)).ParseFile()
	if err != nil {
		fmt.Printf(err.Error())
	} else {
		p.PrintProg(stmt)
	}
	ld := os.Args[2]
	f, _ := p.OpenLinkerScript(ld)
	pt, _ := f.GetPartitionList()
	for _, val := range pt {
		println(val)
	}
}
