package barreldb

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
	"strconv"
)

func (barrelDB *BarrelDatabase) InsertTxWithHash(hash common.Hash, tx *types.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(types.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable("hash-tx").Put(hash.ToSlice(), buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) InsertTxWithNumber(number uint32, tx *types.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(types.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable("number-tx").Put([]byte(strconv.Itoa(int(number))), buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) InsertLastTx(tx *types.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(types.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable("lastTx").Put([]byte{}, buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) InsertLastTxNumber(number uint32) error {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, number)
	if err := barrelDB.GetTable("lastTxNumber").Put([]byte{}, b); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) SelectTxByHash(hash common.Hash) (*types.Transaction, error) {
	data, err := barrelDB.GetTable("hash-tx").Get(hash.ToSlice())
	if err != nil {
		return nil, err
	}

	tx := new(types.Transaction)
	err = tx.Decode(types.NewGobTxDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (barrelDB *BarrelDatabase) SelectTxByNumber(number uint32) (*types.Transaction, error) {
	data, err := barrelDB.GetTable("number-tx").Get([]byte(strconv.Itoa(int(number))))
	if err != nil {
		return nil, err
	}

	tx := new(types.Transaction)
	err = tx.Decode(types.NewGobTxDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (barrelDB *BarrelDatabase) SelectLastTx() (*types.Transaction, error) {
	data, err := barrelDB.GetTable("lastTx").Get([]byte{})
	if err != nil {
		return nil, err
	}

	tx := new(types.Transaction)
	err = tx.Decode(types.NewGobTxDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (barrelDB *BarrelDatabase) SelectLastTxNumber() (*uint32, error) {
	data, err := barrelDB.GetTable("lastTxNumber").Get([]byte{})
	if err != nil {

		fmt.Println("jijiji: ", err)
		return nil, err
	}

	number, err := strconv.Atoi(hex.EncodeToString(data))
	if err != nil {
		return nil, err
	}

	num := uint32(number)
	return &num, nil
}
