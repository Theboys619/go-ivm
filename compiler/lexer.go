package compiler

import "fmt"

// Token - A basic Token struct for identifier token
type Token struct {
	tokenType string
	value string

	line int
	index int

	file string
}

// NewToken - Creates a new Token
func NewToken(tokenType string, value string) *Token {
	return &Token{
		tokenType: tokenType,
		value: value,
		line: 0,
		index: 0,
		file: "unknown",
	}
}

// NewNullToken - Creates a new Token a Null type
func NewNullToken() *Token {
	return NewToken("Nothing", "Nothing")
}

// SetFile - Set file inside token
func (token *Token) SetFile(file string) {
	token.file = file
}

// SetPos - Set line and index
func (token *Token) SetPos(line int, index int) {
	token.line = line
	token.index = index
}

// GetValue - Get the value of the token as a string
func (token *Token) GetValue() string {
	return token.value
}

// Lexer - Tokenizes given input
type Lexer struct {
	input string
	length int
	file string

	tokens []Token

	pos int
	line int
	index int
	curchar byte
}

// NewLexer - Creates a new Lexer
func NewLexer() *Lexer {
	return &Lexer{
		input: "",
		length: 0,
		file: "unknown",
		tokens: make([]Token, 0, 100),
		pos: 1,
		line: 1,
		index: 0,
		curchar: '0',
	}
}

var keywords map[string]string = map[string]string{
	"make": "Keyword",
	"send": "Keyword",
	"oftype": "Keyword",
	"class": "Keyword",
	"any": "typeKW",
	"if": "Keyword",
	"else": "Keyword",
	"import": "Keyword",
	"from": "Keyword",
	"as": "Keyword",
    "true": "Keyword",
    "false": "Keyword",
	"loop": "Keyword",
	
	// Type Keywords
	"function": "typeKW",
	"number": "typeKW",
	"boolean": "typeKW",
	"string": "typeKW",
	"nothing": "typeKW",
}
var operators []map[string]string = []map[string]string{
	0: nil,
	1: map[string]string{
		"!": "NotOp",
		"=": "Assign",
		"+": "Plus",
		"-": "Minus",
		"*": "Multiply",
		"/": "Divide",
		"<": "LessThan",
		">": "GreaterThan",
		"%": "Modulus",
	},
	2: map[string]string{
		"==": "isEqual",
		"!=": "NotEqual",
		"+=": "PlusEqual",
		"-=": "MinusEqual",
		"*=": "MultiplyEqual",
		"/=": "DivideEqual",
		"%=": "ModulusEqual",
		">=": "GreaterEqual",
		"<=": "LessEqual",
		"&&": "AndOp",
		"||": "OrOp",
	},
}

func (lexer *Lexer) reset() {
	lexer.input = ""
	lexer.length = 0
	lexer.pos = 1
	lexer.index = 0
	lexer.line = 1
	lexer.tokens = make([]Token, 0, 100)
	lexer.curchar = '0'
	lexer.file = "unknown"
}

func (lexer *Lexer) advance(amt... int) byte {
	amount := 1
	if len(amt) > 0 {
		amount = amt[0]
	}

	lexer.index += amount
	lexer.pos += amount

	if lexer.index >= lexer.length {
		lexer.curchar = 0
		return 0
	}

	lexer.curchar = lexer.input[lexer.index]
	return lexer.curchar
}

func (lexer *Lexer) peek(amt... int) byte {
	amount := 1
	if len(amt) > 0 {
		amount = amt[0]
	}

	lookupIndex := lexer.index + amount;

	if lookupIndex >= lexer.length {
		return 0
	}

	return lexer.input[lookupIndex]
}

func (lexer *Lexer) grab(amt int) string {
	value := string(lexer.curchar)

	for i := 1; i < amt; i++ {
		value += string(lexer.peek())
	}

	return value
}

func (lexer *Lexer) isDelimiter(char byte) bool {
	switch char {
	case '(', ')', '{', '}', '[', ']', ':', ';', ',', '.':
		return true
	default:
		return false
	}
}

