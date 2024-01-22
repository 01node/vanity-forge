package main

import (
	"regexp"
	"strings"
)

// MatchWithMode matches the candidate string with the specified mode in the matcher.
// It returns true if the candidate matches the mode, otherwise false.
func (m matcher) MatchWithMode(candidate string) bool {
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

// Match checks if the candidate string matches the criteria specified in the matcher.
// It trims the prefix from the candidate, checks the required amount of digits and letters,
// and then calls MatchWithMode to perform the matching based on the mode.
// It returns true if the candidate matches the criteria, otherwise false.
func (m matcher) Match(candidate string) bool {
	candidate = strings.TrimPrefix(candidate, m.Chain.PrefixFull)

	if !m.CheckRequiredDigits(candidate, m.RequiredDigits) {
		return false
	}

	if !m.CheckRequiredLetters(candidate, m.RequiredLetters) {
		return false
	}

	return m.MatchWithMode(candidate)
}

// ValidateInput validates the input parameters of the matcher and returns any validation errors.
// It dynamically selects the appropriate generator based on the encryption type in the chain,
// and then calls the generator's ValidateInput method.
// It returns a slice of validation error messages.
func (m matcher) ValidateInput() []string {
	var generatorValidate func(SearchString string, RequiredLetters, RequiredDigits int) []string

	switch m.Chain.Encryption {
	case Secp256k1:
		var secp256k1generator = secp256k1Wallet{
			Chain: m.Chain,
		}

		generatorValidate = secp256k1generator.ValidateInput
	case ECSDA:
		var ecsdagenerator = ecsdaWallet{
			Chain: m.Chain,
		}

		generatorValidate = ecsdagenerator.ValidateInput
	default:
		var secp256k1generator = secp256k1Wallet{
			Chain: m.Chain,
		}
		generatorValidate = secp256k1generator.ValidateInput
	}

	return generatorValidate(m.SearchString, m.RequiredLetters, m.RequiredDigits)
}

// CheckRequiredDigits checks if the candidate string contains the required amount of digits.
// It dynamically selects the appropriate generator based on the encryption type in the chain,
// and then calls the generator's CheckRequiredDigits method.
// It returns true if the candidate contains the required amount of digits, otherwise false.
func (m matcher) CheckRequiredDigits(candidate string, required int) bool {
	var gcrd func(candidate string, required int) bool

	switch m.Chain.Encryption {
	case Secp256k1:
		var secp256k1generator = secp256k1Wallet{
			Chain: m.Chain,
		}

		gcrd = secp256k1generator.CheckRequiredDigits
	case ECSDA:
		var ecsdagenerator = ecsdaWallet{
			Chain: m.Chain,
		}

		gcrd = ecsdagenerator.CheckRequiredDigits
	default:
		var secp256k1generator = secp256k1Wallet{
			Chain: m.Chain,
		}
		gcrd = secp256k1generator.CheckRequiredDigits
	}

	return gcrd(candidate, required)
}

// CheckRequiredLetters checks if the candidate string contains the required amount of letters.
// It dynamically selects the appropriate generator based on the encryption type in the chain,
// and then calls the generator's CheckRequiredLetters method.
// It returns true if the candidate contains the required amount of letters, otherwise false.
func (m matcher) CheckRequiredLetters(candidate string, required int) bool {
	var gcrl func(candidate string, required int) bool

	switch m.Chain.Encryption {
	case Secp256k1:
		var secp256k1generator = secp256k1Wallet{
			Chain: m.Chain,
		}

		gcrl = secp256k1generator.CheckRequiredLetters
	case ECSDA:
		var ecsdagenerator = ecsdaWallet{
			Chain: m.Chain,
		}

		gcrl = ecsdagenerator.CheckRequiredLetters
	default:
		var secp256k1generator = secp256k1Wallet{
			Chain: m.Chain,
		}
		gcrl = secp256k1generator.CheckRequiredLetters
	}

	return gcrl(candidate, required)
}

// GenerateWallet generates a wallet based on the encryption type in the chain.
// It dynamically selects the appropriate generator based on the encryption type in the chain,
// and then calls the generator's GenerateWallet method.
// It returns the generated wallet.
func (m matcher) GenerateWallet() wallet {
	var generate func() wallet

	switch m.Chain.Encryption {
	case Secp256k1:
		var secp256k1generator = secp256k1Wallet{
			Chain: m.Chain,
		}

		generate = secp256k1generator.GenerateWallet
	case ECSDA:
		var ecsdagenerator = ecsdaWallet{
			Chain: m.Chain,
		}

		generate = ecsdagenerator.GenerateWallet
	default:
		var secp256k1generator = secp256k1Wallet{
			Chain: m.Chain,
		}
		generate = secp256k1generator.GenerateWallet
	}

	return generate()
}

// findMatchingWallets finds matching wallets based on the matcher criteria and sends them to the channel.
// It runs in a loop until the quit signal is received.
// It generates a wallet using the GenerateWallet method and checks if it matches the criteria using the Match method.
// If a match is found, it sends the wallet to the channel.
func findMatchingWallets(ch chan wallet, quit chan struct{}, m matcher) {
	for {
		select {
		case <-quit:
			return
		default:
			w := m.GenerateWallet()
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

// findMatchingWalletConcurrent finds a matching wallet concurrently using multiple goroutines.
// It creates a channel for sending and receiving wallets, and a quit channel for signaling the goroutines to stop.
// It spawns the specified number of goroutines, each running the findMatchingWallets function.
// It returns the first matching wallet received from the channel.
func findMatchingWalletConcurrent(m matcher, goroutines int) wallet {
	ch := make(chan wallet)
	quit := make(chan struct{})
	defer close(quit)

	for i := 0; i < goroutines; i++ {
		go findMatchingWallets(ch, quit, m)
	}
	return <-ch
}
