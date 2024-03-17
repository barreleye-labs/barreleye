package types

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/common/util"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"io"
	"math/big"
)

type PrivateKey struct {
	Key       *ecdsa.PrivateKey
	PublicKey PublicKey
}

func (k *PrivateKey) Decode(dec Decoder[*PrivateKey]) error {
	return dec.Decode(k)
}

func (k *PrivateKey) Encode(enc Encoder[*PrivateKey]) error {
	return enc.Encode(k)
}

func (k *PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.Key, data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		R: r,
		S: s,
	}, nil
}

func NewPrivateKeyFromReader(r io.Reader) *PrivateKey {
	key, err := ecdsa.GenerateKey(secp256k1.S256(), r) //ether 말고도 있음
	if err != nil {
		panic(err)
	}

	publicKey := PublicKey{
		&key.PublicKey,
	}

	return &PrivateKey{
		Key:       key,
		PublicKey: publicKey,
	}
}

func GeneratePrivateKey() *PrivateKey {

	return NewPrivateKeyFromReader(rand.Reader)
}

func CreatePrivateKey(hexKey string) (*PrivateKey, error) {
	key, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		return nil, err
	}

	publicKey := PublicKey{
		&key.PublicKey,
	}

	return &PrivateKey{
		Key:       key,
		PublicKey: publicKey,
	}, nil
}

type PublicKey struct {
	Key *ecdsa.PublicKey
}

func (k *PublicKey) String1() string {
	return "tempPublicKeyString"
}

func (k *PublicKey) Address() common.Address {
	h := sha256.Sum256(append(k.Key.X.Bytes(), k.Key.Y.Bytes()...))
	return common.NewAddressFromBytes(h[:20])
}

func GetPublicKey(xHex string, yHex string) (*PublicKey, error) {
	x := new(big.Int)
	x, ok := x.SetString(util.Rm0x(xHex), 16)
	if !ok {
		return nil, fmt.Errorf("invalid signerX")
	}

	y := new(big.Int)
	y, ok = y.SetString(util.Rm0x(yHex), 16)
	if !ok {
		return nil, fmt.Errorf("invalid signerY")
	}

	ecdsaPublicKey := ecdsa.PublicKey{
		Curve: secp256k1.S256(),
		X:     x,
		Y:     y,
	}

	return &PublicKey{
		Key: &ecdsaPublicKey,
	}, nil
}

type Signature struct {
	S *big.Int
	R *big.Int
}

func (sig *Signature) String() string {
	b := append(sig.S.Bytes(), sig.R.Bytes()...)
	return hex.EncodeToString(b)
}

func (sig *Signature) Verify(publicKey PublicKey, data []byte) bool {
	return ecdsa.Verify(publicKey.Key, data, sig.R, sig.S)
}

func GetSignature(rHex string, sHex string) (*Signature, error) {
	r := new(big.Int)
	r, ok := r.SetString(util.Rm0x(rHex), 16)
	if !ok {
		return nil, fmt.Errorf("invalid signatureR")
	}

	s := new(big.Int)
	s, ok = s.SetString(sHex, 16)
	if !ok {
		return nil, fmt.Errorf("invalid signatureS")
	}

	return &Signature{
		R: r,
		S: s,
	}, nil
}
