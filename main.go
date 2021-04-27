package main

import (
	// "fmt"
	// "main/compiler"
	"main/ivm"
)

func main() {
	// fmt.Println("Hello, World!")
	vm := ivm.NewVM(256)
	vm.LoadProgram([]uint8{
		ivm.PUSHL, 0x0002,
		ivm.PUSHL, 0x0002,
		ivm.ADD,
		ivm.HALT,
	})
	vm.Run()

	// lexer := compiler.NewLexer()
	// tokens := lexer.Tokenize("make x = \"Test Hello\";\n log(x);")
	// fmt.Println(tokens)

	// parser := compiler.NewParser()

	// ast := parser.Parse(tokens)

	// fmt.Println(ast.Block)
}
