package barreldb

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
	"strconv"
)

func (barrelDB *BarrelDatabase) InsertTxWithHash(hash common.Hash, tx *types.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(types.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable(HashTxTableName).Put(hash.ToSlice(), buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) InsertTxWithNumber(number uint32, tx *types.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(types.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable(NumberTxTableName).Put([]byte(strconv.Itoa(int(number))), buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) InsertLastTx(tx *types.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(types.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable(LastTxTableName).Put([]byte{}, buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) InsertLastTxNumber(number uint32) error {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, number)
	if err := barrelDB.GetTable(LastTxNumberTableName).Put([]byte{}, b); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) SelectTxByHash(hash common.Hash) (*types.Transaction, error) {
	data, err := barrelDB.GetTable(HashTxTableName).Get(hash.ToSlice())
	if err != nil {
		if err.Error() != common.LevelDBNotFoundError {
			return nil, err
		}
		return nil, nil
	}

	tx := new(types.Transaction)
	err = tx.Decode(types.NewGobTxDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (barrelDB *BarrelDatabase) SelectTxByNumber(number uint32) (*types.Transaction, error) {
	data, err := barrelDB.GetTable(NumberTxTableName).Get([]byte(strconv.Itoa(int(number))))
	if err != nil {
		if err.Error() != common.LevelDBNotFoundError {
			return nil, err
		}
		return nil, nil
	}

	tx := new(types.Transaction)
	err = tx.Decode(types.NewGobTxDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (barrelDB *BarrelDatabase) SelectLastTx() (*types.Transaction, error) {
	data, err := barrelDB.GetTable(LastTxTableName).Get([]byte{})
	if err != nil {
		if err.Error() != common.LevelDBNotFoundError {
			return nil, err
		}
		return nil, nil
	}

	tx := new(types.Transaction)
	err = tx.Decode(types.NewGobTxDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (barrelDB *BarrelDatabase) SelectLastTxNumber() (*uint32, error) {
	data, err := barrelDB.GetTable(LastTxNumberTableName).Get([]byte{})
	if err != nil {
		if err.Error() != common.LevelDBNotFoundError {
			return nil, err
		}
		return nil, nil
	}

	number, err := strconv.Atoi(hex.EncodeToString(data))
	if err != nil {
		return nil, err
	}

	num := uint32(number)
	return &num, nil
}
