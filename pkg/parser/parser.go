package parser

import (
	"fmt"
	"github.com/RyanCopley/expression-parser/pkg/ast/expressions"
	"github.com/RyanCopley/expression-parser/pkg/tokens"
	"github.com/RyanCopley/expression-parser/pkg/types"
	"strings"

	"github.com/RyanCopley/expression-parser/pkg/ast"
	"github.com/RyanCopley/expression-parser/pkg/errors"
)

// TokenStream represents a stream of tokens.
type TokenStream interface {
	NextToken() (tokens.Token, error)
}

// Parser holds the state for parsing.
type Parser struct {
	lexer     TokenStream
	curToken  tokens.Token
	peekToken tokens.Token
	errors    []string
}

// NewParser creates a new parser.
func NewParser(l TokenStream) (*Parser, error) {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}
	if err := p.nextToken(); err != nil {
		return nil, err
	}
	if err := p.nextToken(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Parser) nextToken() error {
	p.curToken = p.peekToken
	tok, err := p.lexer.NextToken()
	if err != nil {
		return err
	}
	p.peekToken = tok
	return nil
}

func (p *Parser) ParseExpression() (ast.Expression, error) {
	return p.parseOrExpression()
}

const (
	_ int = iota
	LOWEST
	OR
	AND
	EQUALS
	GTR
	SUM
	PRODUCT
	CALL
	MEMBER
)

var precedences = map[tokens.TokenType]int{
	tokens.TokenOr:              OR,
	tokens.TokenAnd:             AND,
	tokens.TokenEq:              EQUALS,
	tokens.TokenNeq:             EQUALS,
	tokens.TokenLt:              GTR,
	tokens.TokenGt:              GTR,
	tokens.TokenLte:             GTR,
	tokens.TokenGte:             GTR,
	tokens.TokenPlus:            SUM,
	tokens.TokenMinus:           SUM,
	tokens.TokenMultiply:        PRODUCT,
	tokens.TokenDivide:          PRODUCT,
	tokens.TokenLparen:          CALL,
	tokens.TokenDot:             MEMBER,
	tokens.TokenLeftBracket:     MEMBER,
	tokens.TokenQuestionDot:     MEMBER,
	tokens.TokenQuestionBracket: MEMBER,
}

func (p *Parser) curPrecedence() int {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) parseOrExpression() (ast.Expression, error) {
	left, err := p.parseAndExpression()
	if err != nil {
		return nil, err
	}
	for p.curTokenIs(tokens.TokenOr) || (p.curTokenIs(tokens.TokenIdent) && strings.ToUpper(p.curToken.Literal) == "OR") {
		operator := p.curToken
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		right, err := p.parseAndExpression()
		if err != nil {
			return nil, err
		}
		left = &expressions.BinaryExpr{
			Left:     left,
			Operator: operator.Type,
			Right:    right,
			Line:     operator.Line,
			Column:   operator.Column,
		}
	}
	return left, nil
}

func (p *Parser) parseAndExpression() (ast.Expression, error) {
	left, err := p.parseEqualityExpression()
	if err != nil {
		return nil, err
	}
	for p.curTokenIs(tokens.TokenAnd) || (p.curTokenIs(tokens.TokenIdent) && strings.ToUpper(p.curToken.Literal) == "AND") {
		operator := p.curToken
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		right, err := p.parseEqualityExpression()
		if err != nil {
			return nil, err
		}
		left = &expressions.BinaryExpr{
			Left:     left,
			Operator: operator.Type,
			Right:    right,
			Line:     operator.Line,
			Column:   operator.Column,
		}
	}
	return left, nil
}

func (p *Parser) parseEqualityExpression() (ast.Expression, error) {
	left, err := p.parseRelationalExpression()
	if err != nil {
		return nil, err
	}
	for p.curTokenIs(tokens.TokenEq) || p.curTokenIs(tokens.TokenNeq) {
		operator := p.curToken
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		right, err := p.parseRelationalExpression()
		if err != nil {
			return nil, err
		}
		left = &expressions.BinaryExpr{
			Left:     left,
			Operator: operator.Type,
			Right:    right,
			Line:     operator.Line,
			Column:   operator.Column,
		}
	}
	return left, nil
}

