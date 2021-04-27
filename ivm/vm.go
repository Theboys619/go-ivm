package ivm

import "fmt"

// Instructions
const (
	HALT     = 0x10
	PUSHL    = 0x11 // PUSH Literal onto stack
	PUSHR    = 0x12 // PUSH Register address onto stack
	PUSHREGV = 0x13 // PUSH Register value onto stack
	POP      = 0x14
	ADD      = 0x15 // ADD From stack
	ADDRL    = 0x16 // ADD Register and Literal
	ADDRR    = 0x17 // ADD Register and Register
	MOVL     = 0x18 // MOV Literal into a register
	CALL     = 0x19 // CALL jump to address and save registers
	CALLREG     = 0x1a // CALL jump to address from register and save registers
	RET      = 0x1b // RET from call
)

// StackT - The stack struct
type StackT struct {
	Data []uint16
	Size int
	FrameSize int
}

// NewStack - Creates a Stack struct
func NewStack(size int) *StackT {
	return &StackT{
		Data: make([]uint16, size/2+(size/4), size),
		Size: size,
		FrameSize: 0,
	}
}

// GetValue - Gets a value in the stack at an address (index)
func (stack *StackT) GetValue(address uint16) uint16 {
	return stack.Data[address]
}

// GetValue8 - Gets a value in the stack at an address (index)
func (stack *StackT) GetValue8(address uint8) uint8 {
	return uint8(stack.Data[address] & 0xff00)
}

// Memory - Data for memory
type Memory struct {
	Data        []uint8
	Registers   map[string]uint16
	RegisterMap map[int]string
	Size        int
}

// NewMemory - Creates a Memory struct
func NewMemory(size int) *Memory {
	return &Memory{
		Data:        make([]uint8, size, size+1),
		Registers:   make(map[string]uint16),
		RegisterMap: make(map[int]string),
		Size:        size,
	}
}

// LoadMem - Loads memory into the data property
func (mem *Memory) LoadMem(data []uint8) {
	mem.Data = data
	mem.Data = append(mem.Data, make([]uint8, mem.Size-len(data), mem.Size)...)
}

// GetValue8 - Gets the value at an adress (uint8)
func (mem *Memory) GetValue8(address uint16) uint8 {
	return mem.Data[address]
}

// GetValue16 - Gets the value at an address (uint16)
func (mem *Memory) GetValue16(address uint16) uint16 {
	return uint16((mem.Data[address] << 8) | mem.Data[address+1])
}

// GetRegister - Get a register address
func (mem *Memory) GetRegister(name string) uint16 {
	return mem.Registers[name]
}

// GetRegisterVal8 - Gets the value the register is pointing to. (uint8)
func (mem *Memory) GetRegisterVal8(name string) uint8 {
	return mem.GetValue8(mem.Registers[name])
}

// GetRegisterVal16 - Gets the value the register is pointing to. (uint16)
func (mem *Memory) GetRegisterVal16(name string) uint16 {
	return mem.GetValue16(mem.Registers[name])
}

// SetRegister - Sets the registers address
func (mem *Memory) SetRegister(name string, address uint16) {
	mem.Registers[name] = address
}

// VM - The Virtual Machine
type VM struct {
	Stack *StackT
	Mem   *Memory
}

// NewVM - Creates a new Virtual Machine
func NewVM(size int) *VM {
	mem := NewMemory(size * size)

	mem.SetRegister("ip", 0x0000)
	mem.SetRegister("sp", 0x0000)
	mem.SetRegister("bp", 0x0000)
	mem.SetRegister("r1", 0x0000)
	mem.SetRegister("r2", 0x0000)
	mem.SetRegister("r3", 0x0000)
	mem.SetRegister("r4", 0x0000)
	mem.SetRegister("r5", 0x0000)
	mem.SetRegister("r6", 0x0000)
	mem.SetRegister("r7", 0x0000)
	mem.SetRegister("r8", 0x0000)

	mem.RegisterMap[0] = "ip"
	mem.RegisterMap[1] = "sp"
	mem.RegisterMap[2] = "bp"
	mem.RegisterMap[3] = "r1"
	mem.RegisterMap[4] = "r2"
	mem.RegisterMap[5] = "r3"
	mem.RegisterMap[6] = "r4"
	mem.RegisterMap[7] = "r5"
	mem.RegisterMap[8] = "r6"
	mem.RegisterMap[9] = "r7"
	mem.RegisterMap[10] = "r8"

	return &VM{
		Stack: NewStack(size * 2),
		Mem:   mem,
	}
}

