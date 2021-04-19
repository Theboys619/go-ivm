package ivm

import "fmt"

// Instruction - Instruction type (uint16)
type Instruction uint16

// Instructions
const (
	HALT Instruction = iota
	ADD				// ADD - Add two registers
	ADDL			// ADDL - Add two registers into local register
	SET				// SET - Set Register
	SETL			// SETG - Set Local Register
	STORE			// STORE - Stores an register into argument list for calling functions
	LOAD			// LOAD - loads value into register from arglist
	PUTINT			// PUTINT - Print integer from register
	CALL			// CALL - Jumps to instruction
	SEND			// SEND - Returns from function
)

// Register for storing data
type Register struct {
	data uint16
}

// SetValue - Sets value of register
func (reg *Register) SetValue(val uint16) {
	reg.data = val
}

// GetValue - Get value of register
func (reg *Register) GetValue() uint16 {
	return reg.data
}

// Frame - Frame for locals
type Frame struct {
	locals map[int]Register
}

// NewFrame - Creates new Frame struct
func NewFrame() *Frame {
	return &Frame{
		locals: make(map[int]Register),
	}
}

// GetLocal - Get local from frame
func (frame *Frame) GetLocal(num int) Register {
	return frame.locals[num]
}

func (frame *Frame) SetLocal(regnum int, reg Register) {
	frame.locals[regnum] = reg
}

// SetLocals - Set local map
func (frame *Frame) SetLocals(locals map[int]Register) {
	frame.locals = locals
}

// VM for interpretation
type VM struct {
	program []Instruction
	registers []Register
	args []uint16
	frames []Frame

	ip int
	fp int
}

// NewVM - Creates a new Virtual Machine
func NewVM() *VM {
	vm := &VM{
		program: make([]Instruction, 0),
		registers: make([]Register, 50),
		args: make([]uint16, 0, 15),
		frames: make([]Frame, 1),
		ip: 0,
		fp: 0,
	}

	return vm
}

// SetRegisters - Sets all registers from frame
func (vm *VM) SetRegisters(frame Frame) {
	for regnum, reg := range frame.locals {
		vm.registers[regnum] = reg
	}
}

// SetRegister - Sets a register with a uint16 value
func (vm *VM) SetRegister(reg int, val uint16) {
	vm.registers[reg].SetValue(val)
}

// GetRegister - Gets a register
func (vm *VM) GetRegister(reg int) *Register {
	return &vm.registers[reg]
}

// GetArg - Gets a argument register
func (vm *VM) GetArg(arg int) uint16 {
	return vm.args[arg]
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

// GetFP - Gets the frame pointer
func (vm *VM) GetFP() int {
	return vm.fp
}

// SetFP - Sets the frame pointer
func (vm *VM) SetFP(num int) {
	vm.fp = num
}

// advanceFP - Increment frame pointer
func (vm *VM) advanceFP(adv... int) int {
	amount := 1

	if len(adv) > 0 {
		amount = adv[0]
	}

	vm.fp += amount
	return vm.fp
}

// GetFrame - Gets frame from fp
func (vm *VM) GetFrame(adv... int) Frame {
	amount := 0

	if len(adv) > 0 {
		amount = adv[0]
	}

	return vm.frames[vm.fp + amount]
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

		case ADDL:
			r1 := vm.GetRegister(int(vm.nextInstruction()))
			r2 := vm.GetRegister(int(vm.nextInstruction()))
			regnum3 := int(vm.nextInstruction())
			r3 := vm.GetRegister(regnum3)
			sum := r1.GetValue() + r2.GetValue()
			
			r3.SetValue(sum)
			// vm.locals = append(vm.locals, regnum3)
			frame := vm.GetFrame()
			frame.SetLocal(regnum3, *r3)
			vm.advanceIP()

		case SET:
			regnum := int(vm.nextInstruction())
			reg := vm.GetRegister(regnum)
			val := vm.GetInstruction(1)

			reg.SetValue(uint16(val))
			vm.advanceIP(2)

		case SETL:
			regnum := int(vm.nextInstruction())
			reg := vm.GetRegister(regnum)
			val := vm.GetInstruction(1)

			// vm.locals = append(vm.locals, regnum)
			frame := vm.GetFrame()
			frame.SetLocal(regnum, *reg)
			
			reg.SetValue(uint16(val))
			vm.advanceIP(2)
			
		case PUTINT:
			reg := vm.GetRegister(int(vm.nextInstruction()))
			fmt.Print(reg.GetValue())
			vm.advanceIP()

		case STORE:
			reg := vm.GetRegister(int(vm.nextInstruction()))
			vm.args = append(vm.args, reg.GetValue())
			vm.advanceIP()

		case LOAD:
			reg := vm.GetRegister(int(vm.nextInstruction()))
			item := vm.args[0]
			vm.args = vm.args[1:]

			reg.SetValue(item)
			vm.advanceIP()

		case CALL:
			instr := int(vm.nextInstruction())
			returninstr := uint16(vm.advanceIP())

			vm.args = append(vm.args, returninstr)
			vm.advanceFP()
			vm.SetIP(instr)

		case SEND:
			returninstr := vm.args[0]
			vm.args = vm.args[1:]

			vm.SetFP(vm.fp-1)
			
			vm.SetRegisters(vm.GetFrame())

			vm.SetIP(int(returninstr))

		default:
			panic("Illegal Instruction")
		}
	}

	fmt.Println()
	fmt.Println(vm.registers)
	fmt.Println()

	fmt.Println("Program exited with exit code", exitCode)
}