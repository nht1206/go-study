package qb

import "strings"

type Builder interface {
	QuoteSimpleTableName(tableName string) string
	QuoteSimpleColumnName(colName string) string
	GeneratePlaceholder(int) string
	QueryBuilder() QueryBuilder
	NewQuery(sql string, params Params) *Query
	Select(cols ...string) *SelectQuery
}

type BaseBuilder struct {
	db       *DB
	executor Executor
}

func NewBaseBuilder(db *DB, executor Executor) *BaseBuilder {
	return &BaseBuilder{
		db:       db,
		executor: executor,
	}
}

func (b *BaseBuilder) QuoteSimpleTableName(tableName string) string {
	if strings.Contains(tableName, `"`) {
		return tableName
	}

	return `"` + tableName + `"`
}

func (b *BaseBuilder) QuoteSimpleColumnName(colName string) string {
	if colName == "*" || strings.Contains(colName, `"`) {
		return colName
	}
	return `"` + colName + `"`
}

func (b *BaseBuilder) GeneratePlaceholder(int) string {
	return "?"
}

func (b *BaseBuilder) NewQuery(sql string, params Params) *Query {
	return NewQuery(b.db, b.executor, sql, params)
}

func (b *BaseBuilder) Select(cols ...string) *SelectQuery {
	return NewSelectQuery(b.db).Select(cols...)
}

type StandardBuilder struct {
	*BaseBuilder
	qb QueryBuilder
}

func NewStandardBuilder(db *DB) *StandardBuilder {
	return &StandardBuilder{
		BaseBuilder: NewBaseBuilder(db, db.db),
		qb:          NewBaseQueryBuilder(db),
	}
}

func (b *StandardBuilder) QueryBuilder() QueryBuilder {
	return b.qb
}
