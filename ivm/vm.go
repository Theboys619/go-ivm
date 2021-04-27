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
)

// StackT - The stack struct
type StackT struct {
	Data []uint16
	Size int
}

// NewStack - Creates a Stack struct
func NewStack(size int) *StackT {
	return &StackT{
		Data: make([]uint16, size/2+(size/4), size),
		Size: size,
	}
}

// Push - Push data onto stack
func (stack *StackT) Push(item uint16) {
	stack.Data = append(stack.Data, item)
}

// Pop - Pop item of of stack
func (stack *StackT) Pop() uint16 {
	stackLen := len(stack.Data)
	item := stack.Data[stackLen-1]
	stack.Data = stack.Data[:stackLen-1]

	return item
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
func (mem *Memory) GetValue8(address uint8) uint8 {
	return mem.Data[address]
}

// GetValue16 - Gets the value at an address (uint16)
func (mem *Memory) GetValue16(address uint16) uint16 {
	return uint16((mem.Data[address&0xff00] << 8) | mem.Data[address&0x00ff])
}

// GetRegister - Get a register address
func (mem *Memory) GetRegister(name string) uint16 {
	return mem.Registers[name]
}

// GetRegisterVal8 - Gets the value the register is pointing to. (uint8)
func (mem *Memory) GetRegisterVal8(name string) uint8 {
	return mem.GetValue8(uint8(mem.Registers[name] & 0x00ff))
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
}

func (vm *VM) Pop() uint16 {
	sp := vm.Mem.GetRegister("sp")
	val := vm.Stack.Data[sp]
	vm.Mem.SetRegister("sp", sp+1)

	return val
}

// NextInstruction - Gets the next instruction
func (vm *VM) NextInstruction() uint8 {
	ip := vm.Mem.GetRegister("ip")
	ip++

	vm.Mem.SetRegister("ip", ip)

	return vm.Mem.GetValue8(uint8(ip & 0x00ff))
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

	return vm.Mem.GetValue8(uint8(ip & 0x00ff))
}

// Fetch16 - Fetches the current instruction as uint17
func (vm *VM) Fetch16() uint16 {
	ip := vm.Mem.GetRegister("ip")

	return vm.Mem.GetValue16(ip)
}

func (vm *VM) pushLiteral() {
	value := vm.NextInstruction16()
	vm.Push(value)
	vm.NextInstruction()
}

func (vm *VM) pushRegister() {
	vm.NextInstruction()
	regnum := vm.Fetch()
	registerName := vm.Mem.RegisterMap[int(regnum)]

	vm.Push(vm.Mem.GetRegister(registerName))

	vm.NextInstruction()
}

// If register contains a pointer
func (vm *VM) pushRegValue() {
	vm.NextInstruction()
	regnum := vm.Fetch()
	registerName := vm.Mem.RegisterMap[int(regnum)]

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
	registerName := vm.Mem.RegisterMap[int(regnum)]
	vm.Mem.SetRegister(registerName, x)

	vm.NextInstruction()
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

		case PUSHL:
			vm.pushLiteral()

		case PUSHR:
			vm.pushRegister()

		case PUSHREGV:
			vm.pushRegValue()

		case POP:
			vm.iPOP()

		case ADD:
			vm.iADD()

		default:
			halted = true
			break
		}
	}
}
