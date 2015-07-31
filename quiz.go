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
	wordfile  = flag.String("wordlist", "word.list", "path to the file with analyzed")
	wordmap   = make(map[string]int)
	words     []string
	wordqueue []AnalyzedWord
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

func AnalyzeWords(queue []AnalyzedWord) {
	var longestword AnalyzedWord
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
				processedword := new(AnalyzedWord)
				processedword.Original = analyzedword.Original
				for _, component := range analyzedword.Components {
					processedword.Components = append(processedword.Components, component)
				}
				processedword.Components = append(processedword.Components, subword)
				processedword.Remaining = strings.TrimPrefix(analyzedword.Remaining, subword)
				if subword == analyzedword.Remaining {
					fmt.Printf("Compound word %v\n", processedword)
					if len([]rune(processedword.Original)) > len([]rune(longestword.Original)) && len(processedword.Components) > 1 {
						longestword = *processedword
					}
					break
				}
				temporalqueue = append(temporalqueue, *processedword)
			}
		}
		queue = temporalqueue
	}
	fmt.Printf("Stop analyzing words at %v\n", time.Now())
	fmt.Printf("Longest compound word: %v\n", longestword)
}

func main() {
	flag.Parse()

	fmt.Printf("Start reading word file %v at %v\n", *wordfile, time.Now())
	words, err := ReadAllWords(*wordfile)
	fmt.Printf("Stop reading word file %v at %v\n", *wordfile, time.Now())
	if err != nil {
		fmt.Printf("Error during reading word file %v\n", err)
		return
	}
	fmt.Printf("Read %v words\n", len(words))
	fmt.Printf("Start forming word map and filling queue at %v\n", time.Now())
	for _, word := range words {
		wordmap[word] = len([]rune(word))
		wordqueue = append(wordqueue, AnalyzedWord{Original: word, Remaining: word})
	}
	fmt.Printf("Stop forming word map and filling queue at %v\n", time.Now())
	AnalyzeWords(wordqueue)

	for {
		input := ""
		fmt.Println("Enter new word to the list, enter ! to find a longest compound word, enter . to exit")
		fmt.Scanf("%s", &input)
		switch input {
		case "!":
			AnalyzeWords(wordqueue)
		case ".":
			return
		default:
			if input != "" {
				input = strings.ToLower(input)
				wordmap[input] = len([]rune(input))
				wordqueue = append(wordqueue, AnalyzedWord{Original: input, Remaining: input})
			}
		}
	}
}
