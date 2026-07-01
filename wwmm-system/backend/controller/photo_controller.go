package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"wwmm/config"
	"wwmm/model"
	"wwmm/service"
	"wwmm/utils"
)

func ListPhotos(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "12"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 12
	}
	list, err := service.ListApprovedPhotos(size, (page-1)*size)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}
	for i := range list {
		if c2, e := c1SessionUserID(c); e == nil {
			list[i].HasVoted, _ = hasVoted(c2, list[i].PhotoID)
		}
	}
	utils.Success(c, gin.H{"list": list, "page": page, "size": size})
}

func ListPendingPhotos(c *gin.Context) {
	if _, ok := utils.MustRole(c, model.RoleAdmin); !ok {
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	list, err := service.ListPendingPhotos(size, (page-1)*size)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}
	utils.Success(c, gin.H{"list": list, "page": page, "size": size})
}

func ListMyPhotos(c *gin.Context) {
	s, ok := utils.MustSession(c)
	if !ok {
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	list, err := service.ListMyPhotos(s.UserID, size, (page-1)*size)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}
	utils.Success(c, gin.H{"list": list, "page": page, "size": size})
}

func PhotoDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if id < 1 {
		utils.Fail(c, 400, "无效的ID")
		return
	}
	uid := 0
	if s, err := utils.ResolveSession(utils.GetToken(c)); err == nil {
		uid = s.UserID
	}
	p, err := service.GetPhotoDetail(id, uid)
	if err != nil {
		utils.Fail(c, 404, "作品不存在")
		return
	}
	utils.Success(c, p)
}

func UploadPhoto(c *gin.Context) {
	s, ok := utils.MustRole(c, model.RolePhotographer, model.RoleAdmin)
	if !ok {
		return
	}
	title := strings.TrimSpace(c.PostForm("title"))
	description := c.PostForm("description")
	category := c.PostForm("category")
	location := c.PostForm("shootLocation")
	shootTime := strings.TrimSpace(c.PostForm("shootTime"))
	cameraInfo := c.PostForm("cameraInfo")
	if shootTime != "" {
		if _, err := time.Parse("2006-01-02", shootTime); err != nil {
			utils.Fail(c, 400, "拍摄时间格式错误，应为 YYYY-MM-DD")
			return
		}
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		utils.Fail(c, 400, "请上传图片文件")
		return
	}
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		utils.Fail(c, 400, "仅支持 JPG/PNG 格式")
		return
	}
	if fileHeader.Size > 10*1024*1024 {
		utils.Fail(c, 400, "图片大小不能超过 10MB")
		return
	}

	_ = os.MkdirAll(config.App.UploadDir, 0755)
	_ = os.MkdirAll(filepath.Join(config.App.UploadDir, "photos"), 0755)

	src, err := fileHeader.Open()
	if err != nil {
		utils.Fail(c, 500, "打开上传文件失败")
		return
	}
	defer src.Close()

	hasher := sha256.New()
	buf, err := io.ReadAll(src)
	if err != nil {
		utils.Fail(c, 500, "读取文件失败")
		return
	}
	hasher.Write(buf)
	hash := hex.EncodeToString(hasher.Sum(nil))

	newName := fmt.Sprintf("photo_%d_%s%s", time.Now().UnixNano(), hash[:12], ext)
	dst := filepath.Join(config.App.UploadDir, "photos", newName)
	if err := os.WriteFile(dst, buf, 0644); err != nil {
		utils.Fail(c, 500, "保存文件失败")
		return
	}

	p := &model.Photo{
		Title:         title,
		Description:   description,
		ImageURL:      "/static/photos/" + newName,
		ImageHash:     hash,
		FileSize:      int(fileHeader.Size),
		Category:      category,
		ShootLocation: location,
		ShootTime:     shootTime,
		CameraInfo:    cameraInfo,
	}
	created, err := service.SubmitPhoto(s.UserID, p)
	if err != nil {
		_ = os.Remove(dst)
		utils.Fail(c, 400, err.Error())
		return
	}
	utils.Success(c, created)
}

type AuditReq struct {
	Approve bool   `json:"approve"`
	Comment string `json:"comment"`
}

func AuditPhoto(c *gin.Context) {
	if _, ok := utils.MustRole(c, model.RoleAdmin); !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if id < 1 {
		utils.Fail(c, 400, "无效ID")
		return
	}
	var req AuditReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误")
		return
	}
	if err := service.AuditPhoto(id, 0, req.Approve, req.Comment); err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}
	utils.SuccessMsg(c, "审核完成")
}

func c1SessionUserID(c *gin.Context) (int, error) {
	s, err := utils.ResolveSession(utils.GetToken(c))
	if err != nil {
		return 0, err
	}
	return s.UserID, nil
}
