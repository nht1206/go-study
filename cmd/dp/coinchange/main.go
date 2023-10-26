package main

import "log"

func calcCoinChangeWays(coins []int, target int) int {
	wayMemo := make([]int, target+1)
	wayMemo[0] = 1

	for i := 0; i < len(coins); i++ {
		for j := 0; j <= target; j++ {
			if coins[i] <= j {
				wayMemo[j] = wayMemo[j-coins[i]] + wayMemo[j]
			}
		}
	}

	return wayMemo[target]
}

func main() {
	log.Println(calcCoinChangeWays([]int{1, 5, 10}, 10))
}
