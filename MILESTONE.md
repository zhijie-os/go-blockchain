# Basic Prototype

`block.go` and `blockchain.go` contains the prototype of the blockchain. The `blockchain` use an array (`[]block`) as the fundamental data structure. Each `block` is linked to the previous one.

However, there is a significiant flaw: andding block to the chain is easy and cheap. In real-life application like *Bitcoin , Raven, and Ethereum*, adding new blocks should be a hard work.

# Proof-of-Work

Although *Ethereum* is going to merge and use *Proof of Stake* in September 2023 (this time, it is real!) and I would probably no longer getting profit from mining for a very long time, most of popular coins (e.g *Bitcoin, Raven, Flux, Ergo, ETC, and etc.*) are using *Proof of Stake*.

## Hashcash

Bitcoin uses Hashcash, a PoW algorithm:

1. Take some publicly known data.
2. Add a counter(starts at 0) to it.
3. Get a hash of the `data + counter` combination.
4. Check that the hash meets certain requirements. If it doesn't, repeat steps 3 and 4.

The `counter` is called `nonce` and treated as a solution to the PoW problem's answer.

