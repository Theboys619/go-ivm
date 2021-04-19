package main

import (
	"fmt"
	"main/ivm"
)

func main() {
	fmt.Println("Hello, World!")
	vm := ivm.NewVM()
	vm.LoadProgram([]ivm.Instruction{
		ivm.SET, 0, 2,
		ivm.SET, 1, 3,
		ivm.ADD, 0, 1, 2,
		ivm.PUTINT, 2,
		ivm.HALT,
	})

	vm.Run(0)
}