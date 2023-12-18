package qb

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

type (
	BuilderNewFunc func(*DB, Executor) Builder
)

var BuilderNewFuncMap map[string]BuilderNewFunc

type DB struct {
	Builder
	db         *sql.DB
	driverName string
}

func FromDB(sqlDB *sql.DB, driverName string) *DB {
	db := &DB{
		db:         sqlDB,
		driverName: driverName,
	}

	db.Builder = db.newBuilder()

	return db
}

func (d *DB) DB() *sql.DB {
	return d.db
}

func (d *DB) Close() error {
	if d.db != nil {
		return d.db.Close()
	}

	return nil
}

func (d *DB) newBuilder() Builder {
	f, ok := BuilderNewFuncMap[d.driverName]
	if ok {
		return f(d, d.db)
	}

	return NewStandardBuilder(d)
}

func (d *DB) QuoteColumnName(colName string) string {
	if strings.Contains(colName, "(") || strings.Contains(colName, "[[") {
		return colName
	}

	if !strings.Contains(colName, ".") {
		return d.QuoteSimpleColumnName(colName)
	}

	tableName := ""
	if pos := strings.LastIndex(colName, "."); pos != -1 {
		tableName = d.QuoteTableName(colName[:pos])
		colName = colName[pos+1:]
	}

	return fmt.Sprintf("%s.%s", tableName, d.QuoteSimpleColumnName(colName))
}

func (d *DB) QuoteTableName(tableName string) string {
	if strings.Contains(tableName, "{{") || strings.Contains(tableName, "[[") {
		return tableName
	}

	if !strings.Contains(tableName, ".") {
		return d.QuoteSimpleTableName(tableName)
	}

	parts := strings.Split(tableName, ".")
	for i, p := range parts {
		parts[i] = d.QuoteSimpleTableName(p)
	}

	return strings.Join(parts, ".")
}

var (
	plRegex  = regexp.MustCompile(`\{:\w+\}`)
	colRegex = regexp.MustCompile(`(\{\{[\w\.\- ]+\}\}|\[\[[\w\.\- ]+\]\])`)
)

func (d *DB) processSQL(sql string) (string, []string) {
	placeholders := []string{}

	count := 0
	sql = plRegex.ReplaceAllStringFunc(sql, func(s string) string {
		count++
		placeholders = append(placeholders, s[2:len(s)-1])

		return d.GeneratePlaceholder(count)
	})

	sql = colRegex.ReplaceAllStringFunc(sql, func(s string) string {
		return d.QuoteColumnName(s[2 : len(s)-2])
	})

	return sql, placeholders
}
