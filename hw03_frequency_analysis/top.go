package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(input string) []string {
	topCount := 10

	words := strings.Fields(input)
	countByWord := make(map[string]int, len(words))
	for _, w := range words {
		countByWord[w]++
	}

	wordCounts := make([]int, 0, len(countByWord))
	wordsByCount := make(map[int][]string, len(countByWord))
	for k, v := range countByWord {
		wordsByCount[v] = append(wordsByCount[v], k)
		if len(wordsByCount[v]) == 1 {
			wordCounts = append(wordCounts, v)
		}
	}

	sort.Ints(wordCounts)
	result := make([]string, 0, topCount)
	counter := 1
	for i := len(wordCounts) - 1; i >= 0; i-- {
		if counter > topCount {
			break
		}
		sort.Strings(wordsByCount[wordCounts[i]])
		for _, w := range wordsByCount[wordCounts[i]] {
			result = append(result, w)
			counter++
			if counter > topCount {
				break
			}
		}
	}

	return result
}
