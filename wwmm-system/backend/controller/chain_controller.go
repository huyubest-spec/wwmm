package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"wwmm/dao"
	"wwmm/service"
	"wwmm/utils"
)

func ChainState(c *gin.Context) {
	st, _ := dao.GetChainState()
	totalBlocks, _ := dao.CountBlocks()
	totalTxs, _ := dao.CountTxs()
	txCertify, _ := dao.CountTxsByType(1)
	txVote, _ := dao.CountTxsByType(2)
	utils.Success(c, gin.H{
		"latestIndex":     st.LatestIndex,
		"latestHash":      st.LatestHash,
		"totalBlocks":     totalBlocks,
		"totalTxs":        totalTxs,
		"txCertifyCount":  txCertify,
		"txVoteCount":     txVote,
		"updatedAt":       st.UpdatedAt,
	})
}

func ListBlocks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 10
	}
	list, err := dao.ListBlocks(size, (page-1)*size)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}
	total, _ := dao.CountBlocks()
	utils.Success(c, gin.H{"list": list, "total": total, "page": page, "size": size})
}

func BlockDetail(c *gin.Context) {
	idx, _ := strconv.Atoi(c.Param("index"))
	if idx < 1 {
		utils.Fail(c, 400, "无效index")
		return
	}
	b, err := dao.GetBlockByIndex(idx)
	if err != nil {
		utils.Fail(c, 404, "区块不存在")
		return
	}
	txs, _ := dao.ListTxsByBlock(b.Index)
	for i := range txs {
		txs[i].Payload = prettyJSON(txs[i].Payload)
	}
	utils.Success(c, gin.H{"block": b, "txs": txs})
}

func TxDetail(c *gin.Context) {
	hash := c.Param("hash")
	t, err := service.GetTxByHash(hash)
	if err != nil {
		utils.Fail(c, 404, "交易不存在")
		return
	}
	t.Payload = prettyJSON(t.Payload)
	utils.Success(c, t)
}

func VerifyImage(c *gin.Context) {
	hash := c.Param("hash")
	if len(hash) != 64 {
		utils.Fail(c, 400, "无效哈希")
		return
	}
	p, err := dao.GetPhotoByHash(hash)
	if err != nil {
		utils.Fail(c, 404, "该哈希在链上未找到对应作品")
		return
	}
	t, _ := dao.GetTxByHash(p.ChainTxHash)
	utils.Success(c, gin.H{
		"photo":      p,
		"chainTx":    t,
		"verified":   true,
		"message":    "哈希校验通过，作品存证信息真实有效",
	})
}

func ListAllTxs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "15"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 15
	}
	list, err := dao.ListTxs(size, (page-1)*size)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}
	for i := range list {
		list[i].Payload = prettyJSON(list[i].Payload)
	}
	total, _ := dao.CountTxs()
	utils.Success(c, gin.H{"list": list, "total": total, "page": page, "size": size})
}
