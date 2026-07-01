package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"wwmm/config"
	"wwmm/controller"
	"wwmm/service"
	"wwmm/utils"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	utils.InitDB()
	if err := seedDefaultUsers(); err != nil {
		log.Printf("种子用户写入失败: %v", err)
	}
	if err := service.InitGenesis(); err != nil {
		log.Fatalf("创世区块初始化失败: %v", err)
	}

	_ = os.MkdirAll(config.App.UploadDir+"/photos", 0755)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(utils.Cors())
	r.Use(func(c *gin.Context) {
		t := time.Now()
		c.Next()
		log.Printf("[HTTP] %s %s %d %v", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(t))
	})

	r.Static("/static", "./uploads")
	controller.RegisterRoutes(r)

	addr := ":" + config.App.ServerPort
	log.Printf("[OK] WWMM 后端启动于 %s", addr)
	if err := r.Run(addr); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
