package main

import (
	"fmt"
	"main/compiler"
	"main/ivm"
)

func main() {
	// fmt.Println("Hello, World!")
	vm := ivm.NewVM()
	vm.LoadProgram([]ivm.Instruction{
		ivm.SETL, 0, 2, // 2
		ivm.SETL, 1, 3, // 5
		ivm.STORE, 0, // 7
		ivm.STORE, 1, // 9
		ivm.CALL, 31, // 11
		ivm.LOAD, 3, // 13
		ivm.PUTINT, 3, // 15
		ivm.NEW, 4, 1, 0, // 19
		ivm.SETPROP, 4, 0, 6, 0, // 5
		ivm.PROP, 4, 5, 0, 0, // 5
		ivm.HALT,    // 20
		ivm.LOAD, 0, // 22
		ivm.LOAD, 1, // 24
		ivm.ADDL, 0, 1, 2, // 28
		ivm.CAST, 2, 0, // 31
		ivm.STORE, 2, // 33
		ivm.SEND, // 34
	})

	vm.Run(0)

	lexer := compiler.NewLexer()
	tokens := lexer.Tokenize("make x = \"Test Hello\";\n log(x);")
	fmt.Println(tokens)

	parser := compiler.NewParser()

	ast := parser.Parse(tokens)

	fmt.Println(ast.Block)
}
