package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/bahadylbekov/go-blockchain/blockchain"
)

type CommandLine struct {
}

func (cli *CommandLine) printUsage() {
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("getbalance -address ADDRESS - get the balance of address")
	fmt.Println("createblockchain -address ADDRESS - create new blockchain and init account by address")
	fmt.Println("chaindata - print all blockchain data")
	fmt.Println("transfer -from FROM -to TO -amount AMOUNT - transfer money from one account to another account")
	fmt.Println()
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) printBlockchain() {
	chain := blockchain.ContinueBlockchain("")
	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Prev. Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) createBlockchain(address string) {
	chain := blockchain.InitBlockchain(address)
	defer chain.Database.Close()

	fmt.Printf("Blockchain created by %s\n", address)
}

func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockchain(address)
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CommandLine) transfer(from, to string, amount int) {
	chain := blockchain.ContinueBlockchain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})

	fmt.Printf("Successfully transfered: %d from %s to %s\n", amount, from, to)
}

func (cli *CommandLine) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	transferCmd := flag.NewFlagSet("transfer", flag.ExitOnError)
	chainDataCmd := flag.NewFlagSet("chaindata", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	transferFrom := transferCmd.String("from", "", "Source wallet address")
	transferTo := transferCmd.String("to", "", "Destination wallet address")
	transferAmount := transferCmd.Int("amount", 0, "Amount to transfer")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		blockchain.HandleErr(err)

	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		blockchain.HandleErr(err)

	case "transfer":
		err := transferCmd.Parse(os.Args[2:])
		blockchain.HandleErr(err)

	case "chaindata":
		err := chainDataCmd.Parse(os.Args[2:])
		blockchain.HandleErr(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if transferCmd.Parsed() {
		if *transferFrom == "" || *transferTo == "" || *transferAmount == 0 {
			transferCmd.Usage()
			runtime.Goexit()
		}
		cli.transfer(*transferFrom, *transferTo, *transferAmount)
	}

	if chainDataCmd.Parsed() {
		cli.printBlockchain()
	}
}
