package hw02unpackstring

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}
	if unicode.IsDigit(rune(input[0])) {
		return "", ErrInvalidString
	}
	digitByRuneIndexMap := make(map[int]int)
	for _, r := range "0123456789" {
		startIdx := 0
		for {
			idx := strings.IndexRune(input[startIdx:], r)
			if idx < 0 {
				break
			}
			idx += startIdx
			digitByRuneIndexMap[idx], _ = strconv.Atoi(string(r))
			startIdx = idx + utf8.RuneLen(r)
			if startIdx >= len(input) {
				break
			}
		}
	}

	sb := strings.Builder{}
	sb.Grow(len(input))

	digitIndexes := sortMapKeys(digitByRuneIndexMap)
	var startIdx int
	for _, digitIdx := range digitIndexes {
		if _, ok := digitByRuneIndexMap[digitIdx-1]; ok {
			return "", ErrInvalidString
		}
		partBetweenDigits := input[startIdx:digitIdx]
		lastRune, size := utf8.DecodeLastRuneInString(partBetweenDigits)
		partBetweenDigitsWithoutLastRune := partBetweenDigits[:len(partBetweenDigits)-size]
		sb.WriteString(partBetweenDigitsWithoutLastRune)
		multiplyCount := digitByRuneIndexMap[digitIdx]
		if multiplyCount > 0 {
			sb.WriteString(strings.Repeat(string(lastRune), multiplyCount))
		}
		startIdx = digitIdx + 1
	}

	if startIdx < len(input) {
		sb.WriteString(input[startIdx:])
	}

	return sb.String(), nil
}

func sortMapKeys(input map[int]int) []int {
	keys := make([]int, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	return keys
}
