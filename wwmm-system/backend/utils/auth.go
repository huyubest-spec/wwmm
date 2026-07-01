package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Session struct {
	SessionID string
	UserID    int
	Username  string
	Role      int
	ExpireAt  time.Time
}

const tokenTTL = 24 * time.Hour

func genToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func CreateSession(userID int, username string, role int) (string, error) {
	sid := genToken()
	tok := genToken()
	expire := time.Now().Add(tokenTTL)
	_, err := DB.Exec(
		"INSERT INTO session(session_id, user_id, token, expire_at) VALUES(?,?,?,?)",
		sid, userID, tok, expire)
	if err != nil {
		return "", err
	}
	return tok + "." + sid, nil
}

func DestroySession(token string) error {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return errors.New("invalid token")
	}
	_, err := DB.Exec("DELETE FROM session WHERE session_id=?", parts[1])
	return err
}

func ResolveSession(token string) (*Session, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid token")
	}
	row := DB.QueryRow(
		"SELECT s.session_id, s.user_id, s.expire_at, u.username, u.role FROM session s JOIN `user` u ON s.user_id=u.user_id WHERE s.session_id=? AND s.token=?",
		parts[1], parts[0])
	var s Session
	var username string
	var role int
	if err := row.Scan(&s.SessionID, &s.UserID, &s.ExpireAt, &username, &role); err != nil {
		return nil, errors.New("session not found")
	}
	if time.Now().After(s.ExpireAt) {
		_, _ = DB.Exec("DELETE FROM session WHERE session_id=?", s.SessionID)
		return nil, errors.New("session expired")
	}
	s.Username = username
	s.Role = role
	return &s, nil
}

func GetToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		return ""
	}
	auth = strings.TrimPrefix(auth, "Bearer ")
	auth = strings.TrimPrefix(auth, "Token ")
	return strings.TrimSpace(auth)
}

func MustSession(c *gin.Context) (*Session, bool) {
	tok := GetToken(c)
	if tok == "" {
		Unauthorized(c, "未登录")
		return nil, false
	}
	s, err := ResolveSession(tok)
	if err != nil {
		Unauthorized(c, "登录已过期，请重新登录")
		return nil, false
	}
	return s, true
}

func MustRole(c *gin.Context, roles ...int) (*Session, bool) {
	s, ok := MustSession(c)
	if !ok {
		return nil, false
	}
	for _, r := range roles {
		if s.Role == r {
			return s, true
		}
	}
	Fail(c, 403, "无权限访问该资源")
	return nil, false
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization,Token")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Content-Type")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
