package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type AnalyzedWord struct {
	Original   string
	Remaining  string
	Components []string
}

var (
	wordfile    = flag.String("wordlist", "word.list", "path to the file with analyzed")
	wordmap     = make(map[string]int)
	words       []string
	queue       []AnalyzedWord
	longestword AnalyzedWord
)

func ReadAllWords(filepath string) (words []string, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, strings.ToLower(scanner.Text()))
	}
	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return words, nil
}

func main() {
	flag.Parse()

	fmt.Printf("Start reading word file %v at %v\n", *wordfile, time.Now())
	words, err := ReadAllWords(*wordfile)
	fmt.Printf("Stop reading word file %v at %v\n", *wordfile, time.Now())
	if err != nil {
		fmt.Println("Error during reading word file %v", err)
		return
	}
	fmt.Printf("Read %v words\n", len(words))
	fmt.Printf("Start forming word map and filling queue at %v\n", time.Now())
	for _, word := range words {
		wordmap[word] = len([]rune(word))
		queue = append(queue, AnalyzedWord{Original: word, Remaining: word})
	}
	fmt.Printf("Stop forming word map and filling queue at %v\n", time.Now())
	fmt.Printf("Start analyzing words %v\n", time.Now())
	step := 0
	for len(queue) != 0 {
		step++
		fmt.Printf("Current queue length %v, step %v\n", len(queue), step)
		temporalqueue := []AnalyzedWord{}
		for _, analyzedword := range queue {
			subword := ""
			for _, symbol := range []rune(analyzedword.Remaining) {
				subword += string(symbol)
				if wordmap[subword] == 0 {
					continue
				}
				var processedword AnalyzedWord
				processedword.Original = analyzedword.Original
				processedword.Components = append(analyzedword.Components, subword)
				processedword.Remaining = strings.TrimPrefix(analyzedword.Remaining, subword)
				if subword == analyzedword.Remaining && len(processedword.Components) > 1 {
					fmt.Printf("Compound word %v\n", processedword)
					if len([]rune(processedword.Original)) > len([]rune(longestword.Original)) {
						longestword = processedword
					}
					break
				}
				temporalqueue = append(temporalqueue, processedword)
			}
		}
		queue = temporalqueue
	}
	fmt.Printf("Stop analyzing words at %v\n", time.Now())
	fmt.Printf("Longest compound word: %v\n", longestword)
}
