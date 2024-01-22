package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcsdaWallet_GenerateWallet(t *testing.T) {
	wallet := ecsdaWallet{}
	w := wallet.GenerateWallet()
	assert.NotEmpty(t, w.Address)
	assert.NotNil(t, w.PublicKey)
	assert.NotNil(t, w.PrivateKey)
}
func TestEcsdaWallet_CountUnionChars(t *testing.T) {
	wallet := ecsdaWallet{}

	// Test case 1: Counting union characters in a string with valid letter set
	s1 := "abcdef123456"
	letterSet1 := "abc123"
	expectedCount1 := 6
	assert.Equal(t, expectedCount1, wallet.countUnionChars(s1, letterSet1), "Expected count of 6 for valid letter set")

	// Test case 2: Counting union characters in a string with empty letter set
	s2 := "abcdef123456"
	letterSet2 := ""
	expectedCount2 := 0
	assert.Equal(t, expectedCount2, wallet.countUnionChars(s2, letterSet2), "Expected count of 0 for empty letter set")

	// Test case 3: Counting union characters in an empty string with valid letter set
	s3 := ""
	letterSet3 := "abc123"
	expectedCount3 := 0
	assert.Equal(t, expectedCount3, wallet.countUnionChars(s3, letterSet3), "Expected count of 0 for empty string")

	// Test case 4: Counting union characters in a string with invalid letter set
	s4 := "abcdef123456"
	letterSet4 := "!@#$%"
	expectedCount4 := 0
	assert.Equal(t, expectedCount4, wallet.countUnionChars(s4, letterSet4), "Expected count of 0 for invalid letter set")
}
func TestEcsdaWallet_Bech16Only(t *testing.T) {
	wallet := ecsdaWallet{}

	// Test case 1: All characters are valid bech16 characters
	s1 := "abcdef123456"
	assert.True(t, wallet.bech16Only(s1), "Expected true for valid bech16 characters")

	// Test case 2: Some characters are not valid bech16 characters
	s2 := "abcde!@#$%"
	assert.False(t, wallet.bech16Only(s2), "Expected false for invalid bech16 characters")

	// Test case 3: Empty string
	s3 := ""
	assert.True(t, wallet.bech16Only(s3), "Expected true for empty string")

	// Test case 4: String with valid bech16 characters and additional characters
	s4 := "abcdef123456!"
	assert.False(t, wallet.bech16Only(s4), "Expected false for string with additional characters")
}
func TestEcsdaWallet_CheckRequiredDigits(t *testing.T) {
	wallet := ecsdaWallet{}

	// Test case 1: Candidate string has less required digits
	candidate1 := "abcdef123456"
	required1 := 7
	assert.False(t, wallet.CheckRequiredDigits(candidate1, required1), "Expected false for candidate string with less required digits")

	// Test case 2: Candidate string has exact required digits
	candidate2 := "abcdef123456"
	required2 := 6
	assert.True(t, wallet.CheckRequiredDigits(candidate2, required2), "Expected true for candidate string with exact required digits")

	// Test case 3: Candidate string has more than required digits
	candidate3 := "abcdef123456"
	required3 := 5
	assert.True(t, wallet.CheckRequiredDigits(candidate3, required3), "Expected true for candidate string with more than required digits")

	// Test case 4: Empty candidate string
	candidate4 := ""
	required4 := 3
	assert.False(t, wallet.CheckRequiredDigits(candidate4, required4), "Expected false for empty candidate string")
}
func TestEcsdaWallet_CheckRequiredLetters(t *testing.T) {
	wallet := ecsdaWallet{}

	// Test case 1: Candidate string has enough required letters
	candidate1 := "abcdef123"
	required1 := 5
	assert.True(t, wallet.CheckRequiredLetters(candidate1, required1), "Expected true for candidate string with enough required letters")

	// Test case 2: Candidate string does not have enough required letters
	candidate2 := "abc123"
	required2 := 10
	assert.False(t, wallet.CheckRequiredLetters(candidate2, required2), "Expected false for candidate string without enough required letters")

	// Test case 3: Candidate string is empty
	candidate3 := ""
	required3 := 5
	assert.False(t, wallet.CheckRequiredLetters(candidate3, required3), "Expected false for empty candidate string")

	// Test case 4: Required letters is 0
	candidate4 := "abc123"
	required4 := 0
	assert.True(t, wallet.CheckRequiredLetters(candidate4, required4), "Expected true for required letters equal to 0")
}

func TestEcsdaWallet_ValidateInput(t *testing.T) {
	wallet := ecsdaWallet{}

	// Test case 1: Valid input
	searchString1 := "abcdef123456"
	requiredLetters1 := 6
	requiredDigits1 := 6
	expectedErrs1 := []string(nil)
	assert.Equal(t, expectedErrs1, wallet.ValidateInput(searchString1, requiredLetters1, requiredDigits1), "Expected no errors for valid input")

	// Test case 2: Invalid bech16 characters
	searchString2 := "abcde!@#$%"
	requiredLetters2 := 6
	requiredDigits2 := 6
	expectedErrs2 := []string{"ERROR: abcde!@#$% contains bech16 incompatible characters"}
	assert.Equal(t, expectedErrs2, wallet.ValidateInput(searchString2, requiredLetters2, requiredDigits2), "Expected error for invalid bech16 characters")

	// Test case 3: String too long
	searchString3 := "abcdef123456abcdef123456abcdef123456abcdef1234567"
	requiredLetters3 := 6
	requiredDigits3 := 6
	expectedErrs3 := []string{"ERROR: abcdef123456abcdef123456abcdef123456abcdef1234567 is too long. Must be max 40 characters."}
	assert.Equal(t, expectedErrs3, wallet.ValidateInput(searchString3, requiredLetters3, requiredDigits3), "Expected error for string too long")

	// Test case 4: Negative required characters
	searchString4 := "abcdef123456"
	requiredLetters4 := -1
	requiredDigits4 := 6
	expectedErrs4 := []string{"ERROR: Can't require negative amount of characters."}
	assert.Equal(t, expectedErrs4, wallet.ValidateInput(searchString4, requiredLetters4, requiredDigits4), "Expected error for negative required characters")

	// Test case 5: Required characters exceed string length
	searchString5 := "abcdef123456abcdef123456abcdef123456abcdef123456abcdef123456abcdef123456"
	requiredLetters5 := 10
	requiredDigits5 := 10
	expectedErrs5 := []string{"ERROR: abcdef123456abcdef123456abcdef123456abcdef123456abcdef123456abcdef123456 is too long. Must be max 40 characters."}
	assert.Equal(t, expectedErrs5, wallet.ValidateInput(searchString5, requiredLetters5, requiredDigits5), "Expected error for required characters exceeding string length")
}
