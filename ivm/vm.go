package ivm

import "fmt"

// Instruction - Instruction type (uint16)
type Instruction uint16

// Instructions
const (
	HALT Instruction = iota
	ADD
	SET
	PUTINT
)

// Register for storing data
type Register struct {
	data uint16
}

// SetValue - Sets value of register
func (reg *Register) SetValue(val uint16) {
	reg.data = val;
}

// GetValue - Get value of register
func (reg *Register) GetValue() uint16 {
	return reg.data
}

// VM for interpretation
type VM struct {
	program []Instruction
	registers []Register
	ip int
}

// NewVM - Creates a new Virtual Machine
func NewVM() *VM {
	vm := &VM{
		program: make([]Instruction, 0),
		registers: make([]Register, 50),
		ip: 0,
	}

	return vm
}

// SetRegister - Sets a register with a uint16 value
func (vm *VM) SetRegister(reg int, val uint16) {
	vm.registers[reg].SetValue(val)
}

// GetRegister - Gets a register
func (vm *VM) GetRegister(reg int) *Register {
	return &vm.registers[reg]
}

// LoadProgram - Load bytecode into virutal machine
func (vm *VM) LoadProgram(code []Instruction) {
	vm.program = code;
}

// GetInstruction - Gets current or +amount instruction
func (vm *VM) GetInstruction(adv... int) Instruction {
	amount := 0

	if len(adv) > 0 {
		amount = adv[0]
	}
	
	return vm.program[vm.ip + amount]
}

// GetIP - Gets instruction pointer
func (vm *VM) GetIP() int {
	return vm.ip
}

// SetIP - Sets the instruction pointer
func (vm *VM) SetIP(val int) {
	vm.ip = val
}

// advanceIP - Increment instruction pointer
func (vm *VM) advanceIP(adv... int) int {
	amount := 1

	if len(adv) > 0 {
		amount = adv[0]
	}

	vm.ip += amount
	return vm.ip
}

// nextInstruction - Advance instruction and return current instruction
func (vm *VM) nextInstruction(adv... int) Instruction {
	vm.advanceIP(adv...)
	return vm.GetInstruction()
}

// Run - Runs the current program at ip
func (vm *VM) Run(ip int) {
	exitCode := -1
	exited := false
	vm.SetIP(ip)
	
	for !exited {
		switch vm.GetInstruction() {
		case HALT:
			exitCode = 0
			exited = true
			break

		case ADD:
			r1 := vm.GetRegister(int(vm.nextInstruction()))
			r2 := vm.GetRegister(int(vm.nextInstruction()))
			r3 := vm.GetRegister(int(vm.nextInstruction()))
			sum := r1.GetValue() + r2.GetValue()
			
			r3.SetValue(sum)
			vm.advanceIP()

		case SET:
			reg := vm.GetRegister(int(vm.nextInstruction()))
			val := vm.GetInstruction(1)
			
			reg.SetValue(uint16(val))
			vm.advanceIP(2)
			
		case PUTINT:
			reg := vm.GetRegister(int(vm.nextInstruction()))
			fmt.Print(reg.GetValue())
			vm.advanceIP()

		default:
			panic("Illegal Instruction")
		}
	}

	fmt.Println()
	fmt.Println(vm.registers)
	fmt.Println()

	fmt.Println("Program exited with exit code", exitCode)
}