package node

import (
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
)

type BlockRequestMessage struct {
	Height int32
}

type BlockResponseMessage struct {
	Block *types.Block
}

type ChainInfoRequestMessage struct {
}

type ChainInfoResponseMessage struct {
	To            string
	Version       uint32
	CurrentHeight int32
}

type BlockHashRequestMessage struct {
	Height int32
}

type BlockHashResponseMessage struct {
	Hash          common.Hash
	CurrentHeight int32
}
