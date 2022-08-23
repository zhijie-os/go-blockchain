/*
	Data: the actual data stored in the current block
	PrevBlockHash: the previous block hash
	Hash: the current block hash

	`longest chain wins`
*/

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
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
			b.data, timestamp},
		[]byte{},
	)

	// calculate the current hash
	hash := sha256.Sum256(headers)

	// set Hash value
	b.Hash = hash[:]
}

// Block factory function
func NewBlock(data string, prevBlockHash []byte) *Block {
	// create new block
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	// set block hash
	block.SetHash()
	return block
}
