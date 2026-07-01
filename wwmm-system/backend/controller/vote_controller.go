package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"wwmm/service"
	"wwmm/utils"
)

func CastVote(c *gin.Context) {
	s, ok := utils.MustSession(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if id < 1 {
		utils.Fail(c, 400, "无效ID")
		return
	}
	txHash, err := service.CastVote(s.UserID, id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}
	utils.Success(c, gin.H{
		"txHash":   txHash,
		"photoId":  id,
		"voterId":  s.UserID,
		"onChain":  true,
		"verified": true,
	})
}

func hasVoted(uid, pid int) (bool, error) {
	if uid <= 0 {
		return false, nil
	}
	return hasVotedDB(uid, pid)
}
