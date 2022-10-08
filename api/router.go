package api

import (
	"github.com/gin-gonic/gin"
	"github.com/happyanran/walnut/common"
	"github.com/happyanran/walnut/middleware"
)

var svcCtx *common.ServiceContext

func Router(e *gin.Engine, s *common.ServiceContext) {
	svcCtx = s

	e.MaxMultipartMemory = 100 << 20 // 100 MiB

	api := e.Group("/api")
	{
		api.POST("/auth/signin", Signin)

		user := api.Group("/user", middleware.AuthMw(svcCtx))
		{
			user.GET("/getall", UserGetAll)
			user.GET("/getallcnt", UserGetAllCnt)
			user.POST("/add", UserAdd)
			user.POST("/del", UserDel)
			user.POST("/change", UserChange)
		}

		file := api.Group("/file", middleware.AuthMw(svcCtx))
		{
			file.POST("/diradd", DirAdd)
			file.POST("/dirdel", DirDel)
			file.POST("/dirrename", DirRename)
			file.POST("/dirmove", DirMove)

			file.POST("/upload", FileUpload)
			file.POST("/filedel", FileDel)
			file.POST("/filerename", FileRename)
			file.POST("/filemove", FileMove)

			file.GET("/get", FileGet)
		}
	}
}
