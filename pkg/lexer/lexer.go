package lexer

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/RyanCopley/expression-parser/pkg/tokens"
	"strconv"
	"strings"

	"github.com/RyanCopley/expression-parser/pkg/errors"
)

// isHexDigit returns true if ch is a valid hexadecimal digit.
func isHexDigit(ch byte) bool {
	return ('0' <= ch && ch <= '9') ||
		('a' <= ch && ch <= 'f') ||
		('A' <= ch && ch <= 'F')
}

// Lexer holds the state of the lexer.
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
}

// NewLexer creates a new Lexer for the given input.
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances positions.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

// peekChar returns the next character without advancing the lexer.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// isLetter checks if a character is a letter or underscore.
func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') ||
		('A' <= ch && ch <= 'Z') ||
		ch == '_'
}

// isDigit checks if a character is a digit.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// skipWhitespace skips over spaces, tabs, newlines, and also skips comments (lines starting with "#").
func (l *Lexer) skipWhitespace() {
	// Skip normal whitespace.
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
	// If a comment is encountered (line starts with "#"), skip until newline.
	for l.ch == '#' {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
		// Skip the newline.
		l.readChar()
		// Skip any whitespace after the comment.
		for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
			l.readChar()
		}
	}
}

// NextToken lexes and returns the next token.
func (l *Lexer) NextToken() (tokens.Token, error) {
	var tok tokens.Token

	l.skipWhitespace()
	startLine := l.line
	startColumn := l.column

	switch l.ch {
	case '+':
		tok = tokens.Token{Type: tokens.TokenPlus, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '-':
		tok = tokens.Token{Type: tokens.TokenMinus, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '*':
		tok = tokens.Token{Type: tokens.TokenMultiply, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '/':
		tok = tokens.Token{Type: tokens.TokenDivide, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = tokens.Token{Type: tokens.TokenLte, Literal: "<=", Line: startLine, Column: startColumn}
		} else {
			tok = tokens.Token{Type: tokens.TokenLt, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = tokens.Token{Type: tokens.TokenGte, Literal: ">=", Line: startLine, Column: startColumn}
		} else {
			tok = tokens.Token{Type: tokens.TokenGt, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = tokens.Token{Type: tokens.TokenEq, Literal: "==", Line: startLine, Column: startColumn}
		} else {
			tok = tokens.Token{Type: tokens.TokenIllegal, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = tokens.Token{Type: tokens.TokenNeq, Literal: "!=", Line: startLine, Column: startColumn}
		} else {
			tok = tokens.Token{Type: tokens.TokenNot, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = tokens.Token{Type: tokens.TokenAnd, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			err := errors.NewLexicalError("Unexpected character: &", startLine, startColumn)
			tok = tokens.Token{Type: tokens.TokenIllegal, Literal: string(l.ch), Line: startLine, Column: startColumn}
			l.readChar()
			return tok, err
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = tokens.Token{Type: tokens.TokenOr, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			err := errors.NewLexicalError("Unexpected character: |", startLine, startColumn)
			tok = tokens.Token{Type: tokens.TokenIllegal, Literal: string(l.ch), Line: startLine, Column: startColumn}
			l.readChar()
			return tok, err
		}
	case '(':
		tok = tokens.Token{Type: tokens.TokenLparen, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case ')':
		tok = tokens.Token{Type: tokens.TokenRparen, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '[':
		tok = tokens.Token{Type: tokens.TokenLeftBracket, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case ']':
		tok = tokens.Token{Type: tokens.TokenRightBracket, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '{':
		tok = tokens.Token{Type: tokens.TokenLeftCurly, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '}':
		tok = tokens.Token{Type: tokens.TokenRightCurly, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case ',':
		tok = tokens.Token{Type: tokens.TokenComma, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case ':':
		tok = tokens.Token{Type: tokens.TokenColon, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '.':
		tok = tokens.Token{Type: tokens.TokenDot, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '?':
		if l.peekChar() == '.' {
			l.readChar()
			tok = tokens.Token{Type: tokens.TokenQuestionDot, Literal: "?.", Line: startLine, Column: startColumn}
		} else if l.peekChar() == '[' {
			l.readChar()
			tok = tokens.Token{Type: tokens.TokenQuestionBracket, Literal: "?[", Line: startLine, Column: startColumn}
		} else {
			tok = tokens.Token{Type: tokens.TokenQuestion, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '$':
		tok = tokens.Token{Type: tokens.TokenDollar, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '"', '\'':
		str, err := l.readString(l.ch)
		if err != nil {
			tok = tokens.Token{Type: tokens.TokenIllegal, Literal: err.Error(), Line: startLine, Column: startColumn}
			return tok, err
		}
		tok = tokens.Token{Type: tokens.TokenString, Literal: str, Line: startLine, Column: startColumn}
		return tok, nil
	case 0:
		tok = tokens.Token{Type: tokens.TokenEof, Literal: "", Line: startLine, Column: startColumn}
	default:
		if isLetter(l.ch) {
			lit := l.readIdentifier()
			tok = tokens.Token{Type: lookupIdent(lit), Literal: lit, Line: startLine, Column: startColumn}
			return tok, nil
		} else if isDigit(l.ch) {
			return l.readNumber()
		} else {
			err := errors.NewLexicalError("Unexpected character: "+string(l.ch), startLine, startColumn)
			tok = tokens.Token{Type: tokens.TokenIllegal, Literal: string(l.ch), Line: startLine, Column: startColumn}
			l.readChar()
			return tok, err
		}
	}

	l.readChar()
	return tok, nil
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '-' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func lookupIdent(ident string) tokens.TokenType {
	keywords := map[string]tokens.TokenType{
		"true":  tokens.TokenBool,
		"false": tokens.TokenBool,
		"null":  tokens.TokenNull,
		"AND":   tokens.TokenAnd,
		"OR":    tokens.TokenOr,
		"NOT":   tokens.TokenNot,
	}
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return tokens.TokenIdent
}

func (l *Lexer) readNumber() (tokens.Token, error) {
	start := l.position
	startLine := l.line
	startColumn := l.column

	if l.ch == '-' || l.ch == '+' {
		sign := l.ch
		l.readChar()
		if !isDigit(l.ch) {
			return tokens.Token{
				Type:    tokens.TokenIllegal,
				Literal: l.input[start:l.position],
				Line:    startLine,
				Column:  startColumn,
			}, errors.NewLexicalError(fmt.Sprintf("Invalid number literal: '%c' not followed by a digit", sign), startLine, startColumn)
		}
	}

	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' {
		l.readChar()
		if !isDigit(l.ch) {
			return tokens.Token{
				Type:    tokens.TokenIllegal,
				Literal: l.input[start:l.position],
				Line:    startLine,
				Column:  startColumn,
			}, errors.NewLexicalError("Invalid number literal: missing digits after decimal point", startLine, l.position)
		}
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	if l.ch == 'e' || l.ch == 'E' {
		l.readChar()
		if l.ch == '-' || l.ch == '+' {
			l.readChar()
		}
		if !isDigit(l.ch) {
			return tokens.Token{
				Type:    tokens.TokenIllegal,
				Literal: l.input[start:l.position],
				Line:    startLine,
				Column:  startColumn,
			}, errors.NewLexicalError("Invalid number literal: missing digits in exponent", startLine, startColumn)
		}
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return tokens.Token{
		Type:    tokens.TokenNumber,
		Literal: l.input[start:l.position],
		Line:    startLine,
		Column:  startColumn,
	}, nil
}

func (l *Lexer) readString(quote byte) (string, error) {
	startLine := l.line
	startColumn := l.column
	var sb strings.Builder
	escaped := false

	l.readChar() // skip opening quote
	for l.ch != 0 {
		if escaped {
			if l.ch == 'u' {
				// Read next 4 hexadecimal digits.
				hexDigits := ""
				for i := 0; i < 4; i++ {
					l.readChar()
					if !isHexDigit(l.ch) {
						return "", errors.NewLexicalError("Invalid unicode escape sequence", l.line, l.column)
					}
					hexDigits += string(l.ch)
				}
				code, err := strconv.ParseInt(hexDigits, 16, 32)
				if err != nil {
					return "", errors.NewLexicalError("Invalid unicode escape sequence", l.line, l.column)
				}
				sb.WriteRune(rune(code))
				escaped = false
			} else {
				switch l.ch {
				case 'n':
					sb.WriteByte('\n')
				case 'r':
					sb.WriteByte('\r')
				case 't':
					sb.WriteByte('\t')
				case '\\':
					sb.WriteByte('\\')
				case '"':
					sb.WriteByte('"')
				case '\'':
					sb.WriteByte('\'')
				default:
					return "", errors.NewLexicalError("Invalid escape sequence: \\"+string(l.ch), l.line, l.column)
				}
			}
			escaped = false
		} else {
			if l.ch == '\\' {
				escaped = true
			} else if l.ch == quote {
				l.readChar()
				return sb.String(), nil
			} else {
				sb.WriteByte(l.ch)
			}
		}
		l.readChar()
	}
	return "", errors.NewLexicalError("Unclosed string literal", startLine, startColumn)
}

func (l *Lexer) ExportTokens() ([]byte, error) {
	var buf bytes.Buffer
	for {
		tok, err := l.NextToken()
		if err != nil {
			return nil, err
		}
		code, ok := tokens.TokenTypeToByte[tok.Type]
		if !ok {
			return nil, fmt.Errorf("unknown token type: %v", tok.Type)
		}
		buf.WriteByte(code)

		if fixed, exists := tokens.FixedTokenLiterals[tok.Type]; exists && tok.Literal == fixed {
			// No literal data needed.
		} else {
			literalBytes := []byte(tok.Literal)
			if len(literalBytes) > 255 {
				return nil, fmt.Errorf("literal too long")
			}
			buf.WriteByte(byte(len(literalBytes)))
			buf.Write(literalBytes)
		}

		if tok.Type == tokens.TokenEof {
			break
		}
	}
	return buf.Bytes(), nil
}

func (l *Lexer) ExportTokensSigned(priv *rsa.PrivateKey) ([]byte, error) {
	tokenData, err := l.ExportTokens()
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(tokenData)
	signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hash[:])
	if err != nil {
		return nil, err
	}

	if len(tokenData) > int(^uint32(0)) {
		return nil, fmt.Errorf("token data length %d exceeds maximum allowed size", len(tokenData))
	}

	tokenLen := uint32(len(tokenData))

	var buf bytes.Buffer
	buf.WriteString(tokens.HeaderMagic)

	if err := binary.Write(&buf, binary.LittleEndian, tokenLen); err != nil {
		return nil, err
	}
	buf.Write(tokenData)
	buf.Write(signature)

	return buf.Bytes(), nil
}
