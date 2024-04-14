package hw02unpackstring

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}
	if unicode.IsDigit(rune(input[0])) {
		return "", ErrInvalidString
	}
	inputRunes := []rune(input)
	digitByRuneIndexMap := make(map[int]int)
	for _, r := range "0123456789" {
		startIdx := 0
		for {
			digitIdx := strings.IndexRune(string(inputRunes[startIdx:]), r)
			if digitIdx < 0 {
				break
			}
			digitIdx += startIdx
			digitByRuneIndexMap[digitIdx], _ = strconv.Atoi(string(r))
			startIdx = digitIdx + 1
			if startIdx >= len(inputRunes) {
				break
			}
		}
	}

	digitIndexes := sortMapKeys(digitByRuneIndexMap)
	var resultString string
	var startIdx int
	for _, digitIdx := range digitIndexes {
		if _, ok := digitByRuneIndexMap[digitIdx-1]; ok {
			return "", ErrInvalidString
		}
		partBetweenDigits := inputRunes[startIdx:digitIdx]
		if len(partBetweenDigits) > 1 {
			partBetweenDigitsWithoutLastRune := partBetweenDigits[:len(partBetweenDigits)-1]
			resultString += string(partBetweenDigitsWithoutLastRune)
		}
		lastRune := partBetweenDigits[len(partBetweenDigits)-1:]
		multiplyCount := digitByRuneIndexMap[digitIdx]
		if multiplyCount > 0 {
			resultString += strings.Repeat(string(lastRune), multiplyCount)
		}
		startIdx = digitIdx + 1
	}

	if startIdx < len(inputRunes) {
		resultString += string(inputRunes[startIdx:])
	}

	return resultString, nil
}

func sortMapKeys(input map[int]int) []int {
	keys := make([]int, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	return keys
}
