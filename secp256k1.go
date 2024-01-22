package main

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// secp256k1Wallet represents a secp256k1 wallet.
type secp256k1Wallet struct {
	Chain chain
}

// bech32digits represents the digits allowed in the Bech32 alphabet.
const bech32digits = "023456789"

// bech32letters represents the letters allowed in the Bech32 alphabet.
const bech32letters = "acdefghjklmnpqrstuvwxyzACDEFGHJKLMNPQRSTUVWXYZ"

// bech32chars represents the alphanumeric characters allowed in the Bech32 alphabet.
const bech32chars = bech32digits + bech32letters

// bech32Only checks if a string contains only characters from the Bech32 alphabet.
func (w secp256k1Wallet) bech32Only(s string) bool {
	return w.countUnionChars(s, bech32chars) == len(s)
}

// countUnionChars counts the number of characters in a string that are present in a given letter set.
func (w secp256k1Wallet) countUnionChars(s string, letterSet string) int {
	count := 0
	for _, char := range s {
		if strings.Contains(letterSet, string(char)) {
			count++
		}
	}
	return count
}

// CheckRequiredDigits checks if a candidate string contains the required number of digits.
func (w secp256k1Wallet) CheckRequiredDigits(candidate string, required int) bool {
	if w.countUnionChars(candidate, bech32digits) < required {
		return false
	}
	return true
}

// CheckRequiredLetters checks if a candidate string contains the required number of letters.
func (w secp256k1Wallet) CheckRequiredLetters(candidate string, required int) bool {
	if w.countUnionChars(candidate, bech32letters) < required {
		return false
	}
	return true
}

// ValidateInput validates the search string, required letters, and required digits.
// It returns a list of errors encountered during validation.
func (w secp256k1Wallet) ValidateInput(SearchString string, RequiredLetters int, RequiredDigits int) []string {
	var errs []string
	if !w.bech32Only(SearchString) {
		errs = append(errs, "ERROR: "+SearchString+" contains bech32 incompatible characters.")
	}
	if len(SearchString) > 38 {
		errs = append(errs, "ERROR: "+SearchString+" is too long. Must be max 38 characters.")
	}
	if RequiredDigits < 0 || RequiredLetters < 0 {
		errs = append(errs, "ERROR: Can't require negative amount of characters.")
	}
	if RequiredDigits+RequiredLetters > 38 {
		errs = append(errs, "ERROR: Can't require more than 38 characters.")
	}

	return errs
}

// GenerateWallet generates a new secp256k1 wallet.
func (w secp256k1Wallet) GenerateWallet() wallet {
	var privkey secp256k1.PrivKey = secp256k1.GenPrivKey()
	var pubkey secp256k1.PubKey = privkey.PubKey().(secp256k1.PubKey)
	bech32Addr, err := bech32.ConvertAndEncode(w.Chain.Prefix, pubkey.Address())
	if err != nil {
		panic(err)
	}

	return wallet{bech32Addr, pubkey, privkey}
}
