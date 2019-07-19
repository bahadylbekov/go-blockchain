package blockchain

import (
	"crypto/elliptic"
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

const MiningReward = 50

func (tx *Transaction) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(tx)
	HandleErr(err)

	return res.Bytes()
}

func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		fmt.Sprintf("Coinbase TX to: %s", to)
	}
	txInput := TxInput{[]byte{}, -1, data}
	txOutput := TxOutput{MiningReward, to}

	tx := Transaction{nil, []TxInput{txInput}, []TxOutput{txOutput}}
	tx.SetId()
	return &tx
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	newTx := *tx
	newTx.ID = []byte{}

	hash = sha256.Sum256(newTx.Serialize())

	return hash[:]
}

func (tx *Transaction) SetId() {
	var hash [32]byte
	var encoded bytes.Buffer
	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	HandleErr(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (tx *Transaction) IsCoinBase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (tx *Transaction) SignTransaction(privateKey ecdsa.PrivateKey, prevTx map[string]Transaction) {
	if tx.IsCoinBase() {
		return
	}

	for _, in := range tx.Inputs {
		if prevTx[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("ERROR: Previous transaction doesn't exist")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inId, in := range txCopy.Inputs {
		prevTx := prevTx[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].PubKey = prevTx.Outputs[in.Out].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inId].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.ID)
		HandleErr(err)
		signature := append(r.Bytes(), s.Bytes()...)
		txCopy.Inputs[inId].Signature = signature
	}
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, in := range tx.Inputs {
		inputs = append(inputs, TxInput{in.ID, in.Out, nil, nil})
	}

	for _, out := range tx.Outputs {
		outputs = append(outputs, TxOutput{out.Value, out.PubKeyHash})
	}

	return Transaction{tx.ID, inputs, outputs}
}

func NewTransaction(from, to string, amount int, chain *Blockchain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := chain.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("Error: Account doesn't have enough funds")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		HandleErr(err)

		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetId()

	return &tx
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinBase() {
		return true
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID = nil {
			log.Panic("Previous transaction doesn't exist")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inId, in := range tx.Inputs {
		prevTx := prevTx[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].PubKey = prevTx.Outputs[in.Out].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inId].PubKey = nil

		r, s := big.Int{}
		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:(sigLen/2)])
		s.SetBytes(in.Signature[:(sigLen/2)])

		p, q := big.Int{}
		pubKeyLen := len(in.PubKey)
		p.SetBytes(in.PubKey[:(pubKeyLen/2)])
		q.SetBytes(in.PubKey[:(pubKeyLen/2)])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y} 
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false  {
			return false
		}
	}

}

func (tx Transaction) ToString() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("-- Transaction %x\n:", tx.ID))
	for i, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("Input %d:", i))
		lines = append(lines, fmt.Sprintf("TX.ID %x:", input.ID))
		lines = append(lines, fmt.Sprintf("Output %d:", input.Out))
		lines = append(lines, fmt.Sprintf("Signature %x:", input.Signature))
		lines = append(lines, fmt.Sprintf("Public Key %x:", input.PubKey))
	}

	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("Output %d:", i))
		lines = append(lines, fmt.Sprintf("Value %d:", output.Value))
		lines = append(lines, fmt.Sprintf("Script %x:", output.PubKeyHash))
	}

}