package main

import ("fmt"; "strings"; "unicode")

func PalindromeChecker(input string) bool {
	s := removePunctuationAndSpaces(input)
	N := len(s)

	for i := 0; i < (N / 2); i++{
		if s[i] != s[N - i - 1]{
			return false
		}
	}

	return true

}

func removePunctuationAndSpaces(s string) string {
	var builder strings.Builder
	for _, char := range s {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			builder.WriteRune(unicode.ToLower(char))
		}
	}
	return builder.String()
}

func main() {
	arr := []string{
		"Madam, in Eden, I'm Adam",
		"A man, a plan, a canal, Panama!",
		"Hello, World!",
		"Was it a car or a cat I saw?",
		"dad",
	}

	for _, CurrString := range arr {
		fmt.Printf("'%s' is a palindrome: %v\n", CurrString, PalindromeChecker(CurrString))
	}
}