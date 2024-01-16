package main

func prefixFromChain(chain string) string {
	switch chain {
	case "celestia":
		return "celestia1"
	default:
		return chain + "1"
	}
}

func internalPrefixFromChain(chain string) string {
	switch chain {
	case "celestia":
		return "celestia"
	default:
		return chain
	}
}
