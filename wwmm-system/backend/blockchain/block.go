package blockchain

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"wwmm/utils"
)

type Transaction struct {
	TxID      int    `json:"txId"`
	TxHash    string `json:"txHash"`
	TxType    int    `json:"txType"`
	Sender    string `json:"sender"`
	Payload   string `json:"payload"`
	CreatedAt string `json:"createdAt"`
}

type Block struct {
	Index      int           `json:"index"`
	PrevHash   string        `json:"prevHash"`
	MerkleRoot string        `json:"merkleRoot"`
	Timestamp  int64         `json:"timestamp"`
	Nonce      int64         `json:"nonce"`
	Difficulty int           `json:"difficulty"`
	Hash       string        `json:"hash"`
	TxCount    int           `json:"txCount"`
	Miner      string        `json:"miner"`
	TxList     []Transaction `json:"txList"`
}

func ComputeBlockHash(b *Block) string {
	raw := fmt.Sprintf("%d|%s|%s|%d|%d|%d|%d|%s",
		b.Index, b.PrevHash, b.MerkleRoot, b.Timestamp, b.Nonce, b.Difficulty, b.TxCount, b.Miner)
	return utils.Sha256HexString(raw)
}

func (b *Block) Mine() {
	prefix := strings.Repeat("0", b.Difficulty)
	for {
		b.Hash = ComputeBlockHash(b)
		if strings.HasPrefix(b.Hash, prefix) {
			return
		}
		b.Nonce++
	}
}

func PackBlock(index int, prevHash string, txs []Transaction, difficulty int, miner string) *Block {
	leafs := make([]string, 0, len(txs))
	for _, t := range txs {
		leafs = append(leafs, t.TxHash)
	}
	merkle := MerkleRoot(leafs)
	b := &Block{
		Index:      index,
		PrevHash:   prevHash,
		MerkleRoot: merkle,
		Timestamp:  time.Now().Unix(),
		Nonce:      0,
		Difficulty: difficulty,
		TxCount:    len(txs),
		Miner:      miner,
		TxList:     txs,
	}
	b.Mine()
	return b
}

func TxPayloadPretty(p string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(p), &m); err != nil {
		return p
	}
	b, _ := json.MarshalIndent(m, "", "  ")
	return string(b)
}
