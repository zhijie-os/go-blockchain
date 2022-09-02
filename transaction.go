package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
)

const subsidy = 10

type Transaction struct {
	// ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// there is no inputs, use -1 as indicator
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = "Coinbase"
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{[]TXInput{txin}, []TXOutput{txout}}

	return &tx
}

func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// return the hash of the transaction
func (tx Transaction) GetHash() []byte {
	var encoded bytes.Buffer
	var hash [32]byte

	// create encoder
	enc := gob.NewEncoder(&encoded)
	// encode the transaction into enc
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())

	return hash[:]
}

func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, _ := hex.DecodeString(txid)

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}
	// Build a list of outputs
	// create output locked by the receiver's address
	outputs = append(outputs, TXOutput{amount, to})

	if acc > amount {
		// change, locked by sender's address
		outputs = append(outputs, TXOutput{acc - amount, from})
	}

	tx := Transaction{inputs, outputs}

	return &tx
}
