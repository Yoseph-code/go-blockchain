package main

import (
	"fmt"
	"strconv"

	"github.com/Yoseph-code/go-blockchain/blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()

	chain.AddBlock("first block afet genesis")
	chain.AddBlock("second block afet genesis")
	chain.AddBlock("third block afet genesis")

	for _, block := range chain.Blocks {
		fmt.Printf("Previus hash: %x \n", block.PrevHash)
		fmt.Printf("Block data: %s \n", block.Data)
		fmt.Printf("Block hash: %x \n", block.Hash)
		fmt.Println()

		pow := blockchain.NewProf(block)
		fmt.Printf("Pow %s \n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
