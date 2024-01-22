package main

type Encryption int64

const (
	Undefined Encryption = iota
	Secp256k1
	Ethsecp256k1
	ECSDA
)

type chain struct {
	Name       string
	Prefix     string
	PrefixFull string
	Encryption
}

type settings struct {
	SelectedChain   chain  // chain selector string or nil
	MatcherMode     string // starts-with, ends-with, contains
	SearchString    string // search string
	NumAccounts     string // number of accounts to generate
	RequiredLetters int    // number of letters to generate
	RequiredDigits  int    // number of digits to generate
}

type walletgenerator struct {
	GenerateWallet func() wallet
}

type matcher struct {
	Mode            string
	SearchString    string
	Chain           chain
	RequiredLetters int
	RequiredDigits  int
}

var (
	AvailableChains = []chain{
		{
			Name:       "celestia",
			Prefix:     "celestia",
			PrefixFull: "celestia1",
			Encryption: Secp256k1,
		},
		{
			Name:       "cosmos",
			Prefix:     "cosmos",
			PrefixFull: "cosmos1",
			Encryption: Secp256k1,
		},
		{
			Name:       "dydx",
			Prefix:     "dydx",
			PrefixFull: "dydx1",
			Encryption: Secp256k1,
		},
		{
			Name:       "berachain",
			Prefix:     "0x",
			PrefixFull: "0x",
			Encryption: ECSDA,
		},
	}
	MatcherModes = []string{"contains", "starts-with", "ends-with", "regex"}
)
