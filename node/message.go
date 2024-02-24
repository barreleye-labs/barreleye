package node

import (
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
