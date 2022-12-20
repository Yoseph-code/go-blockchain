package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
	// "github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	Blocks []*Block
	// LastHash []byte
	// Database *badger.DB
}

func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}

func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, newBlock)
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	Hanlde(err)

	return res.Bytes()
}

func Deserializer(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Hanlde(err)

	return &block
}

func Hanlde(err error) {
	if err != nil {
		log.Panic(err)
	}
}
