package main

import (
	"errors"
	"log"
	"wwmm/dao"
	"wwmm/model"
	"wwmm/utils"
)

func seedDefaultUsers() error {
	defaults := []struct {
		username string
		password string
		phone    string
		email    string
		realName string
		role     int
		bio      string
	}{
		{"admin", "admin123", "13800000001", "admin@wwmm.com", "系统管理员", model.RoleAdmin, "平台超级管理员"},
		{"photographer", "photo123", "13800000002", "photo@wwmm.com", "张摄影", model.RolePhotographer, "职业风光摄影师"},
		{"voter", "vote123", "13800000003", "voter@wwmm.com", "王观众", model.RoleVoter, "摄影爱好者"},
		{"alice", "alice123", "13800000004", "alice@wwmm.com", "Alice", model.RolePhotographer, "街拍摄影师"},
		{"bob", "bob123", "13800000005", "bob@wwmm.com", "Bob", model.RolePhotographer, "旅行摄影师"},
	}
	for _, d := range defaults {
		c, _ := dao.CountUserByUsername(d.username)
		if c > 0 {
			continue
		}
		salt := utils.GenSalt()
		ph := utils.HashPassword(d.password, salt)
		u := &model.User{
			Username:     d.username,
			PasswordHash: ph,
			Salt:         salt,
			Phone:        d.phone,
			Email:        d.email,
			RealName:     d.realName,
			Role:         d.role,
			Status:       model.StatusEnabled,
			Bio:          d.bio,
		}
		if _, err := dao.CreateUser(u); err != nil {
			log.Printf("创建种子用户 %s 失败: %v", d.username, err)
			continue
		}
		log.Printf("[SEED] 用户 %s / %s (role=%d)", d.username, d.password, d.role)
	}
	return nil
}

var _ = errors.New
