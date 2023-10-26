package main

import "log"

func min(a, b int) int {
	if a > b {
		return b
	}

	return a
}

func calcMinimumCoins(coins []int, target int) int {
	minimunCoinMemo := make([]int, target+1)
	minimunCoinMemo[0] = 0

	for i := 1; i <= target; i++ {
		for _, c := range coins {
			sub := i - c
			if sub < 0 {
				continue
			}
			if minimunCoinMemo[i] == 0 {
				minimunCoinMemo[i] = minimunCoinMemo[sub] + 1
			} else {
				minimunCoinMemo[i] = min(minimunCoinMemo[i], minimunCoinMemo[sub]+1)
			}
		}
	}

	return minimunCoinMemo[target]
}

func main() {
	log.Println(calcMinimumCoins([]int{2, 3, 1}, 100))
}