func (p *Parser) parseRelationalExpression() (ast.Expression, error) {
	left, err := p.parseAdditiveExpression()
	if err != nil {
		return nil, err
	}
	for p.curTokenIs(tokens.TokenLt) || p.curTokenIs(tokens.TokenGt) || p.curTokenIs(tokens.TokenLte) || p.curTokenIs(tokens.TokenGte) {
		operator := p.curToken
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		right, err := p.parseAdditiveExpression()
		if err != nil {
			return nil, err
		}
		left = &expressions.BinaryExpr{
			Left:     left,
			Operator: operator.Type,
			Right:    right,
			Line:     operator.Line,
			Column:   operator.Column,
		}
	}
	return left, nil
}

func (p *Parser) parseAdditiveExpression() (ast.Expression, error) {
	left, err := p.parseMultiplicativeExpression()
	if err != nil {
		return nil, err
	}
	for p.curTokenIs(tokens.TokenPlus) || p.curTokenIs(tokens.TokenMinus) {
		operator := p.curToken
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		right, err := p.parseMultiplicativeExpression()
		if err != nil {
			return nil, err
		}
		left = &expressions.BinaryExpr{
			Left:     left,
			Operator: operator.Type,
			Right:    right,
			Line:     operator.Line,
			Column:   operator.Column,
		}
	}
	return left, nil
}

func (p *Parser) parseMultiplicativeExpression() (ast.Expression, error) {
	left, err := p.parseUnaryExpression()
	if err != nil {
		return nil, err
	}
	for p.curTokenIs(tokens.TokenMultiply) || p.curTokenIs(tokens.TokenDivide) {
		operator := p.curToken
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		right, err := p.parseUnaryExpression()
		if err != nil {
			return nil, err
		}
		left = &expressions.BinaryExpr{
			Left:     left,
			Operator: operator.Type,
			Right:    right,
			Line:     operator.Line,
			Column:   operator.Column,
		}
	}
	return left, nil
}

func (p *Parser) parseUnaryExpression() (ast.Expression, error) {
	if p.curTokenIs(tokens.TokenNot) || p.curTokenIs(tokens.TokenMinus) {
		operator := p.curToken
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		expr, err := p.parseUnaryExpression()
		if err != nil {
			return nil, err
		}
		return &expressions.UnaryExpr{
			Operator: operator.Type,
			Expr:     expr,
			Line:     operator.Line,
			Column:   operator.Column,
		}, nil
	}
	return p.parseMemberAccessExpression()
}

func (p *Parser) parseMemberAccessExpression() (ast.Expression, error) {
	expr, err := p.parsePrimaryExpressionInner()
	if err != nil {
		return nil, err
	}
	for p.curTokenIs(tokens.TokenDot) || p.curTokenIs(tokens.TokenLeftBracket) || p.curTokenIs(tokens.TokenQuestionDot) || p.curTokenIs(tokens.TokenQuestionBracket) {
		var part expressions.MemberPart
		if p.curTokenIs(tokens.TokenDot) || p.curTokenIs(tokens.TokenQuestionDot) {
			optional := p.curTokenIs(tokens.TokenQuestionDot)
			if err := p.nextToken(); err != nil {
				return nil, err
			}
			if !p.curTokenIs(tokens.TokenIdent) && p.curToken.Type != tokens.TokenString {
				return nil, errors.NewSyntaxError(fmt.Sprintf("Expected identifier after dot at line %d, column %d", p.curToken.Line, p.curToken.Column), p.curToken.Line, p.curToken.Column)
			}
			part = expressions.MemberPart{Optional: optional, IsIndex: false, Key: strings.TrimSpace(p.curToken.Literal), Line: p.curToken.Line, Column: p.curToken.Column}
			if err := p.nextToken(); err != nil {
				return nil, err
			}
		} else {
			optional := p.curTokenIs(tokens.TokenQuestionBracket)
			if err := p.nextToken(); err != nil {
				return nil, err
			}
			exprTmp, err := p.ParseExpression()
			if err != nil {
				return nil, err
			}
			indexExpr := exprTmp
			if !p.curTokenIs(tokens.TokenRightBracket) {
				return nil, errors.NewSyntaxError(fmt.Sprintf("Expected closing bracket at line %d, column %d", p.curToken.Line, p.curToken.Column), p.curToken.Line, p.curToken.Column)
			}
			if err := p.nextToken(); err != nil {
				return nil, err
			}
			part = expressions.MemberPart{Optional: optional, IsIndex: true, Expr: indexExpr, Line: p.curToken.Line, Column: p.curToken.Column}
		}
		if mae, ok := expr.(*expressions.MemberAccessExpr); ok {
			mae.AccessParts = append(mae.AccessParts, part)
		} else {
			expr = &expressions.MemberAccessExpr{Target: expr, AccessParts: []expressions.MemberPart{part}}
		}
	}
	return expr, nil
}

