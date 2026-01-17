package utils

import (
	"crypto/rand"
	"math/big"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// IsEmpty checks if a string is empty
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// IsNotEmpty checks if a string is not empty
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// Truncate truncates a string to the specified length
func Truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length]
}

// TruncateWithEllipsis truncates a string and adds ellipsis
func TruncateWithEllipsis(s string, length int) string {
	if len(s) <= length {
		return s
	}
	if length <= 3 {
		return s[:length]
	}
	return s[:length-3] + "..."
}

// Contains checks if a string contains a substring
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// ContainsAny checks if a string contains any of the substrings
func ContainsAny(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// ContainsAll checks if a string contains all of the substrings
func ContainsAll(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if !strings.Contains(s, substr) {
			return false
		}
	}
	return true
}

// StartsWith checks if a string starts with a prefix
func StartsWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// EndsWith checks if a string ends with a suffix
func EndsWith(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// ToLower converts a string to lowercase
func ToLower(s string) string {
	return strings.ToLower(s)
}

// ToUpper converts a string to uppercase
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// ToTitle converts a string to title case
func ToTitle(s string) string {
	return cases.Title(language.English).String(s)
}

// ToCamelCase converts a string to camelCase
func ToCamelCase(s string) string {
	if IsEmpty(s) {
		return s
	}

	words := strings.Fields(strings.ReplaceAll(s, "_", " "))
	if len(words) == 0 {
		return s
	}

	result := strings.ToLower(words[0])
	for _, word := range words[1:] {
		result += cases.Title(language.English).String(word)
	}

	return result
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(s string) string {
	if IsEmpty(s) {
		return s
	}

	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}

	return string(result)
}

// ToKebabCase converts a string to kebab-case
func ToKebabCase(s string) string {
	return strings.ReplaceAll(ToSnakeCase(s), "_", "-")
}

// RemoveWhitespace removes all whitespace from a string
func RemoveWhitespace(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

// RemoveSpecialChars removes special characters from a string
func RemoveSpecialChars(s string) string {
	var result []rune
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			result = append(result, r)
		}
	}
	return string(result)
}

// Reverse reverses a string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsPalindrome checks if a string is a palindrome
func IsPalindrome(s string) bool {
	s = strings.ToLower(RemoveWhitespace(s))
	return s == Reverse(s)
}

// IsEmail checks if a string is a valid email
func IsEmail(s string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

// IsURL checks if a string is a valid URL
func IsURL(s string) bool {
	pattern := `^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

// IsPhoneNumber checks if a string is a valid phone number
func IsPhoneNumber(s string) bool {
	pattern := `^\+?[1-9]\d{1,14}$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

// IsNumeric checks if a string is numeric
func IsNumeric(s string) bool {
	pattern := `^-?\d+(\.\d+)?$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

// IsAlpha checks if a string contains only alphabetic characters
func IsAlpha(s string) bool {
	pattern := `^[a-zA-Z]+$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

// IsAlphaNumeric checks if a string contains only alphanumeric characters
func IsAlphaNumeric(s string) bool {
	pattern := `^[a-zA-Z0-9]+$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}

	return string(result)
}

// GenerateRandomAlphaString generates a random alphabetic string
func GenerateRandomAlphaString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)

	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}

	return string(result)
}

// GenerateRandomNumericString generates a random numeric string
func GenerateRandomNumericString(length int) string {
	const charset = "0123456789"
	result := make([]byte, length)

	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}

	return string(result)
}

// PadLeft pads a string to the left with a character
func PadLeft(s string, length int, padChar rune) string {
	if len(s) >= length {
		return s
	}

	padding := strings.Repeat(string(padChar), length-len(s))
	return padding + s
}

// PadRight pads a string to the right with a character
func PadRight(s string, length int, padChar rune) string {
	if len(s) >= length {
		return s
	}

	padding := strings.Repeat(string(padChar), length-len(s))
	return s + padding
}

// PadCenter pads a string to the center with a character
func PadCenter(s string, length int, padChar rune) string {
	if len(s) >= length {
		return s
	}

	totalPadding := length - len(s)
	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding

	leftPad := strings.Repeat(string(padChar), leftPadding)
	rightPad := strings.Repeat(string(padChar), rightPadding)

	return leftPad + s + rightPad
}

// SplitAndTrim splits a string and trims each part
func SplitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// JoinNonEmpty joins non-empty strings with a separator
func JoinNonEmpty(sep string, strs ...string) string {
	var result []string

	for _, s := range strs {
		if IsNotEmpty(s) {
			result = append(result, s)
		}
	}

	return strings.Join(result, sep)
}

// RemoveDuplicates removes duplicate strings from a slice
func RemoveDuplicates(strs []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, s := range strs {
		if !keys[s] {
			keys[s] = true
			result = append(result, s)
		}
	}

	return result
}

// Intersection returns the intersection of two string slices
func Intersection(a, b []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, s := range a {
		keys[s] = true
	}

	for _, s := range b {
		if keys[s] {
			result = append(result, s)
		}
	}

	return result
}

// Union returns the union of two string slices
func Union(a, b []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, s := range a {
		if !keys[s] {
			keys[s] = true
			result = append(result, s)
		}
	}

	for _, s := range b {
		if !keys[s] {
			keys[s] = true
			result = append(result, s)
		}
	}

	return result
}

// Difference returns the difference of two string slices
func Difference(a, b []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, s := range b {
		keys[s] = true
	}

	for _, s := range a {
		if !keys[s] {
			result = append(result, s)
		}
	}

	return result
}
