package parser

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"

	"github.com/masahiro331/go-license/lexer"
	"github.com/masahiro331/go-license/token"
)

var (
	ErrInvalidExpression = xerrors.New("invalid expression error")
)

type Pair struct {
	root    *LicenseExpression
	cursor  *LicenseExpression
	bracket token.TokenType
}

type Stack []Pair

func (s *Stack) Push(x Pair) {
	*s = append(*s, x)
}

func (s *Stack) Pop() Pair {
	l := len(*s)
	x := (*s)[l-1]
	*s = (*s)[:l-1]
	return x
}

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

type LicenseExpression struct {
	Node     Node
	Operator string
	Next     *LicenseExpression
}

type Node struct {
	License           string
	LicenseExpression *LicenseExpression
}

func (l *LicenseExpression) String() string {
	cursor := l

	var str string
	for ; cursor != nil; cursor = cursor.Next {
		str = strings.Join([]string{str, cursor.Node.String(), cursor.Operator}, " ")
	}
	return strings.TrimSpace(str)
}

func (n Node) String() string {
	if n.LicenseExpression != nil {
		return fmt.Sprintf("(%s)", n.LicenseExpression)
	}
	return n.License
}

func Parse(lex *lexer.Lexer) (*LicenseExpression, error) {
	root := &LicenseExpression{}
	cursor := root
	stack := Stack{}

	for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
		switch tok.Type {
		case token.IDENT:
			if cursor.Node.License == "" {
				cursor.Node = Node{License: tok.Literal}
			} else {
				cursor.Node.License = fmt.Sprintf("%s %s", cursor.Node.License, tok.Literal)
			}
		case token.AND, token.OR:
			cursor.Operator = string(tok.Type)
			cursor.Next = &LicenseExpression{}
			cursor = cursor.Next
		case token.LPAREN, token.LBRACE:
			p := Pair{root: root, cursor: cursor, bracket: tok.Type}
			stack.Push(p)
			root = &LicenseExpression{}
			cursor = root
		case token.RPAREN, token.RBRACE:
			e := stack.Pop()
			if e.bracket == token.LPAREN {
				if tok.Type != token.RPAREN {
					return nil, ErrInvalidExpression
				}
			} else if e.bracket == token.LBRACE {
				if tok.Type != token.RBRACE {
					return nil, ErrInvalidExpression
				}
			}
			e.cursor.Node.LicenseExpression = root
			cursor = e.cursor
			root = e.root
		}
	}
	if !stack.IsEmpty() {
		return nil, ErrInvalidExpression
	}
	return root, nil
}
