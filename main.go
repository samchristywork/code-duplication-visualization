package main

import (
	"fmt"
)

func longestCommonSubsequence(a, b []string) [][]int {
	lenA, lenB := len(a), len(b)
	dp := make([][]int, lenA+1)
	for i := range dp {
		dp[i] = make([]int, lenB+1)
	}

	for i := 1; i <= lenA; i++ {
		for j := 1; j <= lenB; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	return dp
}

func main() {
	filename := "target/computer.go"
	needle := `TestString`

	ret := checkForStringInFile(filename, needle)
	fmt.Printf("ret: %v\n", ret)
}
