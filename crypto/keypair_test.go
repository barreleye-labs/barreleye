package crypto

import (
	"encoding/hex"
	"fmt"
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

	fmt.Println("바이트슬라이스 -> 빅인트 :", result)
	return result
}

func TestKeypairSignVerifySuccess2(t *testing.T) {
	//msg := []byte("hello world")

	sigBytes, _ := hex.DecodeString("b1413bdc270769a4da265c198bd2ec2b48b05094359fa35fde299f9abad370b70e3c912d94df0e3ee4c855b16f4ee114bd754a67a4519d1e1eed9ad00601670f")
	fmt.Println("sigBytes: ", sigBytes)
	fmt.Println("sigBytesS: ", MakeByteToBigint(sigBytes[32:]))
	fmt.Println("sigBytesR: ", MakeByteToBigint(sigBytes[:32]))
	s := &Signature{
		S: MakeByteToBigint(sigBytes[32:]),
		R: MakeByteToBigint(sigBytes[:32]),
	}
	pubkey := []byte{195, 190, 123, 195, 135, 194, 187, 195, 180, 194, 144, 65,
		194, 170, 44, 195, 181, 53, 37, 194, 141, 11, 20, 194,
		184, 194, 162, 53, 195, 165, 194, 154, 195, 185, 73, 194,
		166, 195, 143, 194, 173, 194, 188, 45, 109, 195, 180, 195,
		191, 195, 163, 79, 43, 194, 145, 85, 115, 194, 171, 195,
		139, 195, 128, 194, 138, 195, 134, 195, 129, 70, 106, 94,
		69, 17, 23, 194, 143, 194, 150, 194, 168, 194, 131, 195,
		179, 42, 195, 177, 26, 195, 157, 99, 46, 195, 185, 91,
		195, 154, 194, 191, 34}
	pub := PublicKey(pubkey)
	fmt.Println("aaa: ", pub)
	fmt.Println(s)
	//sig, err := privKey.Sign(msg)
	//fmt.Println("sig:", sig)
	//fmt.Println("sigstr:", sig.String())
	//assert.Nil(t, err)
	//
	//assert.True(t, sig.Verify(publicKey, msg))
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
