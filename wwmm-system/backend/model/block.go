package model

import "time"

type Block struct {
	Index      int       `json:"index"`
	PrevHash   string    `json:"prevHash"`
	MerkleRoot string    `json:"merkleRoot"`
	Timestamp  int64     `json:"timestamp"`
	Nonce      int64     `json:"nonce"`
	Difficulty int       `json:"difficulty"`
	Hash       string    `json:"hash"`
	TxCount    int       `json:"txCount"`
	Miner      string    `json:"miner"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Transaction struct {
	TxID      int       `json:"txId"`
	TxHash    string    `json:"txHash"`
	BlockID   *int      `json:"blockId"`
	TxType    int       `json:"txType"`
	Sender    string    `json:"sender"`
	Payload   string    `json:"payload"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type ChainState struct {
	LatestIndex int    `json:"latestIndex"`
	LatestHash  string `json:"latestHash"`
	TotalTx     int    `json:"totalTx"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

const (
	TxTypePhotoCertify = 1
	TxTypeVote         = 2

	TxStatusPending  = 0
	TxStatusPackaged = 1
	TxStatusFailed   = 2
)
