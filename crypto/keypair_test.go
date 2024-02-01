package crypto

import (
	"crypto/elliptic"
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

	//fmt.Println("바이트슬라이스 -> 빅인트 :", result)
	return result
}

func TestKeypairSignVerifySuccess2(t *testing.T) {
	msg := []byte("hello")

	sigBytes, _ := hex.DecodeString("79447872796349563076636b574d4348552f575759434c70356357536c64467143444c6b6175554748327a37524c575277417142344643542f4d446b3346636b58373069477747347756373739346455734e617034513d3d")
	fmt.Println("sigBytes: ", sigBytes)
	fmt.Println("sigBytesS: ", MakeByteToBigint(sigBytes[32:]))
	fmt.Println("sigBytesR: ", MakeByteToBigint(sigBytes[:32]))
	signature := &Signature{
		S: MakeByteToBigint(sigBytes[32:]),
		R: MakeByteToBigint(sigBytes[:32]),
	}

	xbytes, _ := hex.DecodeString("eff8b37fa642fca1114d42f8d7d85d0278dde4a8709ef950381e05b5a35a3df6")
	ybytes, _ := hex.DecodeString("42e831f18a317edd43f9ad02aeec95b3142a662b5a140976b07a0aba0257fea3")
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
