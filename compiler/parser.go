package compiler

// PRECEDENCE - The order of operations precedence
var PRECEDENCE map[string]int = map[string]int{
	"=": 1,

	"&&": 4,
	"||": 5,

	"<": 7, ">": 7, "<=": 7, ">=": 7, "==": 7, "!=": 7,

	"+": 10, "-": 10,
	"*": 20, "/": 20, "%": 20,
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
	"None":         None,
	"While":        While,
	"For":          For,
	"If":           If,
	"Int":          Int,
	"Float":        Float,
	"String":       String,
	"Boolean":      Boolean,
	"Class":        Class,
	"Access":       Access,
	"Array":        Array,
	"Variable":     Variable,
	"Import":       Import,
	"Identifier":   Identifier,
	"Assign":       Assign,
	"Binary":       Binary,
	"Scope":        Scope,
	"FunctionCall": FunctionCall,
	"Function":     Function,
	"FunctionDecl": FunctionDecl,
	"Return":       Return,
	"Datatype":     Datatype,
}

// AST - Main ast containing statements and expressions
type AST struct {
	ExprType int
	TokValue Token

	Op Token

	Left  *AST
	Right *AST

	Block []*AST
	Args  []*AST
	Scope *AST

	DataType string
}

// NewAST - Create a AST Object
func NewAST(ExprType int, value Token) *AST {
	return &AST{
		ExprType: ExprType,
		TokValue: value,
		Op:       *NewToken("Nothing", "Nothing"),
	}
}

// Parser - Main parser struct
type Parser struct {
	Tokens []Token
	CurTok Token

	Ast *AST

	pos int
}

func (parser *Parser) advance(amt ...int) Token {
	amount := 1
	if len(amt) > 0 {
		amount = amt[0]
	}

	if parser.pos+amount >= len(parser.Tokens) {
		parser.CurTok = *NewToken("Nothing", "Nothing")
		return parser.CurTok
	}
	parser.pos += amount

	parser.CurTok = parser.Tokens[parser.pos]
	return parser.CurTok
}

func (parser *Parser) peek(amt ...int) Token {
	amount := 1
	if len(amt) > 0 {
		amount = amt[0]
	}

	if parser.pos+amount >= len(parser.Tokens) {
		return *NewToken("Nothing", "Nothing")
	}

	return parser.Tokens[parser.pos+amount]
}

func (parser *Parser) isNull(tok Token) bool {
	return tok.tokenType == "Nothing" || tok.tokenType == "Null"
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
	for parser.isIgnore(parser.CurTok) && !parser.isNull(parser.CurTok) && !parser.isEOF() {
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

		val := parser.pExpression()
		values = append(values, val)
	}
	parser.skipOver("", end, parser.CurTok)

	return values
}

func (parser *Parser) isCallable(callStmt *AST) bool {
	return callStmt.ExprType != Function && callStmt.ExprType != If && callStmt.ExprType != Return
}

func (parser *Parser) checkCall(expr *AST) *AST {
	if parser.isType("Delimiter", "(", parser.peek()) && parser.isCallable(expr) {
		return parser.pCall(expr)
	}

	return expr
}

func (parser *Parser) pCall(expr *AST) *AST {
	funcCall := NewAST(FunctionCall, expr.TokValue)
	// funcCall->dotOp = nullptr;

	parser.advance()

	funcCall.Args = parser.pDelimiters("(", ")", ",")

	// pIndexAccess(funcCall)
	// pDotOp(funcCall)

	return funcCall
}

func (parser *Parser) checkBinary(left *AST, prec int) *AST {
	op := parser.CurTok

	if parser.isType("Operator", "", op) {
		opvalue := op.GetValue()
		newPrec := PRECEDENCE[opvalue]

		if prec < newPrec {
			parser.advance()

			var assigns map[string]int = map[string]int{
				"=":  Assign,
				"+=": Assign,
			}

			optype, ok := assigns[opvalue]

			if !ok {
				optype = Binary
			}

			expr := NewAST(optype, parser.CurTok)
			expr.Left = left
			expr.Op = op
			expr.Right = parser.checkBinary(parser.checkCall(parser.pAll()), newPrec)
			// expr.DataType = expr->right->dataType;

			return parser.checkBinary(expr, prec)
		}
	}

	return left
}

