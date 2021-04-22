package compiler

// PRECEDENCE - The order of operations precedence
var PRECEDENCE map[string]int = map[string]int{
	"=": 1,

	"&&": 4,
	"||": 5,

	"<": 7,  ">": 7,  "<=": 7,  ">=": 7,  "==": 7,  "!=": 7,

	"+": 10,  "-": 10,
	"*": 20,  "/": 20,  "%": 20,
}

// ExprTypes - The different types for each expression or statement
const (
	None int = iota
	While
	For
	If
	Int
	Float
	String
	Boolean
	Class
	Access
	Array
	Variable
	Import
	Identifier
	Assign
	Binary
	Scope
	FunctionCall
	Function
	FunctionDecl
	Return
	Datatype
)