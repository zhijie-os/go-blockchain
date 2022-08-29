package main

import "fmt"

const subsidy = 10

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	// create transaction input
	txin := TXInput{[]byte{}, -1, data}
	// create transaction output
	txout := TXOutput{subsidy, to} // subsidy is the amount of reward
	// in bitcoin the subsidy is not stored but calculated based only on
	// the total numbers of blocks.

	// create transaction
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	// set id
	tx.ID = []byte{1}

	return &tx
}