func (parser *Parser) pIdentifier(expr *AST) *AST {
	expr.ExprType = Identifier

	if !parser.isType("Delimiter", "(", parser.peek()) {
		parser.advance()
	}

	// pIndexAccess(expr);

	// Could / should change
	// if (!pDotOp(expr)) {
	// 	expr->dataType = pDatatype();
	// }

	return expr
}

func (parser *Parser) pFunction() *AST {
	if !parser.isType("Identifier", "", parser.CurTok) {
		panic(parser.CurTok.GetValue())
	}

	funcToken := parser.CurTok

	ifunc := NewAST(Function, parser.CurTok)
	parser.advance()

	ifunc.Args = parser.pDelimiters("(", ")", ",")
	ifunc.Scope = NewAST(Scope, funcToken)

	// if parser.CurTok.getString() == "oftype" {
	// 	ifunc.DataType = parser.pDatatype()
	// } else {
	// 	ifunc.DataType = "any";
	// }

	ifunc.Scope.Block = parser.pDelimiters("{", "}", "")

	// TODO function types and return types

	return ifunc
}

func (parser *Parser) pMake(tok *AST) *AST {
	parser.advance()

	if parser.isType("Delimiter", "(", parser.peek()) {
		return parser.pFunction()
	}

	identifier := parser.pIdentifier(NewAST(None, parser.CurTok))
	identifier.ExprType = Variable

	return identifier
}

func (parser *Parser) pAll() *AST {
	if parser.isType("Delimiter", "(", parser.CurTok) {
		parser.advance()
		expr := parser.pExpression()
		parser.skipOver("Delimiter", ")", parser.CurTok)
		return expr
	}

	token := NewAST(None, parser.CurTok)

	if parser.isType("Keyword", "make", parser.CurTok) {
		return parser.pMake(token)
	}

	if parser.isType("String", "", parser.CurTok) {
		parser.advance()
		token.ExprType = String

		return token
	}

	if parser.isType("Int", "", parser.CurTok) {
		parser.advance()
		token.ExprType = Int

		return token
	}

	if parser.isType("Identifier", "", parser.CurTok) {
		return parser.pIdentifier(token)
	}

	if parser.isType("Linebreak", "", parser.CurTok) {
		for parser.isType("Linebreak", "", parser.CurTok) {
			parser.advance()
		}

		return parser.pAll()
	}

	panic("Unexpected error. Failed parsing. Got token " + parser.CurTok.GetValue())
}

// pExpression
func (parser *Parser) pExpression() *AST {
	return parser.checkCall(parser.checkBinary(parser.checkCall(parser.pAll()), 0))
}

// Parse - Starts the parsing of the tokens
func (parser *Parser) Parse(tokens []Token) *AST {
	parser.LoadTokens(tokens)
	parser.Ast = NewAST(Scope, *NewToken("Nothing", "_MAIN_"))

	var block []*AST = make([]*AST, 0, 50)

	for !parser.isNull(parser.CurTok) && !parser.isEOF() {
		expr := parser.pExpression()
		block = append(block, expr)

		if parser.isIgnore(parser.CurTok) {
			parser.skipIgnore()
		}
	}
	parser.Ast.Block = block

	return parser.Ast
}

// LoadTokens - Loads tokens into parser
func (parser *Parser) LoadTokens(tokens []Token) {
	parser.Tokens = tokens
	parser.CurTok = tokens[0]
}

func NewParser() *Parser {
	return &Parser{
		Tokens: make([]Token, 0, 20),
		Ast:    nil,
		CurTok: *NewNullToken(),
		pos:    0,
	}
}
