package service

import (
	"encoding/json"
	"errors"
	"wwmm/dao"
	"wwmm/model"
	"wwmm/utils"
)

func SubmitPhoto(photographerID int, p *model.Photo) (*model.Photo, error) {
	if p.Title == "" {
		return nil, errors.New("作品标题不能为空")
	}
	if p.ImageURL == "" || p.ImageHash == "" {
		return nil, errors.New("图片未上传")
	}
	if existing, _ := dao.GetPhotoByHash(p.ImageHash); existing != nil {
		return nil, errors.New("该图片已存在（哈希冲突），请勿重复上传")
	}
	p.PhotographerID = photographerID
	p.Status = model.PhotoStatusPending
	p.VoteCount = 0
	p.ViewCount = 0
	p.IsOnChain = 0
	id, err := dao.CreatePhoto(p)
	if err != nil {
		return nil, err
	}
	p.PhotoID = id

	// 同步上链 (作存证)
	payload, _ := json.Marshal(map[string]interface{}{
		"photoId":        p.PhotoID,
		"title":          p.Title,
		"imageHash":      p.ImageHash,
		"photographer":   photographerID,
		"submitTime":     utils.NowTs(),
		"description":    p.Description,
		"category":       p.Category,
		"shootLocation":  p.ShootLocation,
		"cameraInfo":     p.CameraInfo,
	})
	txHash := utils.GenTxHash(model.TxTypePhotoCertify, "photographer", string(payload), utils.NowTs())
	tx := &model.Transaction{
		TxHash:  txHash,
		TxType:  model.TxTypePhotoCertify,
		Sender:  photographerUsername(photographerID),
		Payload: string(payload),
		Status:  model.TxStatusPending,
	}
	txID, _ := dao.CreateTx(tx)

	// 立即打包区块
	block, err := PackPendingTxs()
	if err != nil {
		return p, err
	}
	_ = txID
	_ = block
	return p, nil
}

func photographerUsername(uid int) string {
	u, err := dao.GetUserByID(uid)
	if err != nil {
		return "unknown"
	}
	return u.Username
}

func ListApprovedPhotos(limit, offset int) ([]model.PhotoFull, error) {
	return dao.ListPhotosByStatus(model.PhotoStatusApproved, limit, offset)
}

func ListPendingPhotos(limit, offset int) ([]model.PhotoFull, error) {
	return dao.ListPhotosByStatus(model.PhotoStatusPending, limit, offset)
}

func ListMyPhotos(uid, limit, offset int) ([]model.Photo, error) {
	return dao.ListPhotosByPhotographer(uid, limit, offset)
}

func GetPhotoDetail(id int, viewerID int) (*model.PhotoFull, error) {
	p, err := dao.GetPhotoByID(id)
	if err != nil {
		return nil, err
	}
	full := &model.PhotoFull{Photo: *p}
	if u, err := dao.GetUserByID(p.PhotographerID); err == nil {
		full.PhotographerName = u.Username
		full.PhotographerAvatar = u.Avatar
	}
	if viewerID > 0 {
		has, _ := dao.HasVoted(viewerID, id)
		full.HasVoted = has
	}
	_ = dao.UpdatePhotoViewCount(id)
	return full, nil
}

func AuditPhoto(photoID, adminID int, approve bool, comment string) error {
	status := model.PhotoStatusRejected
	if approve {
		status = model.PhotoStatusApproved
	}
	if err := dao.UpdatePhotoStatus(photoID, status, comment); err != nil {
		return err
	}
	action := "REJECT"
	if approve {
		action = "APPROVE"
	}
	_ = dao.InsertAuditLog(photoID, adminID, action, comment)
	return nil
}
