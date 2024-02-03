package crypto

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/stretchr/testify/assert"
	"log"
	"math/big"
	"testing"
)

func TestKeypairSignVerifySuccess(t *testing.T) {
	privateKey := GeneratePrivateKey()
	publicKey := privateKey.PublicKey()
	msg := []byte("hello world")

	sig, err := privateKey.Sign(msg)
	assert.Nil(t, err)

	assert.True(t, sig.Verify(publicKey, msg))
}

func TestKeypairSignVerifySuccess3(t *testing.T) {
	xVal := new(big.Int)
	xVal.SetString("77285c4c273f3139c875b875b83e99b37fb1b0d9f82026361ef2278972fffc35", 16)
	yVal := new(big.Int)
	yVal.SetString("afb7806f1da92df2d2c3626dcc74daaa108016a409e3f6adceb34cdc0dd522f1", 16)

	rVal := new(big.Int)
	rVal.SetString("25a4cc9d2176c973dbbbc87e3de512735f85c7bfc88913b8d076bab7abf2b721", 16)
	sVal := new(big.Int)
	sVal.SetString("1f67d8aa3bd0ac77ad02bada5609489a2161b6d6c82e606728e4aa19760746b6", 16)

	msgHash := fmt.Sprintf(
		"%x",
		sha256.Sum256([]byte("hello")),
	)

	message, hashDecodeError := hex.DecodeString(msgHash)

	if hashDecodeError != nil {
		log.Println(hashDecodeError)
		panic("internal server error")
	}

	ecdsaPublicKey := ecdsa.PublicKey{
		Curve: secp256k1.S256(),
		X:     xVal,
		Y:     yVal,
	}

	publicKey := PublicKey{
		Key: &ecdsaPublicKey,
	}

	signature := Signature{
		S: sVal,
		R: rVal,
	}
	assert.True(t, signature.Verify(publicKey, message))
}

func MakeByteToBigint(data []byte) *big.Int {
	result := new(big.Int)
	result.SetBytes(data)

	//fmt.Println("바이트슬라이스 -> 빅인트 :", result)
	return result
}

//func TestKeypairSignVerifySuccess2(t *testing.T) {
//	msgHash := fmt.Sprintf(
//		"%x",
//		sha256.Sum256([]byte("hello")),
//	)
//
//	msg, hashDecodeError := hex.DecodeString(msgHash)
//
//	if hashDecodeError != nil {
//		log.Println(hashDecodeError)
//		panic("internal server error")
//	}
//
//	sigBytes, _ := hex.DecodeString("970bde5760aaee9ed846e2df130377bde7232e359bf60c7a1a4d6e0bd8c4ecf38980169f06bbc787bac22b3c4dec3ab63853a2e3efb074d119efd96647ea5de7")
//	fmt.Println("sigBytes: ", sigBytes)
//	fmt.Println("sigBytesS: ", MakeByteToBigint(sigBytes[32:]))
//	fmt.Println("sigBytesR: ", MakeByteToBigint(sigBytes[:32]))
//	signature := &Signature{
//		S: MakeByteToBigint(sigBytes[32:]),
//		R: MakeByteToBigint(sigBytes[:32]),
//	}
//
//	xbytes, _ := hex.DecodeString("b286847d97818b6b7acc377fab09522d3b17954279a6bd236b0a16031e9df818")
//	ybytes, _ := hex.DecodeString("9867d1b6f7b922073c53908111ccea65cc803ce54b15d44ef33a6c40ff3d8c6d")
//	x := MakeByteToBigint(xbytes)
//	y := MakeByteToBigint(ybytes)
//
//	publicKey := elliptic.MarshalCompressed(elliptic.P256(), x, y)
//	fmt.Println("publicKey: ", publicKey)
//	fmt.Println("signature: ", signature)
//	assert.True(t, signature.Verify(publicKey, msg))
//}

//func TestKeypairSignVerifyFail(t *testing.T) {
//	privKey := GeneratePrivateKey()
//	PublicKey := privKey.PublicKey()
//	msg := []byte("hello world")
//
//	sig, err := privKey.Sign(msg)
//	assert.Nil(t, err)
//
//	otherPrivKey := GeneratePrivateKey()
//	otherPublicKey := otherPrivKey.PublicKey()
//
//	assert.False(t, sig.Verify(otherPublicKey, msg))
//	assert.False(t, sig.Verify(PublicKey, []byte("xxxxxx")))
//}
