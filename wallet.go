package main

import "encoding/hex"

type wallet struct {
	Address    string
	PublicKey  []byte
	PrivateKey []byte
}

func (w wallet) String() string {
	return "Private key:\t" + hex.EncodeToString(w.PrivateKey) + "\n" +
		"Public key:\t" + hex.EncodeToString(w.PublicKey) + "\n" +
		"Address:\t" + w.Address
}
