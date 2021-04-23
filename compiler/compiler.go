package compiler

import (
	"main/ivm"
	"strconv"
)

// ImpScope - The scope struct to keep track of locals
type ImpScope struct {
	Variables    map[string]int
	Instructions []uint16
	Calls        map[string]int
	Name         string
}

// NewScope - Creates a new scope struct
func NewScope() *ImpScope {
	return &ImpScope{
		Variables:    make(map[string]int),
		Instructions: make([]uint16, 0, 100),
	}
}

// Compiler - The main compiler
type Compiler struct {
	Scopes map[string]*ImpScope
	Ast    *AST
}

// SETL - SETL Instruction
func (compiler *Compiler) SETL(register string, value string) []uint16 {
	reg, err := strconv.Atoi(register)
	val, err2 := strconv.Atoi(value)

	if err != nil || err2 != nil {
		panic("Invalid instruction parameter.")
	}

	return []uint16{
		uint16(ivm.SETL),
		uint16(reg),
		uint16(val),
	}
}

// LOAD - LOAD Instruction
func (compiler *Compiler) LOAD(register string) []uint16 {
	reg, err := strconv.Atoi(register)

	if err != nil {
		panic("Invalid instruction parameter.")
	}

	return []uint16{
		uint16(ivm.SETL),
		uint16(reg),
	}
}

// STORE - STORE Instruction
func (compiler *Compiler) STORE(register string) []uint16 {
	reg, err := strconv.Atoi(register)

	if err != nil {
		panic("Invalid instruction parameter.")
	}

	return []uint16{
		uint16(ivm.STORE),
		uint16(reg),
	}
}

// CALL - CALL Instruction
func (compiler *Compiler) CALL(instructionnum string) []uint16 {
	ip, err := strconv.Atoi(instructionnum)

	if err != nil {
		panic("Invalid instruction parameter.")
	}

	return []uint16{
		uint16(ivm.CALL),
		uint16(ip),
	}
}

func (compiler *Compiler) cArgs(args []*AST, scope *ImpScope) []uint16 {
	instructions := make([]uint16, 0, 20)
	for i, arg := range args {
		instructions = append(instructions, compiler.STORE(string(i))...)
	}
}

func (compiler *Compiler) cString(expr *AST, scope *ImpScope) []uint16 {

}

func (compiler *Compiler) cVariable(expr *AST, scope *ImpScope) []uint16 {

}

func (compiler *Compiler) cBinary(expr *AST, scope *ImpScope) []uint16 {

}

func (compiler *Compiler) cAssign(expr *AST, scope *ImpScope) []uint16 {

}

func (compiler *Compiler) cSyscall(expr *AST, scope *ImpScope) []uint16 {
	funcname := expr.TokValue.GetValue()

	switch funcname {
	case PUTINT:

	}
}

func (compiler *Compiler) cFunctionCall(expr *AST, scope *ImpScope) []uint16 {
	funcname := expr.TokValue.GetValue()

	if funcname == "syscall" {
		return compiler.cSyscall(expr, scope)
	}

	return compiler.CALL("1")
}

func (compiler *Compiler) cScope(expr *AST, scope *ImpScope) []uint16 {
	iScope := NewScope()

	scopeName := expr.TokValue.GetValue()

	compiler.Scopes[scopeName] = iScope

	block := expr.Block

	for _, exp := range block {
		instructions := compiler.Compile(exp, iScope)

		if len(instructions) > 0 && instructions[0] == uint16(ivm.CALL) {
			iScope.Calls[exp.TokValue.GetValue()] = len(iScope.Instructions)
		}

		iScope.Instructions = append(iScope.Instructions, instructions...)
	}

	return make([]uint16, 0)
}

// Compile - Compile program to instructions
func (compiler *Compiler) Compile(ast *AST, scope *ImpScope) []uint16 {
	switch ast.ExprType {
	case Scope:
		return compiler.cScope(ast, scope)

	case String:
		return compiler.cString(ast, scope)

	case Variable:
		return compiler.cVariable(ast, scope)

	case Binary:
		return compiler.cBinary(ast, scope)

	case Assign:
		return compiler.cAssign(ast, scope)

	case FunctionCall:
		return compiler.cFunctionCall(ast, scope)

	default:
		panic("AHHHH")
	}
}

// LoadAST - Loads an ast into the compiler struct
func (compiler *Compiler) LoadAST(ast *AST) {
	compiler.Ast = ast
}

// NewCompiler - Creates a new compiler
func NewCompiler() *Compiler {
	return &Compiler{
		Scopes: make(map[string]*ImpScope),
	}
}
