package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"io"
	"math/big"
)

type PrivateKey struct {
	key       *ecdsa.PrivateKey
	publicKey PublicKey
}

func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		R: r,
		S: s,
	}, nil
}

func NewPrivateKeyFromReader(r io.Reader) PrivateKey {
	key, err := ecdsa.GenerateKey(secp256k1.S256(), r) //ether 말고도 있음
	if err != nil {
		panic(err)
	}

	publicKey := PublicKey{
		&key.PublicKey,
	}

	return PrivateKey{
		key:       key,
		publicKey: publicKey,
	}
}

func GeneratePrivateKey() PrivateKey {

	return NewPrivateKeyFromReader(rand.Reader)
}

func (k PrivateKey) PublicKey() PublicKey {
	return k.publicKey
}

type PublicKey struct {
	Key *ecdsa.PublicKey
}

func (k PublicKey) String() string {
	return "tempPublicKeyString"
}

func (k PublicKey) Address() common.Address {
	h := sha256.Sum256([]byte("tempPublicKeyString"))

	return common.NewAddressFromBytes(h[len(h)-20:])
}

type Signature struct {
	S *big.Int
	R *big.Int
}

func (sig Signature) String() string {
	b := append(sig.S.Bytes(), sig.R.Bytes()...)
	return hex.EncodeToString(b)
}

func (sig Signature) Verify(publicKey PublicKey, data []byte) bool {
	return ecdsa.Verify(publicKey.Key, data, sig.R, sig.S)
}
