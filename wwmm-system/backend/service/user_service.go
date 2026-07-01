package service

import (
	"errors"
	"strings"
	"unicode"
	"wwmm/dao"
	"wwmm/model"
	"wwmm/utils"
)

func RegisterUser(username, password, phone, email, realName string, sex, role int) (*model.User, error) {
	username = strings.TrimSpace(username)
	if username == "" || len(username) < 3 {
		return nil, errors.New("用户名长度需 ≥ 3")
	}
	if len(password) < 6 {
		return nil, errors.New("密码长度需 ≥ 6")
	}
	if c, _ := dao.CountUserByUsername(username); c > 0 {
		return nil, errors.New("用户名已被占用")
	}
	if role != model.RoleVoter && role != model.RolePhotographer {
		role = model.RoleVoter
	}
	salt := utils.GenSalt()
	ph := utils.HashPassword(password, salt)
	u := &model.User{
		Username:     username,
		PasswordHash: ph,
		Salt:         salt,
		Phone:        phone,
		Email:        email,
		RealName:     realName,
		Sex:          sex,
		Role:         role,
		Status:       model.StatusEnabled,
	}
	id, err := dao.CreateUser(u)
	if err != nil {
		return nil, err
	}
	u.UserID = id
	return u, nil
}

func LoginUser(username, password string) (*model.User, string, error) {
	u, err := dao.GetUserByUsername(username)
	if err != nil {
		return nil, "", errors.New("用户名或密码错误")
	}
	if u.Status == model.StatusDisabled {
		return nil, "", errors.New("账号已被禁用")
	}
	ph := utils.HashPassword(password, u.Salt)
	if ph != u.PasswordHash {
		return nil, "", errors.New("用户名或密码错误")
	}
	tok, err := utils.CreateSession(u.UserID, u.Username, u.Role)
	if err != nil {
		return nil, "", err
	}
	return u, tok, nil
}

func Logout(token string) error {
	return utils.DestroySession(token)
}

func isValidName(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return s != ""
}
