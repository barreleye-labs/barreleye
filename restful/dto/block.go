package dto

type Block struct {
	Hash          string    `json:"hash"`
	Version       uint32    `json:"version"`
	DataHash      string    `json:"dataHash"`
	PrevBlockHash string    `json:"prevBlockHash"`
	Height        int32     `json:"height"`
	Timestamp     int64     `json:"timestamp"`
	Signer        string    `json:"signer"`
	Extra         string    `json:"extra"`
	Signature     Signature `json:"signature"`
	TxCount       uint32    `json:"txCount"`
	Transactions  []string  `json:"transactions"`
}

func CreateBlock(
	hash string,
	version uint32,
	dataHash string,
	prevBlockHash string,
	height int32,
	timestamp int64,
	signer string,
	extra string,
	signature Signature,
	txCount uint32,
	transactions []string) Block {
	return Block{
		Hash:          hash,
		Version:       version,
		DataHash:      dataHash,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     timestamp,
		Signer:        signer,
		Extra:         extra,
		Signature:     signature,
		TxCount:       txCount,
		Transactions:  transactions,
	}
}

type BlockResponse struct {
	Block Block `json:"block"`
}

func CreateBlockResponse(block Block) BlockResponse {
	return BlockResponse{
		Block: block,
	}
}

type BlocksResponse struct {
	Blocks     []Block `json:"blocks"`
	TotalCount uint32  `json:"totalCount"`
}

func CreateBlocksResponse(blocks []Block, totalCount uint32) BlocksResponse {
	return BlocksResponse{
		Blocks:     blocks,
		TotalCount: totalCount,
	}
}
