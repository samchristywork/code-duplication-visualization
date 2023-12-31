package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"os"
	"strings"
)

type Sample struct {
	line      string
	neighbors []string
}

type Line struct {
	Line    string
	Color   color.RGBA
	Tooltip string
}

type Match struct {
	Filename string
	Line     int
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func normalize(s string) string {
	s = strings.ReplaceAll(s, "\t", "")
	return s
}

func compare(a, b string) bool {
	return normalize(a) == normalize(b)
}

func longestCommonSubsequence(a, b []string) [][]int {
	lenA, lenB := len(a), len(b)
	dp := make([][]int, lenA+1)
	for i := range dp {
		dp[i] = make([]int, lenB+1)
	}

	for i := 1; i <= lenA; i++ {
		for j := 1; j <= lenB; j++ {
			if compare(a[i-1], b[j-1]) {
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

	return float64(additions) < float64(needleLen)*threshold
}

func checkForStringInDirectory(dirname string, needle string, threshold float64, ignore []string) []string {
	matches := []string{}

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
			files := checkForStringInDirectory(filename, needle, threshold, ignore)
			matches = append(matches, files...)
		} else {
			if checkForStringInFile(filename, needle, threshold) {
				matches = append(matches, filename)
			}
		}
	}

	return matches
}

func findMatchesInFile(filename string, needle string, threshold float64) []Match {
	matches := []Match{}

	lines := strings.Split(getFileContents(filename), "\n")
	n := strings.Split(needle, "\n")

	for i := 0; i < len(lines)-len(n); i++ {
		score := 0
		for j := 0; j < len(n); j++ {
			if compare(lines[i+j], n[j]) {
				score++
			}
		}

		if score >= int(float64(len(n))*threshold) {
			matches = append(matches, Match{filename, i})
		}
	}

	return matches
}

func findMatches(dirname string, needle string, threshold float64, ignore []string) []Match {
	matches := []Match{}

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
			matches = append(matches, findMatches(filename, needle, threshold, ignore)...)
		} else {
			m := findMatchesInFile(filename, needle, threshold)
			matches = append(matches, m...)
		}
	}

	return matches
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

		windowSize := 4

		for j := -windowSize; j <= windowSize; j++ {
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

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s %s\n", r.RemoteAddr, r.Proto, r.Method, r.URL)

		next.ServeHTTP(w, r)
	})
}

func getFiles(dirname string, extension string) []string {
	files, err := os.ReadDir(dirname)
	if err != nil {
		panic(err)
	}

	var ret []string

	for _, file := range files {
		filename := dirname + "/" + file.Name()

		if file.IsDir() {
			ret = append(ret, getFiles(filename, extension)...)
		} else if strings.HasSuffix(filename, extension) {
			ret = append(ret, filename)
		}
	}

	return ret
}

func startServer(port int) {
	http.Handle("/file", middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("file")

		samples := samplesFromFile(filename)

		w.Header().Set("Content-Type", "application/json")

		lines := []Line{}

		for n, sample := range samples {
			matches := findMatches("target", strings.Join(sample.neighbors, "\n"), 0.7, []string{filename})

			line := Line{
				Line:    sample.line,
				Tooltip: "",
			}

			if len(matches) >= 3 {
				line.Color = color.RGBA{255, 0, 0, 255} // Red
			} else if len(matches) == 2 {
				line.Color = color.RGBA{255, 128, 0, 255} // Orange
			} else if len(matches) == 1 {
				line.Color = color.RGBA{255, 255, 0, 255} // Yellow
			} else {
				line.Color = color.RGBA{0, 255, 0, 255} // Green
			}

			if len(matches) > 0 {
				line.Tooltip += fmt.Sprintf("%s:%d\n\n", filename, n)

				for _, l := range sample.neighbors {
					line.Tooltip += l + "\n"
				}
				line.Tooltip += "\n"

				for _, m := range matches {
					line.Tooltip += fmt.Sprintf("%s (%d)\n", m.Filename, m.Line)
				}
			}

			lines = append(lines, line)
		}

		bytes, err := json.Marshal(lines)
		if err != nil {
			panic(err)
		}

		w.Write(bytes)
	})))

	http.Handle("/files", middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		files := getFiles("target", ".go")

		w.Header().Set("Content-Type", "application/json")

		bytes, err := json.Marshal(files)
		if err != nil {
			panic(err)
		}

		w.Write(bytes)
	})))

	http.Handle("/", middleware(http.FileServer(http.Dir("./static"))))

	fmt.Printf("Starting server on port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	startServer(8000)
}
