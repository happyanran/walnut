package main

import (
	"github.com/gin-gonic/gin"
	"github.com/happyanran/walnut/api"
	"github.com/happyanran/walnut/common"
	"github.com/happyanran/walnut/model"
	"github.com/happyanran/walnut/webdav"
)

func main() {
	cfg := common.LoadConfig()
	svcCtx := common.NewServiceContext(cfg)

	if err := model.Init(svcCtx); err != nil {
		svcCtx.Log.Fatal("Init model failed: ", err)
	}

	go webdav.WebDav(svcCtx)

	gin.SetMode(cfg.ServerConf.GinMode)
	engine := gin.Default()

	api.Router(engine, svcCtx)

	svcCtx.Log.Info("Walnut startup.")

	if svcCtx.Cfg.HttpsConf.Enable {
		svcCtx.Log.Panic(engine.RunTLS(cfg.ServerConf.Addr, svcCtx.Cfg.HttpsConf.Certfile, svcCtx.Cfg.HttpsConf.Keyfile))
	} else {
		svcCtx.Log.Panic(engine.Run(cfg.ServerConf.Addr))
	}

}
