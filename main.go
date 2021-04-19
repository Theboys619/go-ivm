package main

import (
	"main/ivm"
)

func main() {
	// fmt.Println("Hello, World!")
	vm := ivm.NewVM()
	vm.LoadProgram([]ivm.Instruction{
		ivm.SETL, 0, 2,
		ivm.SETL, 1, 3,
		ivm.STORE, 0,
		ivm.STORE, 1,
		ivm.CALL, 17,
		ivm.LOAD, 3,
		ivm.PUTINT, 3,
		ivm.HALT,
		ivm.LOAD, 0,
		ivm.LOAD, 1,
		ivm.ADDL, 0, 1, 2,
		ivm.STORE, 2,
		ivm.SEND,
	})

	vm.Run(0)
}