package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	badger "github.com/dgraph-io/badger/v3"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from genesis"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func ContinueBlockchain(addrress string) *BlockChain {
	if DBexists() == false {
		fmt.Println("No existing blockchain, create one")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.Dir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)

		err = item.Value(func(val []byte) error {
			lastHash = val

			return nil
		})

		return err
	})
	Handle(err)

	chain := BlockChain{LastHash: lastHash, Database: db}

	return &chain
}

func InitBlockChain(address string) *BlockChain {
	var lastHash []byte

	if DBexists() {
		fmt.Println("Blockchain alredy exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.Dir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinBaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)

		err = txn.Set([]byte("lh"), genesis.Hash)

		lastHash = genesis.Hash

		return err
	})
	Handle(err)

	blockchain := &BlockChain{LastHash: lastHash, Database: db}

	return blockchain
}

func (chain *BlockChain) AddBlock(transactions []*Transaction) {
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

	newBlock := CreateBlock(transactions, lastHash)

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

func (chain *BlockChain) FindUnspentTransaction(address string) []Transaction {
	var unspentTxs []Transaction

	spendTxOs := make(map[string][]int)

	iter := chain.Interator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spendTxOs[txID] != nil {
					for _, spentOut := range spendTxOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					for in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)

						spendTxOs[inTxID] = append(spendTxOs[inTxID], in.Out)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTxs
}

func (chain *BlockChain) FindUTXO(address string) []TxOutput {
	var UTOXs []TxOutput
	unSpendTransactions := chain.FindUnspentTransaction(address)

	for _, tx := range unSpendTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTOXs = append(UTOXs, out)
			}
		}
	}

	return UTOXs
}

func (chain *BlockChain) FindSpendableOutPuts(address string, amount int) (int, map[string][]int) {
	unspendOuts := make(map[string][]int)
	unspendTxs := chain.FindUnspentTransaction(address)
	accumulated := 0

Work:
	for _, tx := range unspendTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspendOuts[txID] = append(unspendOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspendOuts
}
