package main
// blockchain stores blocks, the blocks should be linked together by the PreviousHash

type Blockchain struct {
	blocks []*Block
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}


// first block in the chain is called "genesis" block
func NewGenesisBlock() *Block {
	// the data is "Genesis Block"
	// no previous hash
	return NewBlock("Genesis Block", []byte{})
}


func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{ NewGenesisBlock() }}
}