func (p *Parser) parsePrimaryExpressionInner() (ast.Expression, error) {
	switch p.curToken.Type {
	case tokens.TokenLparen:
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		expr, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		if !p.curTokenIs(tokens.TokenRparen) {
			return nil, errors.NewSyntaxError("Expected RPAREN", p.curToken.Line, p.curToken.Column)
		}
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return expr, nil

	case tokens.TokenNumber:
		lit := &expressions.LiteralExpr{
			Value:  types.ParseNumber(p.curToken.Literal),
			Line:   p.curToken.Line,
			Column: p.curToken.Column,
		}
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return lit, nil

	case tokens.TokenString:
		lit := &expressions.LiteralExpr{
			Value:  p.curToken.Literal,
			Line:   p.curToken.Line,
			Column: p.curToken.Column,
		}
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return lit, nil

	case tokens.TokenBool:
		var val bool
		if p.curToken.Literal == "true" {
			val = true
		} else {
			val = false
		}
		lit := &expressions.LiteralExpr{
			Value:  val,
			Line:   p.curToken.Line,
			Column: p.curToken.Column,
		}
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return lit, nil

	case tokens.TokenNull:
		lit := &expressions.LiteralExpr{
			Value:  nil,
			Line:   p.curToken.Line,
			Column: p.curToken.Column,
		}
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return lit, nil

	case tokens.TokenDollar:
		return p.parseContextExpression()
	case tokens.TokenLeftCurly:
		return p.parseObjectLiteral()
	case tokens.TokenLeftBracket:
		return p.parseArrayLiteral()
	case tokens.TokenIdent:
		if p.peekTokenIs(tokens.TokenLparen) || p.peekTokenIs(tokens.TokenDot) {
			return p.parseFunctionCall()
		}
		return nil, errors.NewSyntaxError(fmt.Sprintf("Bare identifier '%s' is not allowed outside of context references or object keys", p.curToken.Literal), p.curToken.Line, p.curToken.Column)
	default:
		return nil, errors.NewSyntaxError(fmt.Sprintf("Unexpected token %s", p.curToken.Literal), p.curToken.Line, p.curToken.Column)
	}
}

func (p *Parser) parseContextExpression() (ast.Expression, error) {
	startToken := p.curToken
	if err := p.nextToken(); err != nil {
		return nil, err
	}
	if p.curTokenIs(tokens.TokenIdent) {
		ident := &expressions.IdentifierExpr{
			Name:   p.curToken.Literal,
			Line:   p.curToken.Line,
			Column: p.curToken.Column,
		}
		ce := &expressions.ContextExpr{
			Ident:  ident,
			Line:   startToken.Line,
			Column: startToken.Column,
		}
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return ce, nil
	} else if p.curTokenIs(tokens.TokenLeftBracket) {
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		expr, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		if !p.curTokenIs(tokens.TokenRightBracket) {
			return nil, errors.NewSyntaxError("Expected RBRACKET in context expression", p.curToken.Line, p.curToken.Column)
		}
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		ce := &expressions.ContextExpr{
			Ident:     nil,
			Subscript: expr,
			Line:      startToken.Line,
			Column:    startToken.Column,
		}
		return ce, nil
	} else {
		ce := &expressions.ContextExpr{
			Ident:     nil,
			Subscript: nil,
			Line:      startToken.Line,
			Column:    startToken.Column,
		}
		return ce, nil
	}
}

