package consensus

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/bahadylbekov/go-blockchain/blockchain"
)

const targetBits = 24

type ProofOfWork struct {
	Block  *blockchain.Block
	Target *big.Int
}

func NewProofOfWork(b *md.Block) *md.ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &md.ProofOfWork{b, target}

	return pow
}

func (pow *md.ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevBlockHash,
			pow.Block.Data,
			IntToHex(pow.Block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

func (pow *md.ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining block containing \"%s\"n", pow.block.Data)

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == 1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Printf("\n\n")

	return nonce, hash[:]
}

func (pow *md.ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
