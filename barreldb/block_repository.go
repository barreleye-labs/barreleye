package barreldb

import (
	"bytes"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
	"strconv"
)

func (barrelDB *BarrelDatabase) InsertBlockWithHash(hash common.Hash, block *types.Block) error {
	buf := &bytes.Buffer{}
	if err := block.Encode(types.NewGobBlockEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable("hash-block").Put(hash.ToSlice(), buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) InsertBlockWithHeight(height uint32, block *types.Block) error {
	buf := &bytes.Buffer{}
	if err := block.Encode(types.NewGobBlockEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable("height-block").Put([]byte(strconv.Itoa(int(height))), buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) InsertLastBlock(block *types.Block) error {
	buf := &bytes.Buffer{}
	if err := block.Encode(types.NewGobBlockEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable("lastBlock").Put([]byte{}, buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) SelectBlockByHash(hash common.Hash) (*types.Block, error) {
	data, err := barrelDB.GetTable("hash-block").Get(hash.ToSlice())
	if err != nil {
		return nil, err
	}

	block := new(types.Block)
	err = block.Decode(types.NewGobBlockDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (barrelDB *BarrelDatabase) SelectBlockByHeight(height uint32) (*types.Block, error) {
	data, err := barrelDB.GetTable("height-block").Get([]byte(strconv.Itoa(int(height))))
	if err != nil {
		return nil, err
	}

	block := new(types.Block)
	err = block.Decode(types.NewGobBlockDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (barrelDB *BarrelDatabase) SelectLastBlock() (*types.Block, error) {
	data, err := barrelDB.GetTable("lastBlock").Get([]byte{})
	if err != nil {
		return nil, err
	}

	block := new(types.Block)
	err = block.Decode(types.NewGobBlockDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	return block, nil
}
