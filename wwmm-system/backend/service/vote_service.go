package service

import (
	"encoding/json"
	"errors"
	"wwmm/dao"
	"wwmm/model"
	"wwmm/utils"
)

func CastVote(uid, photoID int) (string, error) {
	p, err := dao.GetPhotoByID(photoID)
	if err != nil {
		return "", errors.New("作品不存在")
	}
	if p.Status != model.PhotoStatusApproved {
		return "", errors.New("作品未通过审核，不能投票")
	}
	if p.PhotographerID == uid {
		return "", errors.New("不能给自己的作品投票")
	}
	if has, _ := dao.HasVoted(uid, photoID); has {
		return "", errors.New("您已经投过票")
	}

	u, err := dao.GetUserByID(uid)
	if err != nil {
		return "", errors.New("用户不存在")
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"photoId":    photoID,
		"photoTitle": p.Title,
		"imageHash":  p.ImageHash,
		"voterId":    uid,
		"voterName":  u.Username,
		"voteTime":   utils.NowTs(),
	})
	txHash := utils.GenTxHash(model.TxTypeVote, u.Username, string(payload), utils.NowTs())
	vote := &model.Vote{
		UserID:  uid,
		PhotoID: photoID,
		TxHash:  txHash,
	}
	voteID, err := dao.CreateVote(vote)
	if err != nil {
		return "", err
	}

	tx := &model.Transaction{
		TxHash:  txHash,
		TxType:  model.TxTypeVote,
		Sender:  u.Username,
		Payload: string(payload),
		Status:  model.TxStatusPending,
	}
	_, _ = dao.CreateTx(tx)
	_ = dao.UpdatePhotoVoteCount(photoID)
	_, _ = PackPendingTxs()
	_ = voteID
	return txHash, nil
}
