package main

import (
	"bufio"
	"fmt"
	"os"
)

func longestPalindrome(s string) string {
	if len(s) < 2 {
		return s
	}

	palindromeMemo := make([][]bool, len(s))

	for i := 0; i < len(s); i++ {
		palindromeMemo[i] = make([]bool, len(s))
		palindromeMemo[i][i] = true
	}

	startIndex, maxLen := 0, 0
	for i := 1; i < len(s); i++ {
		for j := 0; j < len(s)-i; j++ {
			if (i == 1 || i == 2) && s[j] == s[j+i] {
				palindromeMemo[j][j+i] = true
			} else if palindromeMemo[j+1][j+i-1] && s[j] == s[j+i] {
				palindromeMemo[j][j+i] = true
			}

			if palindromeMemo[j][j+i] && i > maxLen {
				maxLen = i
				startIndex = j
			}
		}
	}

	return s[startIndex : startIndex+maxLen+1]
}

func main() {
	r := bufio.NewReader(os.Stdin)

	s, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}

	fmt.Println(longestPalindrome(s))
}
