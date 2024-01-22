package main

import (
	"crypto/ecdsa"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

type ecsdaWallet struct {
	Chain chain
}

// On Ethereum and other networks compatible with the Ethereum Virtual Machine (EVM), public addresses all share the same format: they begin with 0x, and are followed by 40 alphanumeric characters (numerals and letters), adding up to 42 characters in total. They're also not case sensitive.

// This address is a number, even though it also includes alphabetical characters. This is because the hexadecimal (base 16) system used to generate the address doesn't just use numerals, like our ten-digit decimal system. Instead, the hexadecimal system uses the numerals 0-9 and the letters A-F. This means it has 16 characters at its disposal, hence the name base 16. In computer science and many programming languages, the 0x prefix is used at the start of all hex numbers, as they are known, to differentiate them from decimal values.

// bech16digits is a constant string representing the valid characters in a Bech16 encoding.
const bech16digits = "0123456789"

// bech16letters represents the valid characters for a Bech16 encoding.
const bech16letters = "abcdefABCDEF"

// bech16chars is a constant string that represents the characters used in the bech16 encoding scheme, which includes both digits and letters.
const bech16chars = bech16digits + bech16letters

// bech16Only checks if the given string contains only characters from the bech16 character set.
func (w ecsdaWallet) bech16Only(s string) bool {
	return w.countUnionChars(s, bech16chars) == len(s)
}

// countUnionChars counts the number of characters in the input string 's' that are present in the 'letterSet'.
// It returns the count of such characters.
func (w ecsdaWallet) countUnionChars(s string, letterSet string) int {
	count := 0
	for _, char := range s {
		if strings.Contains(letterSet, string(char)) {
			count++
		}
	}
	return count
}

// CheckRequiredDigits checks if the given candidate string has the required number of digits.
// It counts the number of union characters between the candidate string and the bech16digits string.
// If the count is less than the required number, it returns false; otherwise, it returns true.
func (w ecsdaWallet) CheckRequiredDigits(candidate string, required int) bool {
	if w.countUnionChars(candidate, bech16digits) < required {
		return false
	}

	return true
}

// CheckRequiredLetters checks if a candidate string contains the required number of union characters.
// It returns true if the candidate string meets the requirement, otherwise false.
func (w ecsdaWallet) CheckRequiredLetters(candidate string, required int) bool {
	if w.countUnionChars(candidate, bech16letters) < required {
		return false
	}

	return true
}

// ValidateInput validates the input string based on the specified criteria.
// It checks if the input string contains bech16 incompatible characters,
// if it exceeds the maximum length of 40 characters,
// and if the required number of letters and digits are non-negative and do not exceed 40.
// It returns a slice of error messages indicating the validation errors, if any.
func (w ecsdaWallet) ValidateInput(SearchString string, RequiredLetters int, RequiredDigits int) []string {
	var errs []string
	if !w.bech16Only(SearchString) {
		errs = append(errs, "ERROR: "+SearchString+" contains bech16 incompatible characters")
	}
	if len(SearchString) > 40 {
		errs = append(errs, "ERROR: "+SearchString+" is too long. Must be max 40 characters.")
	}
	if RequiredDigits < 0 || RequiredLetters < 0 {
		errs = append(errs, "ERROR: Can't require negative amount of characters.")
	}
	if RequiredDigits+RequiredLetters > 40 {
		errs = append(errs, "ERROR: Can't require more than 40 characters.")
	}
	return errs
}

// GenerateWallet generates a new wallet by generating a private key and deriving the corresponding public key and address.
// It returns a wallet struct containing the address, public key, and private key bytes.
func (w ecsdaWallet) GenerateWallet() wallet {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return wallet{address, publicKeyBytes, privateKeyBytes}
}
