package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(input string) []string {
	topCount := 10

	words := strings.Fields(input)
	if len(words) == 0 {
		return []string{}
	}

	countByWord := make(map[string]int, len(words))
	for _, w := range words {
		countByWord[w]++
	}

	wordsSlice := make([]string, 0, len(countByWord))
	for k := range countByWord {
		wordsSlice = append(wordsSlice, k)
	}

	sort.Slice(wordsSlice, func(i int, j int) bool {
		a, b := wordsSlice[i], wordsSlice[j]
		if countByWord[a] == countByWord[b] {
			return a < b
		}
		return countByWord[a] > countByWord[b]
	})

	if len(wordsSlice) < topCount {
		topCount = len(wordsSlice)
	}
	return wordsSlice[:topCount]
}