func (vm *VM) LoadProgram(prgm []uint8) {
	vm.Mem.LoadMem(prgm)
}

func (vm *VM) Push(val uint16) {
	sp := vm.Mem.GetRegister("sp")
	sp--
	vm.Stack.Data[sp] = val
	vm.Mem.SetRegister("sp", sp)

	vm.Stack.FrameSize++
}

func (vm *VM) Pop() uint16 {
	sp := vm.Mem.GetRegister("sp")
	val := vm.Stack.Data[sp]
	vm.Mem.SetRegister("sp", sp+1)

	vm.Stack.FrameSize--

	return val
}

// NextInstruction - Gets the next instruction
func (vm *VM) NextInstruction() uint8 {
	ip := vm.Mem.GetRegister("ip")
	ip++

	vm.Mem.SetRegister("ip", ip)

	return vm.Mem.GetValue8(ip)
}

// NextInstruction16 - Gets the next instruction
func (vm *VM) NextInstruction16() uint16 {
	ip := vm.Mem.GetRegister("ip")
	ip++

	vm.Mem.SetRegister("ip", ip)

	return vm.Mem.GetValue16(ip)
}

// Fetch - Fetches the current instruction
func (vm *VM) Fetch() uint8 {
	ip := vm.Mem.GetRegister("ip")

	return vm.Mem.GetValue8(ip)
}

// Fetch16 - Fetches the current instruction as uint16
func (vm *VM) Fetch16() uint16 {
	ip := vm.Mem.GetRegister("ip")

	return vm.Mem.GetValue16(ip)
}

func (vm *VM) pushFrame() {
	vm.Push(vm.Mem.GetRegister("r1"))
	vm.Push(vm.Mem.GetRegister("r2"))
	vm.Push(vm.Mem.GetRegister("r3"))
	vm.Push(vm.Mem.GetRegister("r4"))
	vm.Push(vm.Mem.GetRegister("r5"))
	vm.Push(vm.Mem.GetRegister("r6"))
	vm.Push(vm.Mem.GetRegister("r7"))
	vm.Push(vm.Mem.GetRegister("r8"))
	vm.Push(vm.Mem.GetRegister("ip"))
	vm.Push(vm.Mem.GetRegister("bp"))
	vm.Mem.SetRegister("bp", vm.Mem.GetRegister("sp"))
	vm.Stack.FrameSize = 0
}

func (vm *VM) popFrame() {
	vm.Mem.SetRegister("sp", vm.Mem.GetRegister("bp"))
	vm.Mem.SetRegister("bp", vm.Pop())
	vm.Mem.SetRegister("ip", vm.Pop())
	vm.Mem.SetRegister("r8", vm.Pop())
	vm.Mem.SetRegister("r7", vm.Pop())
	vm.Mem.SetRegister("r6", vm.Pop())
	vm.Mem.SetRegister("r5", vm.Pop())
	vm.Mem.SetRegister("r4", vm.Pop())
	vm.Mem.SetRegister("r3", vm.Pop())
	vm.Mem.SetRegister("r2", vm.Pop())
	vm.Mem.SetRegister("r1", vm.Pop())

	nArgs := vm.Pop()

	var i uint16;
	for i = 0; i < nArgs; i++ {
		vm.Pop()
	}
}

func (vm *VM) pushLiteral() {
	value1 := vm.NextInstruction()
	value2 := vm.NextInstruction()
	vm.Push(uint16((value1 << 8) | value2))
	vm.NextInstruction()
}

