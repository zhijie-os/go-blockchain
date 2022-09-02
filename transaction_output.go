package main

/*
	In Bitcoin, the smallest unit is satoshi. A satoshi is 0.00000001 BTC.
*/

// where "coins" are stored
type TXOutput struct {
	Value        int
	ScriptPubKey string // locking value with a puzzle
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}
