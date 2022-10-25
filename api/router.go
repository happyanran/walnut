package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/happyanran/walnut/common"
	"github.com/happyanran/walnut/middleware"
)

var svcCtx *common.ServiceContext

func Router(e *gin.Engine, s *common.ServiceContext) {
	svcCtx = s

	e.MaxMultipartMemory = 1 << 30 // 1G

	e.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{"content-type", "authorization"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	api := e.Group("/api")
	{
		api.POST("/signin", Signin)

		user := api.Group("/user", middleware.AuthMw(svcCtx))
		{
			user.POST("/add", UserAdd)       //todo
			user.POST("/del", UserDel)       //todo
			user.POST("/change", UserChange) //todo
			user.GET("/getall", UserGetAll)  //todo
		}

		file := api.Group("/file", middleware.AuthMw(svcCtx))
		{
			file.POST("/diradd", DirAdd)
			file.POST("/dirdel", DirDel)
			file.POST("/dirrename", DirRename)
			file.POST("/dirmove", DirMove)
			file.GET("/dirget", DirGet)
			file.GET("/dirgetchilds", DirGetChilds)

			file.POST("/fileupload", FileUpload)
			file.POST("/filedel", FileDel)
			file.POST("/filerename", FileRename)
			file.POST("/filemove", FileMove)
			file.GET("/filegetbydir", FileGetByDir)
			//file.GET("/filegetbytype", FileGetByType)

			file.StaticFS("/staticfs", gin.Dir(svcCtx.Cfg.ServerConf.Data, false))
		}
	}
}
