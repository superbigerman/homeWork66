package service

import "strings"

// MaskData маскирует данные
func MaskData(data []string) []string {
	var result []string

	for _, line := range data {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		words := strings.Fields(line)
		maskedWords := make([]string, len(words))

		for i, word := range words {
			maskedWords[i] = maskWord(word)
		}

		result = append(result, strings.Join(maskedWords, " "))
	}

	return result
}

// maskWord маскирует одно слово
func maskWord(word string) string {
	if len(word) <= 1 {
		return word
	}
	return string(word[0]) + strings.Repeat("*", len(word)-1)
}
