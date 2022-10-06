package api

import "github.com/gin-gonic/gin"

type FileGetAllReq struct {
	DirID int `json:"dirid" validate:"required,min=0"`
}

func FileGetAll(c *gin.Context) {
	var req FileGetAllReq

	c.ShouldBindJSON(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}
}
