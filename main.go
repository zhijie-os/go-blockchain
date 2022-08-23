package main

import (
	"fmt"
	"strconv"
)

func main() {
	// create blocks
	bc := NewBlockchain()
	bc.AddBlock("Send 1 BTC to Zhijie")
	bc.AddBlock("Send 2 more BTC to Zhijie")

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}

	// validate the blockchain
	for _, block := range bc.blocks {
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}

}
