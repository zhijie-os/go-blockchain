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


# Persistence and CLI

One can think blockchain as a distributed database. A "database" should be persistent and always on.


## Database
The <ins>*Bitcoin Core*</ins>(developerd by Satoshi Nakamoto) uses <ins>LevelDB</ins>. However, I am going to use *BoltDB*.

*Bitcoin Core* uses two "buckets" to store data:

1. `blocks` stores metadata describing all the blocks in a chian.
2. `chainstate` stores the state of a blockchain.

### Storage Format
In `blocks`, the `key-value` pairs like look:

```   
'b' + 32-byte block hash -> block index record. Each record stores:
    * The block header
    * The height.
    * The number of transactions.
    * To what extent this block is validated.
    * In which file, and where in that file, the block data is stored.
    * In which file, and where in that file, the undo data is stored.
```
  
```   
'f' + 4-byte file number -> file information record. Each record stores:
    * The number of blocks stored in the block file with that number.
    * The size of the block file with that number ($DATADIR/blocks/blkNNNNN.dat).
    * The size of the undo file with that number ($DATADIR/blocks/revNNNNN.dat).
    * The lowest and highest height of blocks stored in the block file with that number.
    * The lowest and highest timestamp of blocks stored in the block file with that number.
```
  
```   
'l' -> 4-byte file number: the last block file number used.
```

```
'R' -> 1-byte boolean ('1' if true): whether we're in the process of reindexing.
```

```
'F' + 1-byte flag name length + flag name string -> 1 byte boolean ('1' if true, '0' if false): various flags that can be on or off. Currently defined flags include:
    * 'txindex': Whether the transaction index is enabled.
```

```
't' + 32-byte transaction hash -> transaction index record. These are optional and only exist if 'txindex' is enabled (see above). Each record stores:
    * Which block file number the transaction is stored in.
    * Which offset into that file the block the transaction is part of is stored at.
    * The offset from the start of that block to the position where that transaction itself is stored.
```

This is very detailed *Bitcoin Core*'s implementation. However, this project serves as learning purpose, some simplifications can be done.
We are only going to store:

1. `32-byte block hash -> Block structure (serialized)`
2. `l -> the hash of the last block in a chain`


