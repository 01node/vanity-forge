package main

import (
	"regexp"
	"strings"
)

type matcher struct {
	Mode            string
	SearchString    string
	Chain           string
	RequiredLetters int
	RequiredDigits  int
}

// The Bech32 alphabet contains 32 characters, including lowercase letters a-z and the numbers 0-9, excluding the number 1 and the letters 'b', 'i', 'o' to avoid reader confusion.

const bech32digits = "023456789"
const bech32letters = "acdefghjklmnpqrstuvwxyzACDEFGHJKLMNPQRSTUVWXYZ"

// This is alphanumeric chars minus chars "1", "b", "i", "o" (case insensitive)
const bech32chars = bech32digits + bech32letters

func bech32Only(s string) bool {
	return countUnionChars(s, bech32chars) == len(s)
}

func countUnionChars(s string, letterSet string) int {
	count := 0
	for _, char := range s {
		if strings.Contains(letterSet, string(char)) {
			count++
		}
	}
	return count
}

func (m matcher) MatchDefault(candidate string) bool {
	switch m.Mode {
	case "contains":
		return strings.Contains(candidate, m.SearchString)
	case "starts-with":
		return strings.HasPrefix(candidate, m.SearchString)
	case "ends-with":
		return strings.HasSuffix(candidate, m.SearchString)
	case "regex":
		match, err := regexp.MatchString(m.SearchString, candidate)
		if err != nil {
			return false
		}
		return match
	default:
		return strings.Contains(candidate, m.SearchString)
	}
}

func (m matcher) Match(candidate string) bool {
	// Get chain prefix
	prefix := prefixFromChain(m.Chain)
	// Trim prefix from candidate
	candidate = strings.TrimPrefix(candidate, prefix)

	// Check if candidate contains required amount of digits
	if countUnionChars(candidate, bech32digits) < m.RequiredDigits {
		return false
	}

	// Check if candidate contains required amount of letters
	if countUnionChars(candidate, bech32letters) < m.RequiredLetters {
		return false
	}

	return m.MatchDefault(candidate)
}

func (m matcher) ValidationErrors() []string {
	var errs []string
	if !bech32Only(m.SearchString) {
		errs = append(errs, "ERROR: SearchString contains bech32 incompatible characters")
	}
	if len(m.SearchString) > 38 {
		errs = append(errs, "ERROR: SearchString is too long. Must be max 38 characters.")
	}
	if m.RequiredDigits < 0 || m.RequiredLetters < 0 {
		errs = append(errs, "ERROR: Can't require negative amount of characters")
	}
	if m.RequiredDigits+m.RequiredLetters > 38 {
		errs = append(errs, "ERROR: Can't require more than 38 characters")
	}
	return errs
}

func findMatchingWallets(ch chan wallet, quit chan struct{}, m matcher) {
	for {
		select {
		case <-quit:
			return
		default:
			w := generateWallet(m.Chain)
			if m.Match(w.Address) {
				// Do a non-blocking write instead of simple `ch <- w` to prevent
				// blocking when it's time to quit and ch is full.
				select {
				case ch <- w:
				default:
				}
			}
		}
	}
}

func findMatchingWalletConcurrent(m matcher, goroutines int) wallet {
	ch := make(chan wallet)
	quit := make(chan struct{})
	defer close(quit)

	for i := 0; i < goroutines; i++ {
		go findMatchingWallets(ch, quit, m)
	}
	return <-ch
}
