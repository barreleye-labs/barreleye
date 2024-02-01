package crypto

import (
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeypairSignVerifySuccess(t *testing.T) {
	privKey := GeneratePrivateKey()
	fmt.Println("pri:", privKey)
	publicKey := privKey.PublicKey()
	fmt.Println("publicKey:", publicKey)
	msg := []byte("hello world")

	sig, err := privKey.Sign(msg)
	fmt.Println("sig:", sig)
	fmt.Println("sigstr:", sig.String())
	assert.Nil(t, err)

	assert.True(t, sig.Verify(publicKey, msg))
}

func MakeByteToBigint(data []byte) *big.Int {
	result := new(big.Int)
	result.SetBytes(data)

	//fmt.Println("바이트슬라이스 -> 빅인트 :", result)
	return result
}

func TestKeypairSignVerifySuccess2(t *testing.T) {
	msgHash := fmt.Sprintf(
		"%x",
		sha256.Sum256([]byte("hello")),
	)

	msg, hashDecodeError := hex.DecodeString(msgHash)

	if hashDecodeError != nil {
		log.Println(hashDecodeError)
		panic("internal server error")
	}

	sigBytes, _ := hex.DecodeString("970bde5760aaee9ed846e2df130377bde7232e359bf60c7a1a4d6e0bd8c4ecf38980169f06bbc787bac22b3c4dec3ab63853a2e3efb074d119efd96647ea5de7")
	fmt.Println("sigBytes: ", sigBytes)
	fmt.Println("sigBytesS: ", MakeByteToBigint(sigBytes[32:]))
	fmt.Println("sigBytesR: ", MakeByteToBigint(sigBytes[:32]))
	signature := &Signature{
		S: MakeByteToBigint(sigBytes[32:]),
		R: MakeByteToBigint(sigBytes[:32]),
	}

	xbytes, _ := hex.DecodeString("b286847d97818b6b7acc377fab09522d3b17954279a6bd236b0a16031e9df818")
	ybytes, _ := hex.DecodeString("9867d1b6f7b922073c53908111ccea65cc803ce54b15d44ef33a6c40ff3d8c6d")
	x := MakeByteToBigint(xbytes)
	y := MakeByteToBigint(ybytes)

	publicKey := elliptic.MarshalCompressed(elliptic.P256(), x, y)
	fmt.Println("publicKey: ", publicKey)
	fmt.Println("signature: ", signature)
	assert.True(t, signature.Verify(publicKey, msg))
}

func TestKeypairSignVerifyFail(t *testing.T) {
	privKey := GeneratePrivateKey()
	PublicKey := privKey.PublicKey()
	msg := []byte("hello world")

	sig, err := privKey.Sign(msg)
	assert.Nil(t, err)

	otherPrivKey := GeneratePrivateKey()
	otherPublicKey := otherPrivKey.PublicKey()

	assert.False(t, sig.Verify(otherPublicKey, msg))
	assert.False(t, sig.Verify(PublicKey, []byte("xxxxxx")))
}
