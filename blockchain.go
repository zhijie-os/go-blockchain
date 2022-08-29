package main

// blockchain stores blocks, the blocks should be linked together by the PreviousHash
import (
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain_%s.db"
const blocksBucket = "blocks"

type Blockchain struct {
	// tip is the last block, but it is the head of the linked list
	tip []byte   // tip of the chain
	db  *bolt.DB // the connect to the db file
}

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		// get the bucket
		b := tx.Bucket([]byte(blocksBucket))
		// get the last block hash
		lastHash = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	// create an NewBlock with last block as preivous block
	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		// put hash -> block metadata byte stream
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		// update the last block
		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

// first block in the chain is called "genesis" block
func NewGenesisBlock(coinbase *Transaction) *Block {
	// the data is "Genesis Block"
	// no previous hash
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

/* create a new blockchain:

   Open a DB file.
   Check if there’s a blockchain stored in it.
   If there’s a blockchain:
       Create a new Blockchain instance.
       Set the tip of the Blockchain instance to the last block hash stored in the DB.
   If there’s no existing blockchain:
       Create the genesis block.
       Store in the DB.
       Save the genesis block’s hash as the last block hash.
       Create a new Blockchain instance with its tip pointing at the genesis block.

*/
func NewBlockchain() *Blockchain {
	var tip []byte

	// opne boltDB
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		// check if there are blocks Bucket already
		b := tx.Bucket([]byte(blocksBucket))
		// if no bucket, create
		if b == nil {
			// create the genesis block
			genesis := NewGenesisBlock()
			// create bucket
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}
			// put hash -> serialized infor into the DB
			err = b.Put(genesis.Hash, genesis.Serialize())
			// put the 'l' <=> 'last block' -> genesis.Hash; that is saying, the last block on the chain is the genesis block
			err = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	bc := Blockchain{tip, db}
	return &bc
}

// we need an iterator to iterate through the chain
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// get iterator of a blockchain
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
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
