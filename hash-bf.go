package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"strings"
	"time"
)

var letters = [62]string{
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

func main() {
	start := time.Now()
	target := flag.String("target", "", "Result of hash that we are looking for.")
	workersCount := flag.Int("w", 1, "Num of workers")
	flag.Parse()
	if *target == "" {
		panic("No hash provided.")
	}

	inputChan := make(chan string, *workersCount*2)
	outputChan := make(chan string)
	go fillChan(inputChan)
	for i := 0; i <= *workersCount; i++ {
		go hasher(inputChan, outputChan, *target)
	}

	result := <-outputChan
	fmt.Println(result, "\n", time.Since(start))
}

func fillChan(inputChan chan string) {
	s := letters[0]
	for {
		r := increment(s)
		inputChan <- r
		s = r
	}
}

func hasher(inputChan, outputChan chan string, target string) {
	for s := range inputChan {
		sum := sha256.Sum256([]byte(s))
		r := fmt.Sprintf("%x", sum)
		if r == target {
			outputChan <- s
		}
	}
}

func increment(source string) string {
	sourceSlice := strings.Split(source, "")
Loop:
	for lastIndex := len(source) - 1; ; lastIndex-- {
		if lastIndex == -1 {
			sourceSlice = append([]string{letters[0]}, sourceSlice...)
			break
		}
		last := sourceSlice[lastIndex]
		for i, l := range letters {
			if l == last {
				nextI := i + 1
				if nextI < len(letters) {
					sourceSlice[lastIndex] = letters[nextI]
				} else {
					sourceSlice[lastIndex] = letters[0]
					break
				}
				break Loop
			}
		}
	}
	return strings.Join(sourceSlice, "")
}
