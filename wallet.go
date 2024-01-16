package main

import (
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type wallet struct {
	Address string
	Pubkey  []byte
	Privkey []byte
}

func (w wallet) String() string {
	return "Private key:\t" + hex.EncodeToString(w.Privkey) + "\n" +
		"Public key:\t" + hex.EncodeToString(w.Pubkey) + "\n" +
		"Address:\t" + w.Address
}

func generateWallet(chain string) wallet {
	var prefix = internalPrefixFromChain(chain)
	var privkey secp256k1.PrivKey = secp256k1.GenPrivKey()
	var pubkey secp256k1.PubKey = privkey.PubKey().(secp256k1.PubKey)
	bech32Addr, err := bech32.ConvertAndEncode(prefix, pubkey.Address())
	if err != nil {
		panic(err)
	}

	return wallet{bech32Addr, pubkey, privkey}
}
