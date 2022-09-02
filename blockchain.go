package main

// blockchain stores blocks, the blocks should be linked together by the PreviousHash
import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain_%s.db"
const blocksBucket = "blocks"
const genesisCoinbase = "The times 03/Jan/2009 Chancellor on brink of second bailout for banks"

func dbExists() bool {
	// if assignment-statement; condition
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

type Blockchain struct {
	// tip is the last block, but it is the head of the linked list
	tip []byte   // tip of the chain
	db  *bolt.DB // the connect to the db file
}

func (bc *Blockchain) AddBlock(transactions []*Transaction) {
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
	newBlock := NewBlock(transactions, lastHash)

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

func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists")
		os.Exit(1)
	}

	var tip []byte
	// opne boltDB
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		// create coinbase transaction
		cbtx := NewCoinbaseTX(address, genesisCoinbase)
		genesis := NewGenesisBlock(cbtx)

		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}
		// put block hash -> block info into DB
		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}
		// put the last -> genesis block hash
		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}

		tip = genesis.Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}
	return &bc
}

/*
create a new blockchain:

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
func NewBlockchain(address string) *Blockchain {
	if dbExists() == false {
		fmt.Printf("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	// opne boltDB
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		// check if there are blocks Bucket already
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})

	bc := Blockchain{tip, db}
	return &bc
}

// get iterator of a blockchain
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	// unspent transactions
	var unspentTXs []Transaction

	// spent transaction outputs, ID -> value
	spentTXOs := make(map[string][]int)

	// blockchain iterator
	bci := bc.Iterator()

	for {
		// get next
		block := bci.Next()

		// loop through all the current block's transactions
		for _, tx := range block.Transactions {
			// get transaction ID
			txID := hex.EncodeToString(tx.GetHash())
			// break point
		Outputs:
			// loop through all the transaction outputs
			for outIdx, out := range tx.Vout {
				// check if the output is spent, i.e, used as input of another transaction
				if spentTXOs[txID] != nil { // if the transaction output is spent
					for _, spentOut := range spentTXOs[txID] {
						// loop to next transaction output
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				// otherwise, it is not spent, check if it can be unlocked
				// with given address
				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			// gather all inputs that could unlock outputs locked with the provided address
			// if no coinbase
			if tx.IsCoinbase() == false {
				// for every inputs
				for _, in := range tx.Vin {
					// if can unlock the output with given address
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						// inTxID is spent as output
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		// if this is genesis block
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs // list of transactions containing unspent outputs
}

func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(address)

	// for each unspent transaction
	for _, tx := range unspentTransactions {
		// for each transaction output
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs // return unspent transaction outputs
}

func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	// iterate over all unspent transactions and accumulate their values
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.GetHash())

		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				// get unspent outputs that can be unlocked by the
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					// when the accumulated fund is greater than the fund we needed, stop
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}
