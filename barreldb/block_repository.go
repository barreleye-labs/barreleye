package barreldb

import (
	"bytes"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
	"strconv"
)

func (barrelDB *BarrelDatabase) CreateBlockWithHash(hash common.Hash, block *types.Block) error {
	buf := &bytes.Buffer{}
	if err := block.Encode(types.NewGobBlockEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable("hash-block").Put(hash.ToSlice(), buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) CreateBlockWithHeight(height uint32, block *types.Block) error {
	buf := &bytes.Buffer{}
	if err := block.Encode(types.NewGobBlockEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable("height-block").Put([]byte(strconv.Itoa(int(height))), buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) GetBlock(hash common.Hash) (*types.Block, error) {
	data, err := barrelDB.GetTable("hash-block").Get(hash.ToSlice())
	if err != nil {
		return nil, err
	}

	bDecode := new(types.Block)
	err = bDecode.Decode(types.NewGobBlockDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	return bDecode, nil
}

func (barrelDB *BarrelDatabase) GetBlockByHeight(height uint32) (*types.Block, error) {
	data, err := barrelDB.GetTable("height-block").Get([]byte(strconv.Itoa(int(height))))
	if err != nil {
		return nil, err
	}

	bDecode := new(types.Block)
	err = bDecode.Decode(types.NewGobBlockDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	return bDecode, nil
}
