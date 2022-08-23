package main

const targetBits = 24 // difficulty
const maxNonce = math.MaxInt64

type ProofOfWork struct {
	block	*Block
	target	*big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	// left shift by 256-24 = 232 bits <=> 29 Bytes
	target.Lsh(target, uint(256-targetBits))
	// looks like 0x000100...00
	// if the proof is smaller than the target, then it is valid <=> all beginning bits are zeros.
	pow := &ProofOfWork{b, target}

	return pow
}


// think PoW struct as class, (pow *ProofOfWork) indicates it is a `method` of PoW `class`
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	// concatenate prevBlockHash,  Data, Timestamp, targetBits, and nounce
	data := bytes.join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// `mine` the block
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int	// the integer representation of the hash
	var hash [32]byte // store the sha256 result
	nounce := 0		// counter starts at 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	
	// brute force solving the problem
	for nonce < maxNounce {
		// prepare the byte array
		data := pow.prepareData(nonce)
		// hashing
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])
		
		// compare result
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}


