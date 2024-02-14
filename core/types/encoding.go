package types

import (
	"encoding/gob"
	"io"
)

type Encoder[T any] interface {
	Encode(T) error
}

type Decoder[T any] interface {
	Decode(T) error
}

type GobTxEncoder struct {
	w io.Writer
}

func NewGobTxEncoder(w io.Writer) *GobTxEncoder {
	return &GobTxEncoder{
		w: w,
	}
}

func (e *GobTxEncoder) Encode(tx *Transaction) error {
	return gob.NewEncoder(e.w).Encode(tx)
}

type GobTxDecoder struct {
	r io.Reader
}

func NewGobTxDecoder(r io.Reader) *GobTxDecoder {
	return &GobTxDecoder{
		r: r,
	}
}

func (e *GobTxDecoder) Decode(tx *Transaction) error {
	return gob.NewDecoder(e.r).Decode(tx)
}

type GobHeaderEncoder struct {
	w io.Writer
}

func NewGobHeaderEncoder(w io.Writer) *GobHeaderEncoder {
	return &GobHeaderEncoder{
		w: w,
	}
}

func (enc *GobHeaderEncoder) Encode(h *Header) error {
	return gob.NewEncoder(enc.w).Encode(h)
}

type GobHeaderDecoder struct {
	r io.Reader
}

func NewGobHeaderDecoder(r io.Reader) *GobHeaderDecoder {
	return &GobHeaderDecoder{
		r: r,
	}
}

func (dec *GobHeaderDecoder) Decode(h *Header) error {
	return gob.NewDecoder(dec.r).Decode(h)
}

type GobBlockEncoder struct {
	w io.Writer
}

func NewGobBlockEncoder(w io.Writer) *GobBlockEncoder {
	return &GobBlockEncoder{
		w: w,
	}
}

func (enc *GobBlockEncoder) Encode(b *Block) error {
	return gob.NewEncoder(enc.w).Encode(b)
}

type GobBlockDecoder struct {
	r io.Reader
}

func NewGobBlockDecoder(r io.Reader) *GobBlockDecoder {
	return &GobBlockDecoder{
		r: r,
	}
}

func (dec *GobBlockDecoder) Decode(b *Block) error {
	return gob.NewDecoder(dec.r).Decode(b)
}

//////

type GobAccountEncoder struct {
	w io.Writer
}

func NewGobAccountEncoder(w io.Writer) *GobAccountEncoder {
	return &GobAccountEncoder{
		w: w,
	}
}

func (enc *GobAccountEncoder) Encode(a *Account) error {
	return gob.NewEncoder(enc.w).Encode(a)
}

type GobAccountDecoder struct {
	r io.Reader
}

func NewGobAccountDecoder(r io.Reader) *GobAccountDecoder {
	return &GobAccountDecoder{
		r: r,
	}
}

func (dec *GobAccountDecoder) Decode(a *Account) error {
	return gob.NewDecoder(dec.r).Decode(a)
}
