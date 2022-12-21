package blockchain

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v3"
)

const dbPath = "./tmp/blocks"

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.Dir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")

			genesis := Genesis()
			fmt.Println("Genesis proved")

			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)

			err = txn.Set([]byte("lh"), genesis.Hash)
			Handle(err)

			lastHash = genesis.Hash

			return err
		}

		item, err := txn.Get([]byte("lh"))
		Handle(err)

		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})

		return err
	})
	Handle(err)

	blockchain := &BlockChain{LastHash: lastHash, Database: db}

	return blockchain
}

func (chain *BlockChain) AddBlock(data string) {
	fmt.Println("Added")
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)

		err = item.Value(func(val []byte) error {
			lastHash = val

			return nil
		})

		return err
	})
	Handle(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)

		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	Handle(err)
}

func (chain *BlockChain) Interator() *BlockChainIterator {
	inter := &BlockChainIterator{chain.LastHash, chain.Database}

	return inter
}

func (inter *BlockChainIterator) Next() *Block {
	var block *Block

	err := inter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(inter.CurrentHash)
		Handle(err)

		err = item.Value(func(val []byte) error {
			block = Deserializer(val)

			return nil
		})

		return err
	})
	Handle(err)

	inter.CurrentHash = block.PrevHash

	return block
}
