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
	xVal.SetString("aef1d4e3fefd16fe21462e7c50bb2ee93256d7a0ca349e20b03ced22c342f122", 16)
	yVal := new(big.Int)
	yVal.SetString("561b43dcf616dc92d9e5415b3ccab15edae5f1ea40da920f94a49467d82b7104", 16)

	rVal := new(big.Int)
	rVal.SetString("22bf1ed4aab5e9a3dec5d4906770b9e063767d59928fdb6474cf234f400a30a8", 16)
	sVal := new(big.Int)
	sVal.SetString("6b35ece1e5d5378a9debff02232009b9487ecbb7aadd5d62a98f4aa6ad7c2f3", 16)

	msgHash := fmt.Sprintf(
		"%x",
		sha256.Sum256([]byte("hello")),
	)
	fmt.Println("aaaaa: ", []byte("033ad840087368c8cb21dde49cb8210f1898487a3d060b26e884190789537aab"))

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

func TestKeypairSignVerifySuccess4(t *testing.T) {
	xVal := new(big.Int)
	xVal.SetString("ee5f14a3558a740f7da4f87569a94d31e6d889a7f23ab3c497275eddde93364d", 16)
	yVal := new(big.Int)
	yVal.SetString("7dc7287e1f08343a9561340378082e9fb5f07652cc408065943cf5018aaf4998", 16)

	rVal := new(big.Int)
	rVal.SetString("f86e6837e7177c202b00faf35e731e4c8d65baa75c934bc89bad107d1f959f5f", 16)
	sVal := new(big.Int)
	sVal.SetString("14157ddeed9998c4cacbdfeb6938f681b9df97f9ae639def5869e43060b94815", 16)

	msgHash := fmt.Sprintf(
		"%x",
		sha256.Sum256([]byte("ab54d1b47432f5d8bfe6f747611470476225b597031e6be996ac95dc6ccbb9a119fbbc0f3eb2a449fcababee5f14a3558a740f7da4f87569a94d31e6d889a7f23ab3c497275eddde93364d7dc7287e1f08343a9561340378082e9fb5f07652cc408065943cf5018aaf4998")),
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
