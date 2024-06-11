package main

import (
	"fmt"
	"strings"
)

// TinyURL struct to hold the character set and base.
type TinyURL struct {
	charset string
	base    int
}

// NewTinyURL initializes and returns a TinyURL instance.
func NewTinyURL() *TinyURL {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return &TinyURL{charset: charset, base: len(charset)}
}

// Encode encodes an integer to a base62 string.
func (t *TinyURL) Encode(num int) string {
	if num == 0 {
		return string(t.charset[0])
	}

	var encodedStr strings.Builder
	for num > 0 {
		rem := num % t.base
		encodedStr.WriteByte(t.charset[rem])
		num = num / t.base
	}

	// Reverse the encoded string since we built it backwards.
	runes := []rune(encodedStr.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func main() {
	tinyURL := NewTinyURL()

	// Encoding
	originalNumbers := []int{123, 124, 125, 234, 456, 789}
	for _, num := range originalNumbers {
		encodedStr := tinyURL.Encode(num)
		fmt.Printf("Encoded: %s\n", encodedStr)
	}
}
