package main

import (
	"fmt"
	"main/ivm"
)

func main() {
	fmt.Println("Hello, World!")
	vm := ivm.NewVM()
	vm.LoadProgram([]ivm.Instruction{
		ivm.SETL, 0, 2,
		ivm.SETL, 1, 3,
		ivm.ADDL, 0, 1, 2,
		ivm.PUTINT, 2,
		ivm.HALT,
		ivm.ADD, 0,
	})

	vm.Run(0)
}