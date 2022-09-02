package main

type TXInput struct {
	Txid      []byte
	Vout      int    // stores an index of an output in the transaction
	ScriptSig string // a script which provides data to be used in an output's ScriptPubKey
	// if the data is correct, the output can be unlocked
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}
