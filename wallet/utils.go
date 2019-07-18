package wallet

import (
	"github.com/btcsuite/btcutil/base58"
)

func Base58Encode(input []byte) []byte {
	encoded := base58.Encode(input)

	return []byte(encoded)
}

func Base58Decode(input []byte) []byte {
	decoded := base58.Decode(string(input[:]))
	return decoded
}
