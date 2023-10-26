package main

import "log"

func calcGridPath(m, n int) int {
	pathMemo := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		pathMemo[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		pathMemo[i][1] = 1
	}

	for i := 1; i <= n; i++ {
		pathMemo[1][i] = 1
	}

	for i := 2; i <= m; i++ {
		for j := 2; j <= n; j++ {
			pathMemo[i][j] = pathMemo[i-1][j] + pathMemo[i][j-1]
		}
	}

	return pathMemo[m][n]
}

func main() {
	log.Println(calcGridPath(75, 19))
}
