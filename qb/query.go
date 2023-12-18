package qb

import (
	"context"
	"database/sql"
)

type Params map[string]any

type Executor interface {
	// Exec executes a query and returns the result.
	Exec(query string, args ...any) (sql.Result, error)
	// ExecContext executes a query with the given context and returns the result.
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	// Prepare creates a prepared statement for later queries or executions.
	Prepare(query string) (*sql.Stmt, error)
}

type Query struct {
	db       *DB
	executor Executor

	sql, rawSQL  string
	placeholders []string
	params       Params
	stmt         *sql.Stmt
	Err          error
}

func NewQuery(db *DB, executor Executor, sql string, params Params) *Query {
	rawSQL, placeholders := db.processSQL(sql)

	return &Query{
		db:           db,
		executor:     executor,
		sql:          sql,
		rawSQL:       rawSQL,
		placeholders: placeholders,
		params:       params,
	}
}

func (q *Query) Prepare() *Query {
	stmt, err := q.executor.Prepare(q.rawSQL)
	if err != nil {
		q.Err = err
		return q
	}

	q.stmt = stmt

	return q
}

func (q *Query) ExecContext(ctx context.Context) (sql.Result, error) {
	args := q.prepareArgs()
	if q.stmt == nil {
		return q.executor.ExecContext(ctx, q.sql, args)
	} else {
		return q.stmt.ExecContext(ctx, args)
	}
}

func (q *Query) Exec() (sql.Result, error) {
	args := q.prepareArgs()
	if q.stmt == nil {
		return q.executor.Exec(q.rawSQL, args)
	} else {
		return q.stmt.Exec(args)
	}
}

func (q *Query) prepareArgs() []any {
	result := []any{}

	for _, p := range q.placeholders {
		pValue, ok := q.params[p]
		if ok {
			result = append(result, pValue)
		}
	}

	return result
}

func (q *Query) SQL() string {
	return q.rawSQL
}
