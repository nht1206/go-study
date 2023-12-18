package main

import (
	"database/sql"
	"fmt"

	"github.com/nht1206/go-study/qb"
)

func main() {
	db := qb.FromDB(&sql.DB{}, "test")

	fmt.Println(db.Select("a", "b").From("{{test}}").Where("test.id = {:hihi}").Limit(1000).Build().SQL())
}
