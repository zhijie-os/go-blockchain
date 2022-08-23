# Basic Prototype

`block.go` and `blockchain.go` contains the prototype of the blockchain. The `blockchain` use an array (`[]block`) as the fundamental data structure. Each `block` is linked to the previous one.

However, there is a significiant flaw: andding block to the chain is easy and cheap. In real-life application like *Bitcoin , Raven, and Ethereum*, adding new blocks should be a hard work.
