package main

import (
	"fmt"
	"strings"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

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

func diff(a, b []string) (int, int) {
	dp := longestCommonSubsequence(a, b)
	i, j := len(a), len(b)

	additions := 0
	deletions := 0

	for i > 0 && j > 0 {
		if a[i-1] == b[j-1] {
			// fmt.Printf("  %s\n", a[i-1])
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			// fmt.Printf("- %s\n", a[i-1])
			i--
			deletions++
		} else {
			// fmt.Printf("+ %s\n", b[j-1])
			j--
			additions++
		}
	}

	for i > 0 {
		// fmt.Printf("- %s\n", a[i-1])
		i--
		deletions++
	}

	for j > 0 {
		// fmt.Printf("+ %s\n", b[j-1])
		j--
		additions++
	}

	return additions, deletions
}

func main() {
	filename := "target/computer.go"
	needle := `TestString`

	ret := checkForStringInFile(filename, needle)
	fmt.Printf("ret: %v\n", ret)
}
