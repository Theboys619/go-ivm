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

var StringTypes map[string]int = map[string]int{
	"None": None,
	"While": While,
	"For": For,
	"If": If,
	"Int": Int,
	"Float": Float,
	"String": String,
	"Boolean": Boolean,
	"Class": Class,
	"Access": Access,
	"Array": Array,
	"Variable": Variable,
	"Import": Import,
	"Identifier": Identifier,
	"Assign": Assign,
	"Binary": Binary,
	"Scope": Scope,
	"FunctionCall": FunctionCall,
	"Function": Function,
	"FunctionDecl": FunctionDecl,
	"Return": Return,
	"Datatype": Datatype,
}

// AST - Main ast containing statements and expressions
type AST struct {
	ExprType int
	TokValue Token

	Op Token

	Left *AST
	Right *AST

	Block []*AST
	Args []*AST
}

// NewAST - Create a AST Object
func NewAST(ExprType int, value Token) *AST {
	return &AST{
		ExprType: ExprType,
		TokValue: value,
		Op: *NewToken("Nothing", "Nothing"),
	}
}

// Parser - Main parser struct
type Parser struct {
	Tokens []Token
	CurTok Token

	Ast *AST

	pos int
}

func (parser *Parser) advance(amt... int) Token {
	amount := 1
	if len(amt) > 0 {
		amount = amt[0]
	}

	if parser.pos + amount >= len(parser.Tokens) {
		parser.CurTok = *NewToken("Null", "Null")
	}
	parser.pos += amount;

	parser.CurTok = parser.Tokens[parser.pos]
	return parser.CurTok
}

func (parser *Parser) peek(amt... int) Token {
	amount := 1
	if len(amt) > 0 { amount = amt[0] }

	if parser.pos + amount >= len(parser.Tokens) { return *NewToken("Null", "Null") }

	return parser.Tokens[parser.pos + amount];
}

func (parser *Parser) isNull() bool {
	return parser.CurTok.tokenType == "Nothing" || parser.CurTok.tokenType == "Null"
}

func (parser *Parser) isEOF() bool {
	return parser.CurTok.tokenType == "EOF"
}

// isType - checks the type, value, or both on the token
func (parser *Parser) isType(tokentype string, value string, tok Token) bool {
	if tokentype == "" && len(value) > 0 {
		return tok.GetValue() == value
	}
	return tok.tokenType == tokentype && (value == "" || tok.GetValue() == value)
}

func (parser *Parser) isIgnore(tok Token) bool {
	return tok.value == ";" || tok.tokenType == "Linebreak"
}

func (parser *Parser) skipOver(tokentype string, value string, tok Token) Token {
	if tok.tokenType == "Nothing" {
		tok = parser.CurTok
	}

	if tokentype == "" && len(value) > 0 {
		if tok.GetValue() == value {
			return parser.advance()
		}

		msg := "Unexpeced token " + tok.GetValue()
		panic(msg) // CreateError
	}

	if tok.tokenType == tokentype && (value == "" || tok.GetValue() == value) {
		return parser.advance()
	} else {
		msg := "Unexpeced token " + tok.GetValue()
		panic(msg) // CreateError
	}
}

func (parser *Parser) skipIgnore() {
	for parser.isIgnore(parser.CurTok) {
		parser.advance()
	}
}

func (parser *Parser) pDelimiters(start string, end string, separator string) []*AST {
	var values []*AST = make([]*AST, 0, 10)
	isFirst := true

	parser.skipOver("", start, parser.CurTok)

	for !parser.isEOF() {
		if parser.isType("Delimiter", end, parser.CurTok) {
			break
		} else if isFirst {
			isFirst = false
		} else {
			if separator == "" && parser.isIgnore(parser.CurTok) {
				parser.skipIgnore()
			} else {
				parser.skipOver("", separator, parser.CurTok)
			}
		}

		val := parser.pExpression();
		values = append(values, val)
	}
	parser.skipOver("", end, parser.CurTok)

	return values;
}

func (parser *Parser) isCallable(callStmt *AST) bool {
	return callStmt.ExprType != Function && callStmt.ExprType != If && callStmt.ExprType != Return
}

func (parser *Parser) checkCall(expr *AST) *AST {
	if parser.isType("Delimiter", "(", parser.peek()) && parser.isCallable(expr) {
		// return parser.pCall(expr)
	}

	return expr
}

func (parser *Parser) checkBinary(left *AST, prec int) *AST {
	op := parser.CurTok

	if parser.isType("Operator", "", op) {
		opvalue := op.GetValue()
		newPrec := PRECEDENCE[opvalue]

		if prec < newPrec {
			var assigns map[string]int = map[string]int{
				"=": Assign,
				"+=": Assign,
			}

			optype, ok := assigns[opvalue]

			if !ok {
				optype = Binary
			}

			expr := NewAST(optype, parser.CurTok)
			expr.Left = left
			expr.Op = op
			expr.Right = parser.checkBinary(parser.checkCall(parser.pAll()), newPrec);
			// expr.DataType = expr->right->dataType;

			return parser.checkBinary(expr, prec)
		}
	}

	return left
}

func (parser *Parser) pAll() *AST {

}

// pExpression
func (parser *Parser) pExpression() *AST {
	return parser.checkCall(parser.checkBinary(parser.checkCall(parser.pAll()), 0))
}

// Parse - Starts the parsing of the tokens
func (parser *Parser) Parse(tokens []Token) {
	parser.Tokens = tokens
	parser.Ast = NewAST(Scope, *NewToken("Nothing", "Nothing"))

	var block []*AST = make([]*AST, 0, 50)

	for !parser.isNull() && !parser.isEOF() {
		block = append(block, parser.pExpression())
	}
}