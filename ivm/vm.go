package ivm

import "fmt"

// Instruction - Instruction type (uint16)
type Instruction uint16
// CastTypes - A const enum for cast types
type CastTypes uint16

// Instructions
const (
	HALT Instruction = iota
	ADD					// ADD - Add two registers (ADD [REG] [REG] [REG])
	ADDL				// ADDL - Add two registers into local register (ADDL [REG] [REG] [REG])
	SET					// SET - Set Register (SET [REG] [VAL])
	SETL				// SETG - Set Local Register (SETL [REG] [VAL])
	STORE				// STORE - Stores an register into argument list for calling functions (STORE [REG])
	LOAD				// LOAD - loads value into register from arglist (LOAD [REG])
	PUTINT				// PUTINT - Print integer from register (PUTINT [REG])
	SYS					// SYS - A Syscall to an internal function to the VM (SYS [INT])
	CALL				// CALL - Jumps to instruction (CALL [IP])
	SEND				// SEND - Returns from function (SEND)
	CAST				// CAST - Casts to a type (CAST [REG] [TYPE])
	NEW					// NEW - Creates a new object on the heap (NEW [REG] [PROPSIZE] [METHODSIZE])
	SETPROP				// SETPROP - Sets a prop in object pointer / reg (SETPROP [REG] [PROP] [VALUE] [TYPE])
	PROP				// PROP - Get a prop from object pointer / reg (PROP [REG] [REG2] [PROP] [TYPE])
	NUMINSTRUCTIONS		// NUMINSTRUCTIONS - Number of instructions
)

// CastTypes
const (
	INT CastTypes = iota
	CHAR
	NUMTYPES
)

// Register for storing data
type Register struct {
	data interface{}
}

// SetValue - Sets value of register
func (reg *Register) SetValue(val interface{}) {
	reg.data = val
}

// GetValue - Get value of register
func (reg *Register) GetValue() interface{} {
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

// SetLocal - Set local entry in map
func (frame *Frame) SetLocal(regnum int, reg Register) {
	frame.locals[regnum] = reg
}

// SetLocals - Set local map
func (frame *Frame) SetLocals(locals map[int]Register) {
	frame.locals = locals
}

// CastValue - Cast any to a type listen in CastTypes
func CastValue(value *interface{}, casttype uint16) {
	switch CastTypes(casttype) {
	case INT:
		*value = int((*value).(uint16))
	default:
		return
	}
}

// Object - An objects for data
type Object struct {
	properties []interface{}
	methods []int
}

// SetProperty - Sets a property in the object
func (object *Object) SetProperty(prop int, dtype uint16, value interface{}) {
	if prop >= len(object.properties) {
		object.properties = append(object.properties, nil)
	}
	object.properties[prop] = value
}

// GetProperty - Gets a property in the object
func (object *Object) GetProperty(prop int, dtype uint16, value *interface{}) {
	property := object.properties[prop]

	CastValue(&property, dtype)
	*value = property
}

// Heap - For all allocated objects / structs
type Heap struct {
	objects []Object
	refs int
}

// NewHeap - Creates a new Heap struct
func NewHeap(objects int) *Heap {
	return &Heap{
		objects: make([]Object, objects),
	}
}

// NewObject - Creates a new object on the heap
func (heap *Heap) NewObject(propsize int, methodsize int) *Object {
	object := Object{
		properties: make([]interface{}, propsize, propsize + 5),
		methods: make([]int, methodsize, methodsize + 5),
	}

	heap.refs++

	return &object
}

// DestroyObject - Remove object off the heap
func (heap *Heap) DestroyObject(object *Object) {
	object.properties = make([]interface{}, 0)
	object.properties = nil
	object = nil
	heap.refs--
}

// VM for interpretation
type VM struct {
	program []Instruction
	registers []Register

	args []interface{}
	frames []*Frame
	heap *Heap

	ip int
	fp int
}

// NewVM - Creates a new Virtual Machine
func NewVM() *VM {
	vm := &VM{
		program: make([]Instruction, 0),
		registers: make([]Register, 50),
		args: make([]interface{}, 0, 15),
		frames: make([]*Frame, 1),
		heap: NewHeap(200),
		ip: 0,
		fp: 0,
	}

	vm.frames[0] = NewFrame()

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
func (vm *VM) GetArg(arg int) interface{} {
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

	if vm.frames[num] == nil {
		vm.frames = append(vm.frames, NewFrame())
	}
}

// advanceFP - Increment frame pointer
func (vm *VM) advanceFP(adv... int) int {
	amount := 1

	if len(adv) > 0 {
		amount = adv[0]
	}

	vm.fp += amount

	if len(vm.frames) - 1 < vm.fp {
		vm.frames = append(vm.frames, NewFrame())
	}

	return vm.fp
}

// GetFrame - Gets frame from fp
func (vm *VM) GetFrame(adv... int) *Frame {
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
			sum := r1.GetValue().(uint16) + r2.GetValue().(uint16)
			
			r3.SetValue(sum)
			vm.advanceIP()

		case ADDL:
			r1 := vm.GetRegister(int(vm.nextInstruction()))
			r2 := vm.GetRegister(int(vm.nextInstruction()))
			regnum3 := int(vm.nextInstruction())
			r3 := vm.GetRegister(regnum3)
			sum := r1.GetValue().(uint16) + r2.GetValue().(uint16)
			
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
			
			reg.SetValue(uint16(val))
			frame.SetLocal(regnum, *reg)
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
			vm.GetFrame().SetLocal(vm.ip, *reg)
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
			
			vm.SetRegisters(*vm.GetFrame())

			vm.SetIP(int(returninstr.(uint16)))

		case CAST:
			reg := vm.GetRegister(int(vm.nextInstruction()))
			casttype := uint16(vm.GetInstruction(1))
			vm.advanceIP(2)

			if casttype >= uint16(NUMTYPES) {
				panic("Illegal Instruction. Invalid casting type")
			}

			value := reg.GetValue()
			CastValue(&value, casttype)

			reg.SetValue(value)

		case NEW:
			reg := vm.GetRegister(int(vm.nextInstruction()))
			propsize := uint16(vm.nextInstruction())
			methodsize := uint16(vm.nextInstruction())
			vm.advanceIP()

			objectptr := vm.heap.NewObject(int(propsize), int(methodsize))
			reg.SetValue(objectptr)

		case SETPROP:
			reg := vm.GetRegister(int(vm.nextInstruction()))
			prop := uint16(vm.nextInstruction())
			value := uint16(vm.nextInstruction())
			proptype := uint16(vm.nextInstruction())
			vm.advanceIP()

			objectptr := reg.GetValue().(*Object)
			
			objectptr.SetProperty(int(prop), proptype, value)

		case PROP:
			reg := vm.GetRegister(int(vm.nextInstruction()))
			reg2 := vm.GetRegister(int(vm.nextInstruction()))
			prop := uint16(vm.nextInstruction())
			proptype := uint16(vm.nextInstruction())
			vm.advanceIP()

			value := reg2.GetValue()
			objectptr := reg.GetValue().(*Object)

			objectptr.GetProperty(int(prop), proptype, &value)
			reg2.SetValue(value)
			vm.GetFrame().SetLocal(vm.ip, *reg2)

		default:
			panic("Illegal Instruction")
		}
	}

	fmt.Println()
	fmt.Println(vm.registers)
	fmt.Println()

	fmt.Println("Program exited with exit code", exitCode)
}