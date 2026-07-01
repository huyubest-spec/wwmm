package controller

import (
	"github.com/gin-gonic/gin"
	"wwmm/utils"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/health", Health)

	api := r.Group("/api")
	{
		// 用户
		api.POST("/user/register", Register)
		api.POST("/user/login", Login)
		api.POST("/user/logout", Logout)
		api.GET("/user/me", Me)

		// 作品
		api.GET("/photo/list", ListPhotos)
		api.GET("/photo/pending", ListPendingPhotos)
		api.GET("/photo/mine", ListMyPhotos)
		api.GET("/photo/:id", PhotoDetail)
		api.POST("/photo/upload", UploadPhoto)
		api.POST("/photo/:id/audit", AuditPhoto)

		// 投票
		api.POST("/photo/:id/vote", CastVote)

		// 区块链
		api.GET("/chain/state", ChainState)
		api.GET("/chain/blocks", ListBlocks)
		api.GET("/chain/block/:index", BlockDetail)
		api.GET("/chain/tx/:hash", TxDetail)
		api.GET("/chain/verify/:hash", VerifyImage)
		api.GET("/chain/txs", ListAllTxs)
	}

	_ = utils.Cors
}
