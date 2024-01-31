package node

import (
	"github.com/barreleye-labs/barreleye/core/types"
)

type GetBlocksMessage struct {
	From uint32
	// If to is 0 the maximum blocks will be returned.
	To uint32
}

type BlocksMessage struct {
	Blocks []*types.Block
}

type GetStatusMessage struct {
}

type StatusMessage struct {
	// the id of the server
	ID            string
	Version       uint32
	CurrentHeight uint32
}