func (lexer *Lexer) isOperator() int {
	for i := len(operators)-1; i > 0; i-- {
		value := operators[i]

		for op := range value {
			if lexer.grab(i) == op {
				return i
			}
		}
	}

	return 0
}

func (lexer *Lexer) isWhitespace(c byte) bool {
	return c == ' ' || c == '\r' || c == '\t';
}

func (lexer *Lexer) isAlpha(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z');
}

func (lexer *Lexer) isDigit(c byte) bool {
	return c >= '0' && c <= '9';
}

func (lexer *Lexer) isNumber(c byte) bool {
	return c == '-' && lexer.isDigit(lexer.peek()) || lexer.isDigit(c)
}

func (lexer *Lexer) isQuote(c byte) bool { return c == '\'' || c == '"'; };

// Tokenize - Creates tokens from input
func (lexer *Lexer) Tokenize(input string) []Token {
	lexer.reset()
	lexer.input = input
	lexer.length = len(input)
	lexer.curchar = input[0]

	for lexer.curchar != 0 {
		var oldPos int = lexer.pos

		if lexer.isWhitespace(lexer.curchar) {
			lexer.advance()
		}

		if lexer.curchar == '\n' {
			lexer.advance()
			lexer.pos = 1
			lexer.line++

			lexer.tokens = append(lexer.tokens, *NewToken("Linebreak", "\\n"))
		}

		if lexer.isDelimiter(lexer.curchar) {
			lexer.tokens = append(lexer.tokens, *NewToken("Delimiter", string(lexer.curchar)))
			lexer.advance()
		}

		if opval := lexer.isOperator(); opval > 0 {
			value := lexer.grab(opval)
			lexer.tokens = append(lexer.tokens, *NewToken("Operator", value))
			lexer.advance(opval)
		}

		if lexer.isNumber(lexer.curchar) {
			numtype := "Int"
			index := lexer.pos
			line := lexer.line

			value := ""

			if (lexer.curchar == '-') {
				value += string(lexer.curchar)
				lexer.advance()
			}

			for lexer.isNumber(lexer.curchar) {
				value += string(lexer.curchar)
				lexer.advance()

				if (lexer.curchar == '.') {
					numtype = "Float"
					value += "."

					lexer.advance()
				}
			}

			tok := NewToken(numtype, value);
			tok.SetPos(line, index)
			tok.SetFile(lexer.file);

			lexer.tokens = append(lexer.tokens, *tok);
		}

		if lexer.isQuote(lexer.curchar) {
			quote := lexer.curchar
			index := lexer.pos
			line := lexer.line

			value := ""
			lexer.advance()

			for lexer.curchar != 0 && lexer.curchar != quote {
				if lexer.curchar == '\n' { panic("\\n") }
				value += string(lexer.curchar)
				lexer.advance()
			}

			lexer.advance()

			tok := NewToken("String", value)
			tok.SetPos(line, index)
			tok.SetFile(lexer.file)

			lexer.tokens = append(lexer.tokens, *tok)
		}

		if lexer.isAlpha(lexer.curchar) {
			value := ""

			index := lexer.pos
			line := lexer.line

			for lexer.curchar != 0 && lexer.isAlpha(lexer.curchar) {
				value += string(lexer.curchar)
				lexer.advance()
			}

			kwtype, ok := keywords[value]

			if !ok {
				kwtype = "Identifier"
			}

			tok := NewToken(kwtype, value);
			tok.SetPos(line, index)
			tok.SetFile(lexer.file);

			lexer.tokens = append(lexer.tokens, *tok)
		}

		if lexer.pos == oldPos {
			fmt.Println("Invalid tokens somewhere") // CreateError
			break
		}
	}

	tok := NewToken("EOF", "EOF")
	tok.SetFile(lexer.file)
	tok.SetPos(lexer.line, lexer.pos)
	lexer.tokens = append(lexer.tokens, *tok)

	return lexer.tokens
}