package main

import (
	"github.com/gin-gonic/gin"
	"github.com/happyanran/walnut/api"
	"github.com/happyanran/walnut/common"
	"github.com/happyanran/walnut/model"
)

func main() {
	cfg := common.LoadConfig()
	svcCtx := common.NewServiceContext(cfg)

	if err := model.Init(svcCtx); err != nil {
		svcCtx.Log.Fatal("Init model failed: ", err)
	}

	gin.SetMode(cfg.ServerConf.GinMode)
	engine := gin.Default()
	api.Router(engine, svcCtx)

	svcCtx.Log.Panic(engine.Run(cfg.ServerConf.Addr))
}
