package main

import (
	"fmt"

	fqparser "github.com/nht1206/go-study/fq/parser"
)

func main() {
	exprGroups, err := fqparser.Parse(`@request.header != "abc" && (status = "active" || status = "pending")`)
	if err != nil {
		panic(err)
	}

	printExpr(exprGroups)
	fmt.Println("")
}

func printExpr(exprGroups []fqparser.ExprGroup) {
	for i, exprGroup := range exprGroups {
		expr, ok := exprGroup.Item.(fqparser.Expr)
		if ok {
			if i != 0 {
				fmt.Printf(" %v ", exprGroup.Op)
			}
			fmt.Printf("%v %v %v", expr.Left.Literal, expr.Op, expr.Right.Literal)
			continue
		}

		fmt.Printf(" %s (", exprGroup.Op)
		printExpr(exprGroup.Item.([]fqparser.ExprGroup))
		fmt.Print(")")
	}
}
