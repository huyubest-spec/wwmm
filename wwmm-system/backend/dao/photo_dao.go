package dao

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"

	"wwmm/model"
	"wwmm/utils"
)

func CreatePhoto(p *model.Photo) (int, error) {
	var shootTime interface{}
	if s := strings.TrimSpace(p.ShootTime); s != "" {
		shootTime = s
	}
	res, err := utils.DB.Exec(
		"INSERT INTO photo(title,description,image_url,image_hash,file_size,photographer_id,category,shoot_location,shoot_time,camera_info,status,vote_count,view_count,is_on_chain) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		p.Title, p.Description, p.ImageURL, p.ImageHash, p.FileSize, p.PhotographerID, p.Category, p.ShootLocation, shootTime, p.CameraInfo, p.Status, p.VoteCount, p.ViewCount, p.IsOnChain)
	if err != nil {
		if isDuplicateImageHashErr(err) {
			return 0, errors.New("该图片已存在（哈希冲突），请勿重复上传")
		}
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

// isDuplicateImageHashErr 识别 MySQL 1062 且唯一键为 uk_image_hash 的错误。
// 并发/旧 service 漏过查重时，DB 兜底翻译成业务错误，避免 raw SQL 漏到前端。
func isDuplicateImageHashErr(err error) bool {
	var me *mysql.MySQLError
	if !errors.As(err, &me) {
		return false
	}
	if me.Number != 1062 {
		return false
	}
	return strings.Contains(me.Message, "uk_image_hash")
}

// scanPhotoFields 把可能为 NULL 的字符串列安全地扫进 *model.Photo。
// 返回的描述/分类/拍摄时间/相机/审核意见/链上交易哈希均为空串而非报错。
func scanPhotoFields(p *model.Photo,
	desc, category, shootLocation, shootTime, cameraInfo, auditComment, chainTxHash sql.NullString) {
	if desc.Valid {
		p.Description = desc.String
	}
	if category.Valid {
		p.Category = category.String
	}
	if shootLocation.Valid {
		p.ShootLocation = shootLocation.String
	}
	if shootTime.Valid {
		p.ShootTime = shootTime.String
	}
	if cameraInfo.Valid {
		p.CameraInfo = cameraInfo.String
	}
	if auditComment.Valid {
		p.AuditComment = auditComment.String
	}
	if chainTxHash.Valid {
		p.ChainTxHash = chainTxHash.String
	}
}

func GetPhotoByID(id int) (*model.Photo, error) {
	row := utils.DB.QueryRow(
		"SELECT photo_id,title,description,image_url,image_hash,file_size,photographer_id,category,shoot_location,shoot_time,camera_info,status,audit_comment,vote_count,view_count,is_on_chain,chain_tx_hash,created_at,updated_at FROM photo WHERE photo_id=?",
		id)
	var p model.Photo
	var desc, category, shootLocation, shootTime, cameraInfo, auditComment, chainTxHash sql.NullString
	err := row.Scan(&p.PhotoID, &p.Title, &desc, &p.ImageURL, &p.ImageHash, &p.FileSize, &p.PhotographerID,
		&category, &shootLocation, &shootTime, &cameraInfo, &p.Status, &auditComment, &p.VoteCount, &p.ViewCount,
		&p.IsOnChain, &chainTxHash, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	scanPhotoFields(&p, desc, category, shootLocation, shootTime, cameraInfo, auditComment, chainTxHash)
	return &p, nil
}

func GetPhotoByHash(hash string) (*model.Photo, error) {
	row := utils.DB.QueryRow(
		"SELECT photo_id,title,description,image_url,image_hash,file_size,photographer_id,category,shoot_location,shoot_time,camera_info,status,audit_comment,vote_count,view_count,is_on_chain,chain_tx_hash,created_at,updated_at FROM photo WHERE image_hash=?",
		hash)
	var p model.Photo
	var desc, category, shootLocation, shootTime, cameraInfo, auditComment, chainTxHash sql.NullString
	err := row.Scan(&p.PhotoID, &p.Title, &desc, &p.ImageURL, &p.ImageHash, &p.FileSize, &p.PhotographerID,
		&category, &shootLocation, &shootTime, &cameraInfo, &p.Status, &auditComment, &p.VoteCount, &p.ViewCount,
		&p.IsOnChain, &chainTxHash, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	scanPhotoFields(&p, desc, category, shootLocation, shootTime, cameraInfo, auditComment, chainTxHash)
	return &p, nil
}

func ListPhotosByStatus(status int, limit, offset int) ([]model.PhotoFull, error) {
	rows, err := utils.DB.Query(`
		SELECT p.photo_id,p.title,p.description,p.image_url,p.image_hash,p.file_size,p.photographer_id,
		       p.category,p.shoot_location,p.shoot_time,p.camera_info,p.status,p.audit_comment,
		       p.vote_count,p.view_count,p.is_on_chain,p.chain_tx_hash,p.created_at,p.updated_at,
		       u.username,u.avatar
		FROM photo p LEFT JOIN `+"`user`"+` u ON p.photographer_id=u.user_id
		WHERE p.status=?
		ORDER BY p.vote_count DESC, p.photo_id DESC
		LIMIT ? OFFSET ?`,
		status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []model.PhotoFull
	for rows.Next() {
		var p model.PhotoFull
		var desc, category, shootLocation, shootTime, cameraInfo, auditComment, chainTxHash sql.NullString
		if err := rows.Scan(&p.PhotoID, &p.Title, &desc, &p.ImageURL, &p.ImageHash, &p.FileSize,
			&p.PhotographerID, &category, &shootLocation, &shootTime, &cameraInfo,
			&p.Status, &auditComment, &p.VoteCount, &p.ViewCount, &p.IsOnChain, &chainTxHash,
			&p.CreatedAt, &p.UpdatedAt, &p.PhotographerName, &p.PhotographerAvatar); err != nil {
			return nil, err
		}
		scanPhotoFields(&p.Photo, desc, category, shootLocation, shootTime, cameraInfo, auditComment, chainTxHash)
		list = append(list, p)
	}
	return list, nil
}

func ListPhotosByPhotographer(uid int, limit, offset int) ([]model.Photo, error) {
	rows, err := utils.DB.Query(
		"SELECT photo_id,title,description,image_url,image_hash,file_size,photographer_id,category,shoot_location,shoot_time,camera_info,status,audit_comment,vote_count,view_count,is_on_chain,chain_tx_hash,created_at,updated_at FROM photo WHERE photographer_id=? ORDER BY photo_id DESC LIMIT ? OFFSET ?",
		uid, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []model.Photo
	for rows.Next() {
		var p model.Photo
		var desc, category, shootLocation, shootTime, cameraInfo, auditComment, chainTxHash sql.NullString
		if err := rows.Scan(&p.PhotoID, &p.Title, &desc, &p.ImageURL, &p.ImageHash, &p.FileSize,
			&p.PhotographerID, &category, &shootLocation, &shootTime, &cameraInfo,
			&p.Status, &auditComment, &p.VoteCount, &p.ViewCount, &p.IsOnChain, &chainTxHash,
			&p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		scanPhotoFields(&p, desc, category, shootLocation, shootTime, cameraInfo, auditComment, chainTxHash)
		list = append(list, p)
	}
	return list, nil
}

func CountPhotosByStatus(status int) (int, error) {
	var c int
	err := utils.DB.QueryRow("SELECT COUNT(*) FROM photo WHERE status=?", status).Scan(&c)
	return c, err
}

func UpdatePhotoStatus(id int, status int, comment string) error {
	_, err := utils.DB.Exec("UPDATE photo SET status=?, audit_comment=? WHERE photo_id=?", status, comment, id)
	return err
}

func UpdatePhotoVoteCount(id int) error {
	_, err := utils.DB.Exec("UPDATE photo SET vote_count = (SELECT COUNT(*) FROM vote WHERE photo_id=?) WHERE photo_id=?", id, id)
	return err
}

func UpdatePhotoViewCount(id int) error {
	_, err := utils.DB.Exec("UPDATE photo SET view_count = view_count + 1 WHERE photo_id=?", id)
	return err
}

func UpdatePhotoOnChain(id int, txHash string) error {
	_, err := utils.DB.Exec("UPDATE photo SET is_on_chain=1, chain_tx_hash=? WHERE photo_id=?", txHash, id)
	return err
}
