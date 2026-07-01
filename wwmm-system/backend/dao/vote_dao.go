package dao

import (
	"wwmm/model"
	"wwmm/utils"
)

func CreateVote(v *model.Vote) (int, error) {
	res, err := utils.DB.Exec(
		"INSERT INTO vote(user_id,photo_id,tx_hash) VALUES(?,?,?)",
		v.UserID, v.PhotoID, v.TxHash)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func HasVoted(uid, pid int) (bool, error) {
	var c int
	err := utils.DB.QueryRow("SELECT COUNT(*) FROM vote WHERE user_id=? AND photo_id=?", uid, pid).Scan(&c)
	return c > 0, err
}

func CountVotesByPhoto(pid int) (int, error) {
	var c int
	err := utils.DB.QueryRow("SELECT COUNT(*) FROM vote WHERE photo_id=?", pid).Scan(&c)
	return c, err
}

func UpdateVoteTxHash(voteID int, txHash string) error {
	_, err := utils.DB.Exec("UPDATE vote SET tx_hash=? WHERE vote_id=?", txHash, voteID)
	return err
}

func GetVoteByUserAndPhoto(uid, pid int) (*model.Vote, error) {
	row := utils.DB.QueryRow("SELECT vote_id,user_id,photo_id,IFNULL(tx_hash,''),created_at FROM vote WHERE user_id=? AND photo_id=?", uid, pid)
	var v model.Vote
	err := row.Scan(&v.VoteID, &v.UserID, &v.PhotoID, &v.TxHash, &v.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
