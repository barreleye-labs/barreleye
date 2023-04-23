package main

import (
	"fmt"

	b "block/block"
)

func main() {
	chain := b.InitBlockchain()

	chain.AddBlock("first")
	chain.AddBlock("second")
	chain.AddBlock("third")

	for _, block := range chain.Blocks {
		fmt.Printf("\n")
		fmt.Printf("Prev Hash : %v\n", block.PrevHash)
		fmt.Printf("Data : %v\n", block.Data)
		fmt.Printf("Hash : %v\n", block.Hash)
	}
}
