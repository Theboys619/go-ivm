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
		ivm.PUSHL, 0x00, 0x02,
		ivm.PUSHL, 0x00, 0x02,
		ivm.ADD,
		ivm.POP, 0x03,
		ivm.PUSHR, 0x00,
		ivm.PUSHREGV, 0x00,
		ivm.MOVL, 0x04, 0x00, 0x06,
		ivm.PUSHL, 0x00, 0x00,
		ivm.CALL, 0x00, 0x18,
		ivm.HALT,
		ivm.PUSHL, 0x00, 0x03,
		ivm.PUSHL, 0x00, 0x03,
		ivm.ADD,
		ivm.POP, 0x05,
		ivm.RET,
	})
	vm.Run()

	// lexer := compiler.NewLexer()
	// tokens := lexer.Tokenize("make x = \"Test Hello\";\n log(x);")
	// fmt.Println(tokens)

	// parser := compiler.NewParser()

	// ast := parser.Parse(tokens)

	// fmt.Println(ast.Block)
}
