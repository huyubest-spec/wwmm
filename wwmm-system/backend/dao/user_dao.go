package dao

import (
	"database/sql"
	"errors"
	"wwmm/model"
	"wwmm/utils"
)

func CreateUser(u *model.User) (int, error) {
	res, err := utils.DB.Exec(
		"INSERT INTO `user`(username,password_hash,salt,phone,email,real_name,sex,avatar,bio,role,status) VALUES(?,?,?,?,?,?,?,?,?,?,?)",
		u.Username, u.PasswordHash, u.Salt, u.Phone, u.Email, u.RealName, u.Sex, u.Avatar, u.Bio, u.Role, u.Status)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func GetUserByUsername(username string) (*model.User, error) {
	row := utils.DB.QueryRow(
		"SELECT user_id,username,password_hash,salt,phone,email,real_name,sex,avatar,bio,role,status,created_at,updated_at FROM `user` WHERE username=?",
		username)
	var u model.User
	err := row.Scan(&u.UserID, &u.Username, &u.PasswordHash, &u.Salt, &u.Phone, &u.Email,
		&u.RealName, &u.Sex, &u.Avatar, &u.Bio, &u.Role, &u.Status, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByID(id int) (*model.User, error) {
	row := utils.DB.QueryRow(
		"SELECT user_id,username,password_hash,salt,phone,email,real_name,sex,avatar,bio,role,status,created_at,updated_at FROM `user` WHERE user_id=?",
		id)
	var u model.User
	err := row.Scan(&u.UserID, &u.Username, &u.PasswordHash, &u.Salt, &u.Phone, &u.Email,
		&u.RealName, &u.Sex, &u.Avatar, &u.Bio, &u.Role, &u.Status, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func CountUserByUsername(username string) (int, error) {
	var c int
	row := utils.DB.QueryRow("SELECT COUNT(*) FROM `user` WHERE username=?", username)
	err := row.Scan(&c)
	return c, err
}

func ListUsers(role int, limit, offset int) ([]model.User, error) {
	var rows *sql.Rows
	var err error
	if role >= 0 {
		rows, err = utils.DB.Query("SELECT user_id,username,phone,email,real_name,sex,role,status,created_at FROM `user` WHERE role=? ORDER BY user_id DESC LIMIT ? OFFSET ?",
			role, limit, offset)
	} else {
		rows, err = utils.DB.Query("SELECT user_id,username,phone,email,real_name,sex,role,status,created_at FROM `user` ORDER BY user_id DESC LIMIT ? OFFSET ?",
			limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.UserID, &u.Username, &u.Phone, &u.Email, &u.RealName, &u.Sex, &u.Role, &u.Status, &u.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, nil
}

func CountUsers(role int) (int, error) {
	var c int
	var err error
	if role >= 0 {
		err = utils.DB.QueryRow("SELECT COUNT(*) FROM `user` WHERE role=?", role).Scan(&c)
	} else {
		err = utils.DB.QueryRow("SELECT COUNT(*) FROM `user`").Scan(&c)
	}
	return c, err
}

func UpdateUserStatus(id int, status int) error {
	_, err := utils.DB.Exec("UPDATE `user` SET status=? WHERE user_id=?", status, id)
	return err
}

var ErrNotFound = errors.New("not found")
