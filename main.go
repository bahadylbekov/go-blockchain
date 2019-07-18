package main

import (
	"os"

	"github.com/bahadylbekov/go-blockchain/cli"
)

func main() {
	defer os.Exit(0)

	cli := cli.CommandLine{}
	cli.Run()

	// w := wallet.CreateWallet()
	// w.Address()
}
