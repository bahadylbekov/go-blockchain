package blockchain

import (
	"fmt"

	badger "github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

// Blockchain structure contains hash of last block and whole database
type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

// BlockchainIterator helps to get previous blocks
type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// AddBlock function adds one more block into blockchain
func (chain *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		HandleErr(err)
		lastHash, err = item.Value()

		return err
	})
	HandleErr(err)

	newBlock := NewBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		HandleErr(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	HandleErr(err)
}

func InitBlockchain() *Blockchain {
	var lastHash []byte

	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	HandleErr(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := GenesisBlock()
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			HandleErr(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			HandleErr(err)
			lastHash, err = item.Value()
			return err
		}
	})

	HandleErr(err)

	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

// Iterator - parser of blockchain
func (chain *Blockchain) Iterator() *BlockchainIterator {
	iter := &BlockchainIterator{chain.LastHash, chain.Database}

	return iter
}

// Next function - activate function for iterator
func (iter *BlockchainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		encodedBlock, err := item.Value()
		block = Deserialize(encodedBlock)

		return err
	})
	HandleErr(err)

	iter.CurrentHash = block.PrevBlockHash

	return block
}
