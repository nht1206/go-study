package sql_test

import (
	"testing"

	"github.com/nht1206/go-study/fq/sql"
)

func TestQuery(t *testing.T) {
	testcases := []struct {
		name         string
		filter       string
		selectFields []string
		tableName    string
		want         string
	}{
		{
			"simple filter query expression",
			"[a > 0 || b = 1] && c != 0",
			[]string{"a", "b"},
			"test_table",
			"SELECT a, b FROM `test_table` WHERE (a > 0 OR b = 1) AND c <> 0",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := sql.NewQuery(tc.tableName, tc.filter)
			if err != nil {
				t.Fatal(err)
			}

			if err := q.Parse(); err != nil {
				t.Fatal(err)
			}
			if got := q.Select(tc.selectFields...).SQL(); got != tc.want {
				t.Errorf("Query.SQL() = %v, want %v", got, tc.want)
			}
		})
	}
}
