package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type BlockChain struct {
	blocks []*Block
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.Hash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

func (chain *BlockChain) addBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, newBlock)
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}

func main() {
	chain := InitBlockChain()

	chain.addBlock("first block afet genesis")
	chain.addBlock("second block afet genesis")
	chain.addBlock("third block afet genesis")

	for _, block := range chain.blocks {
		fmt.Printf("Previus hash: %x \n", block.PrevHash)
		fmt.Printf("Block data: %s \n", block.Data)
		fmt.Printf("Block hash: %x \n", block.Hash)
		fmt.Println()
	}
}
