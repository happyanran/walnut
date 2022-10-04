package api

import (
	"github.com/gin-gonic/gin"
	"github.com/happyanran/walnut/common"
	"github.com/happyanran/walnut/middleware"
)

var svcCtx *common.ServiceContext

func Router(e *gin.Engine, s *common.ServiceContext) {
	svcCtx = s

	api := e.Group("/api")
	{
		api.POST("/signin", Signin)

		user := api.Group("/user", middleware.AuthMw(svcCtx))
		{
			user.GET("/all", UserGetAll)
			user.GET("/allcnt", UserGetAllCnt)
			user.POST("/add", UserAdd)
			user.POST("/del", UserDel)
			user.POST("/change", UserChange)
		}
	}
}
