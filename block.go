/*
	Data: the actual data stored in the current block
	PrevBlockHash: the previous block hash
	Hash: the current block hash

	`longest chain wins`
*/
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"strconv"
	"time"
)

// every block must store at least one transaction and
// it’s not possible to mine blocks without transactions
type Block struct {
	Timestamp     int64
	Transaction   []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// Calculate and set the hash for a block
func (b *Block) SetHash() {
	// := is for initialization, equivalent to var ... = ...

	// FormatInt returns the string of the block's timestamp in base 10, then store it as byte array
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	// concate PreviousBlockHash, current Data, and Timestamp to create Hash
	headers := bytes.Join(
		// Data, PrevBlockHash, and Timestamp are byte arrays
		[][]byte{b.PrevBlockHash,
			b.Data, timestamp},
		[]byte{},
	)

	// calculate the current hash
	hash := sha256.Sum256(headers)

	// set Hash value
	b.Hash = hash[:]
}

// Block factory function
func NewBlock(transactions []*Transaction, data string, prevBlockHash []byte) *Block {
	// create new block
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	// create the work for mining
	pow := NewProofOfWork(block)
	// mine the block
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// turn a block and its information into byte stream
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	// create gob encoder
	encoder := gob.NewEncoder(&result)
	// encode the block
	err := encoder.Encode(b)

	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// turn an encoding of block into plaintext of block
func DeserializeBlock(raw []byte) *Block {
	var block Block
	// create gob decoder
	decoder := gob.NewDecoder(bytes.NewReader(raw))
	// decode the block
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}

	return &block
}