func (vm *VM) pushRegister() {
	vm.NextInstruction()
	regnum := vm.Fetch()
	registerName, ok := vm.Mem.RegisterMap[int(regnum)]
	if !ok {
		panic("Register " + registerName + " was not resolved")
	}

	vm.Push(vm.Mem.GetRegister(registerName))

	vm.NextInstruction()
}

// If register contains a pointer
func (vm *VM) pushRegValue() {
	vm.NextInstruction()
	regnum := vm.Fetch()
	registerName, ok := vm.Mem.RegisterMap[int(regnum)]
	if !ok {
		panic("Register " + registerName + " was not resolved")
	}

	vm.Push(vm.Mem.GetRegisterVal16(registerName))
	vm.NextInstruction()
}

func (vm *VM) iADD() {
	b := vm.Pop()
	a := vm.Pop()

	vm.Push(a + b)
	vm.NextInstruction()
}

func (vm *VM) iPOP() {
	x := vm.Pop()

	regnum := vm.NextInstruction()
	registerName, ok := vm.Mem.RegisterMap[int(regnum)]
	if !ok {
		panic("Register " + registerName + " was not resolved")
	}

	vm.Mem.SetRegister(registerName, x)

	vm.NextInstruction()
}

func (vm *VM) iMOVL() {
	regnum := vm.NextInstruction()
	regName, ok := vm.Mem.RegisterMap[int(regnum)]

	if !ok {
		panic("Register " + regName + " was not resolved")
	}

	literal := vm.NextInstruction16()
	vm.Mem.SetRegister(regName, literal)
	vm.NextInstruction()
	vm.NextInstruction()
}

func (vm *VM) iCALL() {
	address := vm.NextInstruction16()
	vm.NextInstruction()
	vm.NextInstruction()

	vm.pushFrame()
	vm.Mem.SetRegister("ip", address)
}

func (vm *VM) iCALLREG() {
	regnum := vm.NextInstruction()
	regName, ok := vm.Mem.RegisterMap[int(regnum)]

	if !ok {
		panic("Register " + regName + " was not resolved")
	}

	vm.pushFrame()
	vm.Mem.SetRegister("ip", vm.Mem.GetRegister(regName))
}

// Debug - Get values
func (vm *VM) Debug() {
	fmt.Println("Memory:")
	fmt.Println()
	fmt.Println(vm.Mem.Data[:100])
	fmt.Println()
	fmt.Println("Stack:")
	fmt.Println(vm.Stack.Data)
	fmt.Println()

	fmt.Println("Registers:")
	for key, val := range vm.Mem.Registers {
		fmt.Printf("%s: %#04x\n", key, val)
	}

	fmt.Println()
	fmt.Println("Current Instruction:")
	fmt.Println(vm.Fetch())
	fmt.Println()
}

func (vm *VM) Run() {
	vm.Mem.SetRegister("sp", uint16(len(vm.Stack.Data)))
	vm.Mem.SetRegister("bp", uint16(len(vm.Stack.Data)-1))

	halted := false
	for !halted {
		vm.Debug()
		switch vm.Fetch() {
		case HALT:
			halted = true
			break

		case PUSHL: // [1byte] [2bytes] - 3 bytes
			vm.pushLiteral()

		case PUSHR: // [1byte] [1byte] - 2 bytes
			vm.pushRegister()

		case PUSHREGV: // [1byte] [1byte] - 2 bytes
			vm.pushRegValue()

		case POP: // [1byte] [1byte] - 2 bytes
			vm.iPOP()

		case ADD: // [1byte] - 1 byte
			vm.iADD()

		case MOVL: // [1byte] [1byte] [2byte] - 4 bytes
			vm.iMOVL()

		case CALL:
			vm.iCALL()

		case CALLREG:
			vm.iCALLREG()

		case RET:
			vm.popFrame()

		default:
			halted = true
			panic("Illegal instruction")
		}
	}
}
