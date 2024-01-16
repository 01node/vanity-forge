package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/pflag"
)

var (
	AvailableChains = []string{"cosmos", "celestia", "dydx"}
	MatcherModes    = []string{"contains", "starts-with", "ends-with", "regex"}
)

type settings struct {
	Encryption      string // secp256k1 or ethsecp256k1
	SelectedChain   string // chain selector string or nil
	MatcherMode     string // starts-with, ends-with, contains
	SearchString    string // search string
	NumAccounts     string // number of accounts to generate
	RequiredLetters int    // number of letters to generate
	RequiredDigits  int    // number of digits to generate
}

func main() {
	// Defined flags
	var accountsNumber = pflag.IntP("accounts-number", "n", 0, "Amount of accounts you need")
	var matcherMode = pflag.StringP("mode", "m", "", "Matcher mode (contains, starts-with, ends-with, regex)")
	var searchString = pflag.StringP("search", "s", "", "Search string")
	var chain = pflag.StringP("chain", "c", "", "Chain selector string")
	var letters = pflag.IntP("letters", "l", 0, "Amount of letters (a-z) that the address must contain")
	var digits = pflag.IntP("digits", "d", 0, "Amount of digits (0-9) that the address must contain")
	var verbose = pflag.BoolP("verbose", "v", false, "Verbose output")

	// Parse flags
	pflag.Parse()

	// Validate matcher mode flag if exists
	if *matcherMode != "" {
		if !slices.Contains(MatcherModes, *matcherMode) {
			fmt.Println("ERROR: Invalid matcher mode. Must be one of: contains, starts-with, ends-with, regex")
			os.Exit(1)
		}
	}

	// Validate chain flag
	if *chain != "" {
		if !slices.Contains(AvailableChains, *chain) {
			fmt.Println("ERROR: Invalid chain. Must be one of the available chains. (cross-vanity --help for more info)")
			os.Exit(1)
		}
	}

	// Validate letters flag
	if *letters < 0 || *letters > 38 {
		fmt.Println("ERROR: Invalid letters. Must be between 0 and 38.")
		os.Exit(1)
	}

	// Validate digits flag
	if *digits < 0 || *digits > 38 {
		fmt.Println("ERROR: Invalid digits. Must be between 0 and 38.")
		os.Exit(1)
	}

	// Letters + Digits must be less than 38
	if *letters+*digits > 38 {
		fmt.Println("ERROR: Letters + Digits must be less than 38.")
		os.Exit(1)
	}

	// Create settings struct with details from flags
	settings := settings{
		Encryption:      "secp256k1",
		SelectedChain:   *chain,
		MatcherMode:     *matcherMode,
		SearchString:    *searchString,
		NumAccounts:     strconv.Itoa(*accountsNumber),
		RequiredLetters: *letters,
		RequiredDigits:  *digits,
	}

	// Prompt user for missing settings on encryption
	// TODO add ethsecp256k1
	if settings.Encryption == "" {
		huh.NewSelect[string]().
			Title("Select encryption").
			Options(
				huh.NewOption("secp256k1", "secp256k1"),
			).
			Value(&settings.Encryption).
			Run()
	}

	// Prompt user for missing settings on chain
	selectChainOptions := make([]huh.Option[string], len(AvailableChains))
	for i, chain := range AvailableChains {
		selectChainOptions[i] = huh.NewOption(chain, chain)
	}

	if settings.SelectedChain == "" {
		huh.NewSelect[string]().
			Title("Select Chain").
			Options(selectChainOptions...).
			Value(&settings.SelectedChain).
			Run()
	}

	// Prompt user for missing settings on matcher mode
	selectMatcherModeOptions := make([]huh.Option[string], len(MatcherModes))
	for i, mode := range MatcherModes {
		selectMatcherModeOptions[i] = huh.NewOption(mode, mode)
	}

	if settings.MatcherMode == "" {
		huh.NewSelect[string]().
			Title("Matcher Mode").
			Options(selectMatcherModeOptions...).
			Value(&settings.MatcherMode).
			Run()
	}

	// Prompt user for missing settings on search string
	if settings.SearchString == "" {
		huh.NewInput().
			Title("Search string").
			CharLimit(38).
			Value(&settings.SearchString).
			Run()
	}

	// Prompt user for missing settings on number of accounts to generate
	if settings.NumAccounts == "0" {
		huh.NewInput().
			Title("Number of accounts to generate").
			Validate(func(s string) error {
				_, err := strconv.Atoi(s)
				return err
			}).
			Value(&settings.NumAccounts).
			Run()
	}

	// Initialize Matcher struct
	m := matcher{
		Mode:            settings.MatcherMode,
		SearchString:    strings.ToLower(*&settings.SearchString),
		Chain:           *&settings.SelectedChain,
		RequiredLetters: *&settings.RequiredLetters,
		RequiredDigits:  *&settings.RequiredDigits,
	}

	matcherValidationErrs := m.ValidationErrors()
	if len(matcherValidationErrs) > 0 {
		for i := 0; i < len(matcherValidationErrs); i++ {
			fmt.Println(matcherValidationErrs[i])
		}
		os.Exit(1)
	}

	var matchingWallet wallet
	NumAccountsInt, err := strconv.Atoi(*&settings.NumAccounts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	action := func() {
		for i := 0; i < NumAccountsInt; i++ {
			// TODO limit CPU cores by flag
			matchingWallet = findMatchingWalletConcurrent(m, runtime.NumCPU())

			fmt.Printf("\nFound a new matching wallet (%d out of %d):\n", i+1, NumAccountsInt)
			fmt.Println(matchingWallet)
		}
	}

	// Print out settings
	if *verbose == true {
		fmt.Println("Encryption: " + *&settings.Encryption)
		fmt.Println("Selected Chain: " + *&settings.SelectedChain)
		fmt.Println("Matcher Mode: " + *&settings.MatcherMode)
		fmt.Println("Search String: " + *&settings.SearchString)
		fmt.Println("Number of Accounts to Generate: " + *&settings.NumAccounts)
	}

	spinerr := spinner.New().
		Type(spinner.Meter).
		Action(action).
		Title(" Generating accounts...").
		Run()

	if spinerr != nil {
		fmt.Println(spinerr)
	}
}
