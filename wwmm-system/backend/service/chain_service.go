package service

import (
	"encoding/json"
	"errors"
	"wwmm/blockchain"
	"wwmm/config"
	"wwmm/dao"
	"wwmm/model"
	"wwmm/utils"
)

func PackPendingTxs() (*blockchain.Block, error) {
	rows, err := utils.DB.Query(
		"SELECT tx_id, tx_hash, block_id, tx_type, sender, payload, status, created_at FROM tx WHERE status=? ORDER BY tx_id",
		model.TxStatusPending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pending []model.Transaction
	for rows.Next() {
		var t model.Transaction
		var blockID *int
		if err := rows.Scan(&t.TxID, &t.TxHash, &blockID, &t.TxType, &t.Sender, &t.Payload, &t.Status, &t.CreatedAt); err != nil {
			return nil, err
		}
		t.BlockID = blockID
		pending = append(pending, t)
	}
	if len(pending) == 0 {
		return nil, nil
	}
	latestIdx, _ := dao.GetLatestBlockIndex()
	var prevHash string
	if latestIdx > 0 {
		if b, err := dao.GetBlockByIndex(latestIdx); err == nil {
			prevHash = b.Hash
		}
	}
	var txs []blockchain.Transaction
	for _, t := range pending {
		txs = append(txs, blockchain.Transaction{
			TxID:      t.TxID,
			TxHash:    t.TxHash,
			TxType:    t.TxType,
			Sender:    t.Sender,
			Payload:   t.Payload,
			CreatedAt: t.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	newBlock := blockchain.PackBlock(latestIdx+1, prevHash, txs, config.App.ChainDiff, "wwmm-node-01")
	blockID, err := dao.InsertBlock(&model.Block{
		Index:      newBlock.Index,
		PrevHash:   newBlock.PrevHash,
		MerkleRoot: newBlock.MerkleRoot,
		Timestamp:  newBlock.Timestamp,
		Nonce:      newBlock.Nonce,
		Difficulty: newBlock.Difficulty,
		Hash:       newBlock.Hash,
		TxCount:    newBlock.TxCount,
		Miner:      newBlock.Miner,
	})
	if err != nil {
		return nil, err
	}
	for _, t := range pending {
		_ = dao.UpdateTxBlockAndStatus(t.TxHash, blockID, model.TxStatusPackaged)
		if t.TxType == model.TxTypePhotoCertify {
			markPhotoOnChainByPayload(t.Payload, t.TxHash)
		} else if t.TxType == model.TxTypeVote {
			_, _ = utils.DB.Exec("UPDATE vote SET tx_hash=? WHERE tx_hash=?", t.TxHash, t.TxHash)
		}
	}
	total, _ := dao.CountTxs()
	_ = dao.UpdateChainState(newBlock.Index, newBlock.Hash, total)
	return newBlock, nil
}

func markPhotoOnChainByPayload(payload, txHash string) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &m); err != nil {
		return
	}
	if v, ok := m["photoId"]; ok {
		switch x := v.(type) {
		case float64:
			_ = dao.UpdatePhotoOnChain(int(x), txHash)
		}
	}
}

func InitGenesis() error {
	cnt, _ := dao.CountBlocks()
	if cnt > 0 {
		return nil
	}
	genesis := blockchain.PackBlock(1, "0000000000000000000000000000000000000000000000000000000000000000", nil, config.App.ChainDiff, "wwmm-genesis")
	_, err := dao.InsertBlock(&model.Block{
		Index:      genesis.Index,
		PrevHash:   genesis.PrevHash,
		MerkleRoot: genesis.MerkleRoot,
		Timestamp:  genesis.Timestamp,
		Nonce:      genesis.Nonce,
		Difficulty: genesis.Difficulty,
		Hash:       genesis.Hash,
		TxCount:    genesis.TxCount,
		Miner:      genesis.Miner,
	})
	if err != nil {
		return err
	}
	_ = dao.UpdateChainState(genesis.Index, genesis.Hash, 0)
	return nil
}

var ErrInvalidHash = errors.New("invalid hash format")

func GetBlockByIndex(idx int) (*model.Block, error) {
	return dao.GetBlockByIndex(idx)
}

func GetBlockByHash(hash string) (*model.Block, error) {
	if len(hash) != 64 {
		return nil, ErrInvalidHash
	}
	return dao.GetBlockByHash(hash)
}

func GetTxByHash(hash string) (*model.Transaction, error) {
	if len(hash) != 64 {
		return nil, ErrInvalidHash
	}
	return dao.GetTxByHash(hash)
}
