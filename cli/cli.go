package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/bahadylbekov/go-blockchain/blockchain"
	"github.com/bahadylbekov/go-blockchain/wallet"
)

type CommandLine struct {
}

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func (cli *CommandLine) printUsage() {
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("getbalance -address ADDRESS - Get the balance of address")
	fmt.Println("createblockchain -address ADDRESS - Create new blockchain and init account by address")
	fmt.Println("chaindata - Print all blockchain data")
	fmt.Println("transfer -from FROM -to TO -amount AMOUNT - Transfer money from one account to another account")
	fmt.Println("createwallet - Create new wallet addresss")
	fmt.Println("addresses - List of all addresses in the blockchain network")
	fmt.Println("reindexUTXO - Rebuild the UTXO set")
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
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) createBlockchain(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("Address is not valid")
	}

	chain := blockchain.InitBlockchain(address)
	defer chain.Database.Close()

	UTXOSet := blockchain.UTXOSet{chain}
	UTXOSet.Reindex()

	fmt.Printf("Blockchain created by %s\n", address)
}

func (cli *CommandLine) getBalance(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("Address is not valid")
	}

	chain := blockchain.ContinueBlockchain(address)
	UTXOSet := blockchain.UTXOSet{chain}
	defer chain.Database.Close()

	balance := 0
	pubKeyHash := wallet.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CommandLine) reindexUTXO() {
	chain := blockchain.ContinueBlockchain("")
	defer chain.Database.Close()
	UTXOSet := blockchain.UTXOSet{chain}
	UTXOSet.Reindex()

	count := UTXOSet.CountUTXO()
	fmt.Printf("Done! There are %d transactions in the UTXO set. \n", count)
}

func (cli *CommandLine) transfer(from, to string, amount int) {
	if !wallet.ValidateAddress(from) {
		log.Panic("Address is not valid")
	}

	if !wallet.ValidateAddress(to) {
		log.Panic("Address is not valid")
	}

	chain := blockchain.ContinueBlockchain(from)
	UTXOSet := blockchain.UTXOSet{chain}
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, &UTXOSet)
	block := chain.AddBlock([]*blockchain.Transaction{tx})
	UTXOSet.Update(block)
	fmt.Printf("Successfully transfered: %d from %s to %s\n", amount, from, to)
}

func (cli *CommandLine) listAddresses() {
	wallets, _ := wallet.CreateWallets()
	addresses := wallets.GetAllAddresses()
	fmt.Println()

	for _, address := range addresses {
		fmt.Println(address)
	}
	fmt.Println()
}

func (cli *CommandLine) createWallet() {
	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()

	wallets.SaveFile()
	fmt.Printf("New Wallet created: %s\n", address)
}

func (cli *CommandLine) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	transferCmd := flag.NewFlagSet("transfer", flag.ExitOnError)
	chainDataCmd := flag.NewFlagSet("chaindata", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	addressesCmd := flag.NewFlagSet("addresses", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	transferFrom := transferCmd.String("from", "", "Source wallet address")
	transferTo := transferCmd.String("to", "", "Destination wallet address")
	transferAmount := transferCmd.Int("amount", 0, "Amount to transfer")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		HandleErr(err)

	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		HandleErr(err)

	case "transfer":
		err := transferCmd.Parse(os.Args[2:])
		HandleErr(err)

	case "chaindata":
		err := chainDataCmd.Parse(os.Args[2:])
		HandleErr(err)

	case "addresses":
		err := addressesCmd.Parse(os.Args[2:])
		HandleErr(err)

	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		HandleErr(err)

	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		HandleErr(err)

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

	if addressesCmd.Parsed() {
		cli.listAddresses()
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO()
	}

}
