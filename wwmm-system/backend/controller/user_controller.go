package controller

import (
	"github.com/gin-gonic/gin"
	"strings"
	"wwmm/service"
	"wwmm/utils"
)

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	RealName string `json:"realName"`
	Sex      int    `json:"sex"`
	Role     int    `json:"role"`
}

func Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误: "+err.Error())
		return
	}
	u, err := service.RegisterUser(req.Username, req.Password, req.Phone, req.Email, req.RealName, req.Sex, req.Role)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	utils.Success(c, u)
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误")
		return
	}
	u, tok, err := service.LoginUser(req.Username, req.Password)
	if err != nil {
		utils.Fail(c, 401, err.Error())
		return
	}
	utils.Success(c, gin.H{
		"token":    tok,
		"userId":   u.UserID,
		"username": u.Username,
		"role":     u.Role,
		"realName": u.RealName,
		"avatar":   u.Avatar,
	})
}

func Logout(c *gin.Context) {
	tok := utils.GetToken(c)
	if err := service.Logout(tok); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	utils.SuccessMsg(c, "已退出登录")
}

func Me(c *gin.Context) {
	s, ok := utils.MustSession(c)
	if !ok {
		return
	}
	utils.Success(c, gin.H{
		"userId":   s.UserID,
		"username": s.Username,
		"role":     s.Role,
	})
}

func Health(c *gin.Context) {
	utils.Success(c, gin.H{"status": "ok"})
}

var _ = strings.TrimSpace
