package qb

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type JoinInfo struct {
	Type      string
	TableName string
	OnExpr    string
}

type QueryBuilder interface {
	BuildSelect(distinct bool, cols []string) string
	BuildFrom(tableName string) string
	BuildJoin(joins []JoinInfo) string
	BuildWhere(expr string) string
	BuildGroupBy(cols []string) string
	BuildHaving(havingExpr string) string
	BuildOrderBy(cols []string) string
	BuildLimit(limit int) string
}

type BaseQueryBuilder struct {
	db *DB
}

func NewBaseQueryBuilder(db *DB) *BaseQueryBuilder {
	return &BaseQueryBuilder{
		db: db,
	}
}

var aliasRegex = regexp.MustCompile(`(?i:\s+as\s+|\s+)([\w\-_\.]+)$`)

func (b *BaseQueryBuilder) BuildSelect(distinct bool, cols []string) string {
	s := bytes.Buffer{}
	s.WriteString("SELECT ")
	if distinct {
		s.WriteString("DISTINCT ")
	}

	if len(cols) == 0 {
		s.WriteString("*")
		return s.String()
	}

	for i, c := range cols {
		if i > 0 {
			s.WriteString(", ")
		}

		matches := aliasRegex.FindStringSubmatch(c)
		if len(matches) == 0 {
			s.WriteString(b.db.QuoteColumnName(c))
		} else {
			c = c[:len(c)-len(matches[0])]
			alias := matches[1]
			s.WriteString(b.db.QuoteColumnName(c) + " AS " + b.db.QuoteColumnName(alias))
		}
	}

	return s.String()
}

func (b *BaseQueryBuilder) BuildFrom(tableName string) string {
	if tableName == "" {
		return ""
	}

	f := "FROM "

	return f + b.quoteTableNameAndAlias(tableName)
}

func (b *BaseQueryBuilder) BuildJoin(joinInfos []JoinInfo) string {
	if len(joinInfos) == 0 {
		return ""
	}

	joins := []string{}
	for _, joinInfo := range joinInfos {
		joins = append(joins, fmt.Sprintf("%s %s ON %s", joinInfo.Type, b.quoteTableNameAndAlias(joinInfo.TableName), joinInfo.OnExpr))
	}

	return strings.Join(joins, " ")
}

func (b *BaseQueryBuilder) BuildWhere(expr string) string {
	return expr
}

func (b *BaseQueryBuilder) quoteTableNameAndAlias(tableName string) string {
	matches := aliasRegex.FindStringSubmatch(tableName)
	if len(matches) == 0 {
		return b.db.QuoteTableName(tableName)
	}

	tableName = tableName[:len(tableName)-len(matches[0])]
	alias := matches[1]

	return b.db.QuoteTableName(tableName) + " " + b.db.QuoteTableName(alias)
}

func (b *BaseQueryBuilder) BuildGroupBy(cols []string) string {
	if len(cols) == 0 {
		return ""
	}

	gb := strings.Builder{}
	gb.WriteString("GROUP BY ")
	for i, col := range cols {
		if i > 0 {
			gb.WriteString(", ")
		}
		gb.WriteString(b.db.QuoteColumnName(col))
	}

	return gb.String()
}

func (b *BaseQueryBuilder) BuildHaving(havingExpr string) string {
	return havingExpr
}

func (b *BaseQueryBuilder) BuildOrderBy(cols []string) string {
	if len(cols) == 0 {
		return ""
	}

	ob := strings.Builder{}
	ob.WriteString("ORDER BY ")
	for i, col := range cols {
		if i > 0 {
			ob.WriteString(", ")
		}
		ob.WriteString(b.db.QuoteColumnName(col))
	}

	return ob.String()
}

func (b *BaseQueryBuilder) BuildLimit(limit int) string {
	if limit < 0 {
		return ""
	}

	return fmt.Sprintf("LIMIT %d", limit)
}