func (p *Parser) parseFunctionCall() (ast.Expression, error) {
	var parts []string
	parts = append(parts, p.curToken.Literal)
	startToken := p.curToken

	if err := p.nextToken(); err != nil {
		return nil, err
	}
	for p.curTokenIs(tokens.TokenDot) {
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		if !p.curTokenIs(tokens.TokenIdent) {
			return nil, errors.NewSyntaxError("Expected identifier after dot in function call", p.curToken.Line, p.curToken.Column)
		}
		parts = append(parts, p.curToken.Literal)
		if err := p.nextToken(); err != nil {
			return nil, err
		}
	}
	if !p.curTokenIs(tokens.TokenLparen) {
		return nil, errors.NewSyntaxError("Expected '(' in function call", p.curToken.Line, p.curToken.Column)
	}
	parenToken := p.curToken

	if err := p.nextToken(); err != nil {
		return nil, err
	}

	var args []ast.Expression
	if !p.curTokenIs(tokens.TokenRparen) {
		arg, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		for p.curTokenIs(tokens.TokenComma) {
			if err := p.nextToken(); err != nil {
				return nil, err
			}
			arg, err := p.ParseExpression()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
		}
		if !p.curTokenIs(tokens.TokenRparen) {
			return nil, errors.NewSyntaxError("Expected ')' after arguments in function call", p.curToken.Line, p.curToken.Column)
		}
	}
	if err := p.nextToken(); err != nil {
		return nil, err
	}
	return &expressions.FunctionCallExpr{
		Namespace:   parts,
		Args:        args,
		Line:        startToken.Line,
		Column:      startToken.Column,
		ParenLine:   parenToken.Line,
		ParenColumn: parenToken.Column,
	}, nil
}

func (p *Parser) parseArrayLiteral() (ast.Expression, error) {
	startToken := p.curToken
	var elements []ast.Expression
	if err := p.nextToken(); err != nil {
		return nil, err
	}
	if p.curTokenIs(tokens.TokenRightBracket) {
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return &expressions.ArrayLiteralExpr{
			Elements: elements,
			Line:     startToken.Line,
			Column:   startToken.Column,
		}, nil
	}
	expr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	elements = append(elements, expr)
	for p.curTokenIs(tokens.TokenComma) {
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		expr, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		elements = append(elements, expr)
	}
	if !p.curTokenIs(tokens.TokenRightBracket) {
		return nil, errors.NewSyntaxError("Expected ']' at end of array literal", p.curToken.Line, p.curToken.Column)
	}
	if err := p.nextToken(); err != nil {
		return nil, err
	}
	return &expressions.ArrayLiteralExpr{
		Elements: elements,
		Line:     startToken.Line,
		Column:   startToken.Column,
	}, nil
}

func (p *Parser) parseObjectLiteral() (ast.Expression, error) {
	startToken := p.curToken
	fields := make(map[string]ast.Expression)

	if err := p.nextToken(); err != nil {
		return nil, err
	}

	if p.curTokenIs(tokens.TokenRightCurly) {
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return &expressions.ObjectLiteralExpr{
			Fields: fields,
			Line:   startToken.Line,
			Column: startToken.Column,
		}, nil
	}

	for {
		var key string
		if p.curTokenIs(tokens.TokenIdent) || p.curTokenIs(tokens.TokenString) {
			key = strings.TrimSpace(p.curToken.Literal)
		} else {
			return nil, errors.NewSyntaxError("Expected identifier or string as object key", p.curToken.Line, p.curToken.Column)
		}

		// Check for duplicate key.
		if _, exists := fields[key]; exists {
			return nil, errors.NewSemanticError(fmt.Sprintf("Duplicate key '%s' detected", key), p.curToken.Line, p.curToken.Column)
		}

		if !p.peekTokenIs(tokens.TokenColon) {
			return nil, errors.NewSyntaxError("Expected ':' after object key", p.peekToken.Line, p.peekToken.Column)
		}

		if err := p.nextToken(); err != nil {
			return nil, err
		}
		if err := p.nextToken(); err != nil {
			return nil, err
		}

		valueExpr, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		fields[key] = valueExpr

		if p.curTokenIs(tokens.TokenComma) {
			// Detect trailing comma.
			if p.peekTokenIs(tokens.TokenRightCurly) {
				return nil, errors.NewSyntaxError("Trailing comma not allowed in object literal", p.peekToken.Line, p.peekToken.Column)
			}
			if err := p.nextToken(); err != nil {
				return nil, err
			}
		} else if p.curTokenIs(tokens.TokenRightCurly) {
			break
		} else {
			return nil, errors.NewSyntaxError("Expected ',' or '}' after object field", p.curToken.Line, p.curToken.Column)
		}
	}

	if !p.curTokenIs(tokens.TokenRightCurly) {
		return nil, errors.NewSyntaxError("Expected '}' at end of object literal", p.curToken.Line, p.curToken.Column)
	}

	if err := p.nextToken(); err != nil {
		return nil, err
	}

	return &expressions.ObjectLiteralExpr{
		Fields: fields,
		Line:   startToken.Line,
		Column: startToken.Column,
	}, nil
}

func (p *Parser) curTokenIs(t tokens.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t tokens.TokenType) bool {
	return p.peekToken.Type == t
}
