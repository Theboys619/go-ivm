package main

import (
	// "fmt"
	// "main/compiler"
	// "fmt"
	"main/ivm"
)

func main() {
	// fmt.Println("Hello, World!")
	vm := ivm.NewVM(256)
	prog := make([]uint8, 0, 256)
	
	prog = append(prog, []uint8{
		ivm.PUSHL, 0x00, 0x02,
		ivm.PUSHL, 0x00, 0x02,
		ivm.ADD,
		ivm.POP, 0x03,
		ivm.PUSHR, 0x00,
		ivm.PUSHREGV, 0x00,
		ivm.MOVL, 0x04, 0x00, 0x06,
		ivm.MOVL, 0x05, 0x00, 0x0a,
		ivm.PUSHL, 0x00, 0x07,
		ivm.PUSHL, 0x00, 0x01,
		ivm.CALL, 0x00, 0x64,
		ivm.PUSHL, 0x00, 0x0c,
		ivm.HALT,
	}...)

	subroutine1 := []uint8{
		ivm.PUSHL, 0x00, 0x03,
		ivm.PUSHL, 0x00, 0x03,
		ivm.ADD,
		ivm.POP, 0x0a,
		ivm.RET,
	}

	prog = append(prog, prog[len(prog):100]...)
	prog = append(prog, subroutine1...)

	vm.LoadProgram(prog)
	vm.Run()

	// lexer := compiler.NewLexer()
	// tokens := lexer.Tokenize("make x = \"Test Hello\";\n log(x);")
	// fmt.Println(tokens)

	// parser := compiler.NewParser()

	// ast := parser.Parse(tokens)

	// fmt.Println(ast.Block)
}
