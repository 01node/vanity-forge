package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecp256k1Wallet_Bech32Only(t *testing.T) {
	wallet := secp256k1Wallet{}
	assert.True(t, wallet.bech32Only("acde"))
	assert.True(t, wallet.bech32Only("023456789"))
	assert.True(t, wallet.bech32Only("acdefghjklmnpqrstuvwxyzACDEFGHJKLMNPQRSTUVWXYZ"))
	assert.False(t, wallet.bech32Only("abcde!"))
	assert.False(t, wallet.bech32Only("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"))
	assert.False(t, wallet.bech32Only("abcde"))
}

func TestSecp256k1Wallet_CheckRequiredDigits(t *testing.T) {
	wallet := secp256k1Wallet{}
	assert.True(t, wallet.CheckRequiredDigits("12345", 3))
	assert.False(t, wallet.CheckRequiredDigits("12345", 6))
}

func TestSecp256k1Wallet_CheckRequiredLetters(t *testing.T) {
	wallet := secp256k1Wallet{}
	assert.True(t, wallet.CheckRequiredLetters("abcde", 2))
	assert.False(t, wallet.CheckRequiredLetters("abcde", 6))
}

func TestSecp256k1Wallet_ValidateInput(t *testing.T) {
	wallet := secp256k1Wallet{}
	errs := wallet.ValidateInput("acde", 2, 3)
	assert.Empty(t, errs)

	errs = wallet.ValidateInput("abcde!", 2, 3)
	assert.NotEmpty(t, errs)
	assert.Contains(t, errs[0], "bech32 incompatible characters")

	errs = wallet.ValidateInput("acdefgacdefghjklmnpqrstuvwxyzACDEFGHJKLMNPQRSTUVWXYZ", 2, 3)
	assert.NotEmpty(t, errs)
	assert.Contains(t, errs[0], "is too long")

	errs = wallet.ValidateInput("acde", -2, 3)
	assert.NotEmpty(t, errs)
	assert.Contains(t, errs[0], "Can't require negative amount of characters")

	errs = wallet.ValidateInput("acde", 10, 30)
	assert.NotEmpty(t, errs)
	assert.Contains(t, errs[0], "Can't require more than 38 characters")
}

func TestSecp256k1Wallet_GenerateWallet(t *testing.T) {
	wallet := secp256k1Wallet{}
	w := wallet.GenerateWallet()
	assert.NotEmpty(t, w.Address)
	assert.NotNil(t, w.PublicKey)
	assert.NotNil(t, w.PrivateKey)
}
