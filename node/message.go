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

type GetStatusMessage struct {
}

type StatusMessage struct {
	// the id of the server
	ID            string
	Version       uint32
	CurrentHeight int32
}
