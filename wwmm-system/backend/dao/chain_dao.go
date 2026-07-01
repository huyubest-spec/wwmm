package dao

import (
	"wwmm/model"
	"wwmm/utils"
)

func GetLatestBlockIndex() (int, error) {
	var idx int
	err := utils.DB.QueryRow("SELECT IFNULL(MAX(index_num),0) FROM block").Scan(&idx)
	return idx, err
}

func GetBlockByIndex(idx int) (*model.Block, error) {
	row := utils.DB.QueryRow(
		"SELECT index_num,prev_hash,merkle_root,timestamp,nonce,difficulty,hash,tx_count,IFNULL(miner,''),created_at FROM block WHERE index_num=?",
		idx)
	var b model.Block
	err := row.Scan(&b.Index, &b.PrevHash, &b.MerkleRoot, &b.Timestamp, &b.Nonce, &b.Difficulty, &b.Hash, &b.TxCount, &b.Miner, &b.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func GetBlockByHash(hash string) (*model.Block, error) {
	row := utils.DB.QueryRow(
		"SELECT index_num,prev_hash,merkle_root,timestamp,nonce,difficulty,hash,tx_count,IFNULL(miner,''),created_at FROM block WHERE hash=?",
		hash)
	var b model.Block
	err := row.Scan(&b.Index, &b.PrevHash, &b.MerkleRoot, &b.Timestamp, &b.Nonce, &b.Difficulty, &b.Hash, &b.TxCount, &b.Miner, &b.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func ListBlocks(limit, offset int) ([]model.Block, error) {
	rows, err := utils.DB.Query(
		"SELECT index_num,prev_hash,merkle_root,timestamp,nonce,difficulty,hash,tx_count,IFNULL(miner,''),created_at FROM block ORDER BY index_num DESC LIMIT ? OFFSET ?",
		limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []model.Block
	for rows.Next() {
		var b model.Block
		if err := rows.Scan(&b.Index, &b.PrevHash, &b.MerkleRoot, &b.Timestamp, &b.Nonce, &b.Difficulty, &b.Hash, &b.TxCount, &b.Miner, &b.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, nil
}

func CountBlocks() (int, error) {
	var c int
	err := utils.DB.QueryRow("SELECT COUNT(*) FROM block").Scan(&c)
	return c, err
}

func InsertBlock(b *model.Block) (int, error) {
	res, err := utils.DB.Exec(
		"INSERT INTO block(index_num,prev_hash,merkle_root,timestamp,nonce,difficulty,hash,tx_count,miner) VALUES(?,?,?,?,?,?,?,?,?)",
		b.Index, b.PrevHash, b.MerkleRoot, b.Timestamp, b.Nonce, b.Difficulty, b.Hash, b.TxCount, b.Miner)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func CreateTx(tx *model.Transaction) (int, error) {
	res, err := utils.DB.Exec(
		"INSERT INTO tx(tx_hash,block_id,tx_type,sender,payload,status) VALUES(?,?,?,?,?,?)",
		tx.TxHash, tx.BlockID, tx.TxType, tx.Sender, tx.Payload, tx.Status)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func UpdateTxBlockAndStatus(txHash string, blockID int, status int) error {
	_, err := utils.DB.Exec("UPDATE tx SET block_id=?, status=? WHERE tx_hash=?", blockID, status, txHash)
	return err
}

func GetTxByHash(hash string) (*model.Transaction, error) {
	row := utils.DB.QueryRow(
		"SELECT tx_id,tx_hash,block_id,tx_type,sender,payload,status,created_at FROM tx WHERE tx_hash=?",
		hash)
	var t model.Transaction
	var blockID *int
	err := row.Scan(&t.TxID, &t.TxHash, &blockID, &t.TxType, &t.Sender, &t.Payload, &t.Status, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	t.BlockID = blockID
	return &t, nil
}

func ListTxsByBlock(blockID int) ([]model.Transaction, error) {
	rows, err := utils.DB.Query(
		"SELECT tx_id,tx_hash,block_id,tx_type,sender,payload,status,created_at FROM tx WHERE block_id=? ORDER BY tx_id",
		blockID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []model.Transaction
	for rows.Next() {
		var t model.Transaction
		var blockID *int
		if err := rows.Scan(&t.TxID, &t.TxHash, &blockID, &t.TxType, &t.Sender, &t.Payload, &t.Status, &t.CreatedAt); err != nil {
			return nil, err
		}
		t.BlockID = blockID
		list = append(list, t)
	}
	return list, nil
}

func ListTxs(limit, offset int) ([]model.Transaction, error) {
	rows, err := utils.DB.Query(
		"SELECT tx_id,tx_hash,block_id,tx_type,sender,payload,status,created_at FROM tx ORDER BY tx_id DESC LIMIT ? OFFSET ?",
		limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []model.Transaction
	for rows.Next() {
		var t model.Transaction
		var blockID *int
		if err := rows.Scan(&t.TxID, &t.TxHash, &blockID, &t.TxType, &t.Sender, &t.Payload, &t.Status, &t.CreatedAt); err != nil {
			return nil, err
		}
		t.BlockID = blockID
		list = append(list, t)
	}
	return list, nil
}

func CountTxs() (int, error) {
	var c int
	err := utils.DB.QueryRow("SELECT COUNT(*) FROM tx").Scan(&c)
	return c, err
}

func CountTxsByType(txType int) (int, error) {
	var c int
	err := utils.DB.QueryRow("SELECT COUNT(*) FROM tx WHERE tx_type=?", txType).Scan(&c)
	return c, err
}

func GetChainState() (*model.ChainState, error) {
	row := utils.DB.QueryRow("SELECT latest_index,IFNULL(latest_hash,''),total_tx,updated_at FROM chain_state WHERE id=1")
	var s model.ChainState
	err := row.Scan(&s.LatestIndex, &s.LatestHash, &s.TotalTx, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func UpdateChainState(index int, hash string, totalTx int) error {
	_, err := utils.DB.Exec("UPDATE chain_state SET latest_index=?, latest_hash=?, total_tx=? WHERE id=1", index, hash, totalTx)
	return err
}

func InsertAuditLog(photoID, adminID int, action, comment string) error {
	_, err := utils.DB.Exec(
		"INSERT INTO photo_audit_log(photo_id,admin_id,action,comment) VALUES(?,?,?,?)",
		photoID, adminID, action, comment)
	return err
}
