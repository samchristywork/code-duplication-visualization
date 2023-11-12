package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
)

type Sample struct {
	line      string
	neighbors []string
}

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

func getFileContents(filename string) string {
	//fmt.Printf("Reading file %s\n", filename)

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	var contents string
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf[:])
		if err != nil {
			break
		}
		contents += string(buf[0:n])
	}

	return contents
}

func checkForStringInFile(filename string, needle string, threshold float64) bool {
	haystack := getFileContents(filename)

	alines := strings.Split(haystack, "\n")
	blines := strings.Split(needle, "\n")
	needleLen := len(blines)

	additions, _ := diff(alines, blines)

	//fmt.Printf("File %s ", filename)
	//fmt.Printf("additions: %d, deletions: %d\n", additions, deletions)

	return float64(additions) < float64(needleLen)*threshold
}

func checkForStringInDirectory(dirname string, needle string, threshold float64, ignore []string) bool {
	files, err := os.ReadDir(dirname)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		filename := dirname + "/" + file.Name()

		ignoreFlag := false
		for _, ignoreFile := range ignore {
			if filename == ignoreFile {
				ignoreFlag = true
			}
		}

		if ignoreFlag {
			continue
		}

		if file.IsDir() {
			if checkForStringInDirectory(filename, needle, threshold, ignore) {
				return true
			}
		} else {
			if checkForStringInFile(filename, needle, threshold) {
				return true
			}
		}
	}

	return false
}

func samplesFromFile(filename string) []Sample {
	content := getFileContents(filename)
	lines := strings.Split(content, "\n")
	samples := make([]Sample, len(lines))

	for i := 0; i < len(lines); i++ {
		sample := Sample{
			line:      lines[i],
			neighbors: make([]string, 9),
		}

		for j := -4; j <= 4; j++ {
			if i+j < 0 || i+j >= len(lines) {
				sample.neighbors[j+4] = ""
			} else {
				sample.neighbors[j+4] = lines[i+j]
			}
		}

		samples[i] = sample
	}

	return samples
}

func scanFile(filename string) {
	samples := samplesFromFile(filename)
	for _, sample := range samples {
		neighbors := strings.Join(sample.neighbors, "\n")
		ret := checkForStringInDirectory("target", neighbors, 0.3, []string{filename})
		fmt.Printf("%v\t%s\n", ret, sample.line)
	}
}

func createPNGFromSource(filename string) {
	content := getFileContents(filename)
	lines := strings.Split(content, "\n")

	height := len(lines)
	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	background := color.RGBA{0, 0, 0, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, background)
		}
	}

	for y, line := range lines {
		for x, char := range line {
			if char == ' ' || char == '\t' {
				img.Set(x, y, color.RGBA{0, 0, 0, 255})
			} else {
				img.Set(x, y, color.RGBA{255, 255, 255, 255})
			}
		}
	}

	f, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	png.Encode(f, img)

	fmt.Printf("Wrote image to target/computer.png\n")
}

func main() {
	scanFile("target/computer.go")
}
