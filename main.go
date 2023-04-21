package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type Blockchain struct {
	blocks []*Block
}

type Block struct {
	Hash 	 []byte
	Data 	 []byte
	PrevHash []byte
}

func (b * Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func CreateBlock(data string,  PrevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), PrevHash}
	block.DeriveHash()
	return block
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}}
}

func (chain *Blockchain) AddBlock(data string){
	beforeBlock := chain.blocks[len(chain.blocks)-1]
	newBlock := CreateBlock(data, beforeBlock.Hash)

	chain.blocks = append(chain.blocks, newBlock)
}

func main() {
	chain := InitBlockchain()

	chain.AddBlock("first")
	chain.AddBlock("second")
	chain.AddBlock("third")

	for _, block := range chain.blocks {
		fmt.Printf("\n")
		fmt.Printf("Prev Hash : %v\n", block.PrevHash)
		fmt.Printf("Data : %v\n", block.Data)
		fmt.Printf("Hash : %v\n", block.Hash)
	}
}