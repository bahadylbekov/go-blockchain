package main

import (
	"fmt"

	"github.com/bahadylbekov/go-blockchain/blockchain"
)

func main() {
	blockchain := blockchain.NewBlockchain()

	blockchain.AddBlock("Ivan send 1 BTC to Greg")
	blockchain.AddBlock("Greg send 13 BTC to Chryssi")

	for _, block := range blockchain.Blocks {
		fmt.Printf("Previous Block Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Current Block Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
