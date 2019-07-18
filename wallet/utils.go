package wallet

import (
	"log"

	"github.com/mr-tron/base58"
)

func Base58Encoding(input []byte) []byte {
	encode := base58.Encode(input)
	return []byte(encode)
}

func Base58Decoding(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	if err != nil {
		log.Panic(err)
	}
	return []byte(decode)
}
