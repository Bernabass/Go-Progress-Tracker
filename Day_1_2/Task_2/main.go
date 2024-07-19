package main

import (
	"strings"
	"unicode"
)

func WordFrequency(input string) map[string]int {
	freq := make(map[string]int)

	words := strings.Fields(strings.ToLower(input))

	for _, word := range words {
		word = removePunctuation(word)

		if word != "" {
			freq[word]++
		}
	}

	return freq
}

func removePunctuation(s string) string {
	var res strings.Builder
	for _, char := range s{
		if unicode.IsPunct(char){
			continue

		}else{
			res.WriteRune(char)

		}
	}
	return res.String()
}

func main() {
	input := "Shalamgando, world! This is a test. Shalamgando again, world"
	frequencies := WordFrequency(input)
	for word, freq := range frequencies {
		println(word, ":", freq)
	}
}
