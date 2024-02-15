package barreldb

import (
	"bytes"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
	"strconv"
)

// HashHeader Repository
func (barrelDB *BarrelDatabase) InsertHeaderWithHash(hash common.Hash, header *types.Header) error {
	buf := &bytes.Buffer{}
	if err := header.Encode(types.NewGobHeaderEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable(HashHeaderTableName).Put(hash.ToSlice(), buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) DeleteHeaderByHash(hash common.Hash) error {
	if err := barrelDB.GetTable(HashHeaderTableName).Delete(hash.ToSlice()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) SelectHeaderByHash(hash common.Hash) (*types.Header, error) {
	data, err := barrelDB.GetTable(HashHeaderTableName).Get(hash.ToSlice())
	if err != nil {
		if err.Error() != common.LevelDBNotFoundError {
			return nil, err
		}
		return nil, nil
	}

	header := new(types.Header)
	err = header.Decode(types.NewGobHeaderDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	return header, nil
}

// HeightHeader Repository
func (barrelDB *BarrelDatabase) InsertHeaderWithHeight(height int32, header *types.Header) error {
	buf := &bytes.Buffer{}
	if err := header.Encode(types.NewGobHeaderEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable(HeightHeaderTableName).Put([]byte(strconv.Itoa(int(height))), buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) DeleteHeaderByHeight(height int32) error {
	if err := barrelDB.GetTable(HeightHeaderTableName).Delete([]byte(strconv.Itoa(int(height)))); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) SelectHeaderByHeight(height int32) (*types.Header, error) {
	data, err := barrelDB.GetTable(HeightHeaderTableName).Get([]byte(strconv.Itoa(int(height))))
	if err != nil {
		if err.Error() != common.LevelDBNotFoundError {
			return nil, err
		}
		return nil, nil
	}

	header := new(types.Header)
	err = header.Decode(types.NewGobHeaderDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	return header, nil
}

// LastHeader Repository
func (barrelDB *BarrelDatabase) InsertLastHeader(header *types.Header) error {
	buf := &bytes.Buffer{}
	if err := header.Encode(types.NewGobHeaderEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable(LastHeaderTableName).Put([]byte{}, buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) DeleteLastHeader() error {
	if err := barrelDB.GetTable(LastHeaderTableName).Delete([]byte{}); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) SelectLastHeader() (*types.Header, error) {
	data, err := barrelDB.GetTable(LastHeaderTableName).Get([]byte{})
	if err != nil {
		if err.Error() != common.LevelDBNotFoundError {
			return nil, err
		}
		return nil, nil
	}

	header := new(types.Header)
	err = header.Decode(types.NewGobHeaderDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	return header, nil
}
