package main

import (
	"log"

	"github.com/boltdb/bolt"
)

// we need an iterator to iterate through the chain
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		// get the byte stream of the current block
		encodedBlock := b.Get(i.currentHash)
		// deserialize the information to plaintext
		block = DeserializeBlock(encodedBlock)
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	// previous block hash becomes next iteration's block
	i.currentHash = block.PrevBlockHash
	// return current iteration's block
	return block
}
