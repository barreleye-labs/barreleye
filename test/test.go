package main

import (
	"fmt"
	"math/big"
)

var Difficulty = 18

func main() {
	target := big.NewInt(1) //1
	fmt.Println("target: ", target)
	target.Lsh(target, uint(238))
	fmt.Println("targetLsh: ", target)
}