package tokens

const HeaderMagic = "STOK" // 4-byte header magic

// TokenType defines the type for tokens.
type TokenType uint8

const (
	TokenEof TokenType = iota
	TokenIllegal
	TokenIdent
	TokenNumber
	TokenString
	TokenBool
	TokenNull
	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide
	TokenLt
	TokenGt
	TokenLte
	TokenGte
	TokenEq
	TokenNeq
	TokenAnd
	TokenOr
	TokenNot
	TokenLparen
	TokenRparen
	TokenLeftBracket
	TokenRightBracket
	TokenLeftCurly
	TokenRightCurly
	TokenComma
	TokenColon
	TokenDot
	TokenQuestion
	TokenQuestionDot
	TokenQuestionBracket
	TokenDollar
)

// Token represents a lexical token.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// TokenTypeToByte maps each TokenType to a unique byte code.
var TokenTypeToByte = map[TokenType]byte{
	TokenEof:             0,
	TokenIllegal:         1,
	TokenIdent:           2,
	TokenNumber:          3,
	TokenString:          4,
	TokenBool:            5,
	TokenNull:            6,
	TokenPlus:            7,
	TokenMinus:           8,
	TokenMultiply:        9,
	TokenDivide:          10,
	TokenLt:              11,
	TokenGt:              12,
	TokenLte:             13,
	TokenGte:             14,
	TokenEq:              15,
	TokenNeq:             16,
	TokenAnd:             17,
	TokenOr:              18,
	TokenNot:             19,
	TokenLparen:          20,
	TokenRparen:          21,
	TokenLeftBracket:     22,
	TokenRightBracket:    23,
	TokenLeftCurly:       24,
	TokenRightCurly:      25,
	TokenComma:           26,
	TokenColon:           27,
	TokenDot:             28,
	TokenQuestionDot:     30,
	TokenQuestionBracket: 31,
	TokenDollar:          32,
}

// FixedTokenLiterals defines fixed literal strings for tokens.
var FixedTokenLiterals = map[TokenType]string{
	TokenPlus:            "+",
	TokenMinus:           "-",
	TokenMultiply:        "*",
	TokenDivide:          "/",
	TokenLt:              "<",
	TokenGt:              ">",
	TokenLte:             "<=",
	TokenGte:             ">=",
	TokenEq:              "==",
	TokenNeq:             "!=",
	TokenAnd:             "AND",
	TokenOr:              "OR",
	TokenNot:             "NOT",
	TokenLparen:          "(",
	TokenRparen:          ")",
	TokenLeftBracket:     "[",
	TokenRightBracket:    "]",
	TokenLeftCurly:       "{",
	TokenRightCurly:      "}",
	TokenComma:           ",",
	TokenColon:           ":",
	TokenDot:             ".",
	TokenQuestionDot:     "?.",
	TokenQuestionBracket: "?[",
	TokenDollar:          "$",
}
