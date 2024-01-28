package barreldb

import (
	"bytes"
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
)

func (barrelDB *BarrelDatabase) CreateBlock(hash common.Hash, block *types.Block) error {
	buf := &bytes.Buffer{}
	if err := block.Encode(types.NewGobBlockEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable("block").Put(hash.ToSlice(), buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) GetBlock(hash common.Hash) (*types.Block, error) {
	data, err := barrelDB.GetTable("block").Get(hash.ToSlice())
	if err != nil {
		return nil, err
	}

	bDecode := new(types.Block)
	err = bDecode.Decode(types.NewGobBlockDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, fmt.Errorf("fail to decode block")
	}

	return bDecode, nil
}
