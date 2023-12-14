package sql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nht1206/go-study/fq/parser"
	"github.com/nht1206/go-study/fq/tokenizer"
)

var JoinOpMap = map[parser.JoinOp]string{
	parser.JoinAnd: "AND",
	parser.JoinOr:  "OR",
}

var SignOpMap = map[parser.SignOp]string{
	parser.SignEq:  "=",
	parser.SignNeq: "<>",
	parser.SignGt:  ">",
	parser.SignLt:  "<",
	parser.SignGe:  ">=",
	parser.SignLe:  "<=",
}

type Query struct {
	filter       string
	selectFields []string
	tableName    string
	whereExpr    string
}

func NewQuery(tableName string, filter string) (*Query, error) {
	if tableName == "" {
		return nil, errors.New("the table name must not be empty")
	}

	if filter == "" {
		return nil, errors.New("the filter query must not be empty")
	}

	return &Query{
		filter:    filter,
		tableName: tableName,
	}, nil
}

func (q *Query) Select(fieldNames ...string) *Query {
	q.selectFields = append(q.selectFields, fieldNames...)
	return q
}

func (q *Query) Parse() error {
	exprGroups, err := parser.Parse(q.filter)
	if err != nil {
		return fmt.Errorf("parse filter query failed: %w", err)
	}

	q.whereExpr = parseExprGroup(exprGroups)

	return nil
}

func parseExprGroup(exprGroups []parser.ExprGroup) string {
	whereExprs := strings.Builder{}
	for i, exprGroup := range exprGroups {
		if i > 0 {
			whereExprs.WriteString(fmt.Sprintf(" %s ", JoinOpMap[exprGroup.Op]))
		}
		switch item := exprGroup.Item.(type) {
		case parser.Expr:
			expr := convertToSQLExpr(item)
			whereExprs.WriteString(expr)
		case []parser.ExprGroup:
			expr := parseExprGroup(item)
			whereExprs.WriteString(fmt.Sprintf("(%s)", expr))
		}
	}

	return whereExprs.String()
}

func convertToSQLExpr(expr parser.Expr) string {
	lResult := resolveToken(expr.Left)
	rResult := resolveToken(expr.Right)

	return fmt.Sprintf("%s %s %s", lResult, SignOpMap[expr.Op], rResult)
}

func resolveToken(token *tokenizer.Token) string {
	switch token.Type {
	case tokenizer.TokenIdentifier:
		return token.Literal
	case tokenizer.TokenText, tokenizer.TokenNumber:
		return token.Literal
	default:
		return ""
	}
}

func (q *Query) Where() string {
	return q.whereExpr
}

func (q *Query) SQL() string {
	sqlQuery := strings.Builder{}
	sqlQuery.WriteString("SELECT ")
	if len(q.selectFields) > 0 {
		sqlQuery.WriteString(strings.Join(q.selectFields, ", ") + " ")
	} else {
		sqlQuery.WriteString("* ")
	}
	sqlQuery.WriteString(fmt.Sprintf("FROM `%s` ", q.tableName))

	where := q.Where()
	if where != "" {
		sqlQuery.WriteString("WHERE " + q.Where())
	}

	return sqlQuery.String()
}
