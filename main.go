package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	block "goChain/blockChain"
	blockchain "goChain/blockChain"
)

type Commandline struct {
	blockchain *blockchain.BlockChain
}

func (cli *Commandline) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block Block_data - add a block to the chain")
	fmt.Println(" print - Prints the blocks in the chain")
	fmt.Println(" printAll - Prints All the blocks in the chain")
}

func (cli *Commandline) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *Commandline) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Added Block!")
}

func (cli *Commandline) printChain() {
	iter := cli.blockchain.Iterator()
	block := iter.Next()

	fmt.Println("\n")
	fmt.Printf("prev Hash %x\n", block.PrevHash)
	fmt.Printf("Data In Block %s\n", block.Data)
	fmt.Printf("Hash  %x\n", block.Hash)

	pow := blockchain.NewProof(block)
	fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
	fmt.Println("\n")
}

func (cli *Commandline) printAllChain() {
	fmt.Println("Print All Blocks")
	iter := cli.blockchain.Iterator()
	block := iter.Next()

	for {

		fmt.Println("\n")
		fmt.Printf("prev Hash  ----> %x \n", block.PrevHash)
		fmt.Printf("Data In Block ----> %s \n", block.Data)
		fmt.Printf("Hash  ---->  %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("Pow:  ----> %s\n", strconv.FormatBool(pow.Validate()))

		block = iter.GetByPrevHash(block.PrevHash)
		if block == nil {
			return
		}
	}

}

func (cli *Commandline) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	printAllChainCmd := flag.NewFlagSet("printAll", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.ErrorHandle(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.ErrorHandle(err)
	case "printAll":
		err := printAllChainCmd.Parse(os.Args[2:])
		blockchain.ErrorHandle(err)
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		blockData := strings.Join(addBlockCmd.Args(), " ")
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData + blockData)
	}

	if printChainCmd.Parsed() {
		// print를 입력했을 떄
		cli.printChain()
	}

	if printAllChainCmd.Parsed() {
		// printAll를 입력했을 떄
		cli.printAllChain()
	}
}

func main() {
	chain := block.InitBlockChain()
	defer chain.Database.Close()

	cli := Commandline{chain}
	cli.run()
}
