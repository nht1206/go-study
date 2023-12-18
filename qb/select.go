package qb

type SelectQuery struct {
	builder Builder

	selects   []string
	distinct  bool
	from      string
	joinInfos []JoinInfo
	where     string
	groupBy   []string
	having    string
	orderBy   []string
	limit     int
	params    Params
}

func NewSelectQuery(builder Builder) *SelectQuery {
	return &SelectQuery{
		builder: builder,
		limit:   -1,
	}
}

func (q *SelectQuery) Select(cols ...string) *SelectQuery {
	q.selects = append(q.selects, cols...)
	return q
}

func (q *SelectQuery) Distinct() *SelectQuery {
	q.distinct = true
	return q
}

func (q *SelectQuery) From(tableName string) *SelectQuery {
	q.from = tableName
	return q
}

func (q *SelectQuery) Join(joinInfo JoinInfo) *SelectQuery {
	q.joinInfos = append(q.joinInfos, joinInfo)
	return q
}

func (q *SelectQuery) InnerJoin(tableName string, onExpr string) *SelectQuery {
	return q.Join(JoinInfo{Type: "INNER JOIN", TableName: tableName, OnExpr: onExpr})
}

func (q *SelectQuery) LeftJoin(tableName string, onExpr string) *SelectQuery {
	return q.Join(JoinInfo{Type: "LEFT JOIN", TableName: tableName, OnExpr: onExpr})
}

func (q *SelectQuery) RightJoin(tableName string, onExpr string) *SelectQuery {
	return q.Join(JoinInfo{Type: "RIGHT JOIN", TableName: tableName, OnExpr: onExpr})
}

func (q *SelectQuery) Where(expr string) *SelectQuery {
	q.where = expr
	return q
}

func (q *SelectQuery) GroupBy(cols ...string) *SelectQuery {
	q.groupBy = append(q.groupBy, cols...)
	return q
}

func (q *SelectQuery) Having(expr string) *SelectQuery {
	q.having = expr
	return q
}

func (q *SelectQuery) OrderBy(cols ...string) *SelectQuery {
	q.orderBy = append(q.orderBy, cols...)
	return q
}

func (q *SelectQuery) Limit(limit int) *SelectQuery {
	q.limit = limit
	return q
}

func (q *SelectQuery) Bind(params Params) *SelectQuery {
	if len(q.params) == 0 {
		q.params = params
	} else {
		for k, v := range params {
			q.params[k] = v
		}
	}

	return q
}

func (q *SelectQuery) Build() *Query {
	qb := q.builder.QueryBuilder()

	clauses := []string{
		qb.BuildSelect(q.distinct, q.selects),
		qb.BuildFrom(q.from),
		qb.BuildJoin(q.joinInfos),
		qb.BuildWhere(q.where),
		qb.BuildGroupBy(q.groupBy),
		qb.BuildHaving(q.having),
		qb.BuildOrderBy(q.orderBy),
		qb.BuildLimit(q.limit),
	}

	sql := ""
	for i, clause := range clauses {
		if clause == "" {
			continue
		}

		if i > 0 {
			sql += " "
		}

		sql += clause
	}

	return q.builder.NewQuery(sql, make(Params))
}
