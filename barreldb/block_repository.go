package barreldb

import (
	"bytes"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
	"strconv"
)

// HashBlock Repository
func (barrelDB *BarrelDatabase) InsertHashBlock(hash common.Hash, block *types.Block) error {
	buf := &bytes.Buffer{}
	if err := block.Encode(types.NewGobBlockEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable(HashBlockTableName).Put(hash.ToSlice(), buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) DeleteHashBlock(hash common.Hash) error {
	if err := barrelDB.GetTable(HashBlockTableName).Delete(hash.ToSlice()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) SelectHashBlock(hash common.Hash) (*types.Block, error) {
	data, err := barrelDB.GetTable(HashBlockTableName).Get(hash.ToSlice())
	if err != nil {
		if err.Error() != common.LevelDBNotFoundError {
			return nil, err
		}
		return nil, nil
	}

	block := new(types.Block)
	err = block.Decode(types.NewGobBlockDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	return block, nil
}

// HeightBlock Repository
func (barrelDB *BarrelDatabase) InsertHeightBlock(height int32, block *types.Block) error {
	buf := &bytes.Buffer{}
	if err := block.Encode(types.NewGobBlockEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable(HeightBlockTableName).Put([]byte(strconv.Itoa(int(height))), buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) DeleteHeightBlock(height int32) error {
	if err := barrelDB.GetTable(HeightBlockTableName).Delete([]byte(strconv.Itoa(int(height)))); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) SelectHeightBlock(height int32) (*types.Block, error) {
	data, err := barrelDB.GetTable(HeightBlockTableName).Get([]byte(strconv.Itoa(int(height))))
	if err != nil {
		if err.Error() != common.LevelDBNotFoundError {
			return nil, err
		}
		return nil, nil
	}

	block := new(types.Block)
	err = block.Decode(types.NewGobBlockDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	return block, nil
}

// LastBlock Repository
func (barrelDB *BarrelDatabase) InsertLastBlock(block *types.Block) error {
	buf := &bytes.Buffer{}
	if err := block.Encode(types.NewGobBlockEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable(LastBlockTableName).Put([]byte{}, buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) DeleteLastBlock() error {
	if err := barrelDB.GetTable(LastBlockTableName).Delete([]byte{}); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) SelectLastBlock() (*types.Block, error) {
	data, err := barrelDB.GetTable(LastBlockTableName).Get([]byte{})
	if err != nil {
		if err.Error() != common.LevelDBNotFoundError {
			return nil, err
		}
		return nil, nil
	}

	block := new(types.Block)
	err = block.Decode(types.NewGobBlockDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	return block, nil
}
