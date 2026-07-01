package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func Sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func Sha256HexString(s string) string {
	return Sha256Hex([]byte(s))
}

func HashPassword(password, salt string) string {
	h := sha256.New()
	h.Write([]byte(salt + ":" + password))
	return hex.EncodeToString(h.Sum(nil))
}

func GenSalt() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func GenTxHash(txType int, sender string, payload string, ts int64) string {
	raw := fmt.Sprintf("%d|%s|%s|%d", txType, sender, payload, ts)
	return Sha256HexString(raw)
}

func NowTs() int64 {
	return time.Now().Unix()
}
