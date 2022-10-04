package api

import (
	"github.com/gin-gonic/gin"
	"github.com/happyanran/walnut/model"
)

func Signin(c *gin.Context) {
	var req UserReq

	c.ShouldBindJSON(&req)

	if val := svcCtx.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	user := model.User{
		Username: req.Username,
	}

	if err := user.UserFindByName(); err != nil {
		ResponseServerErr(c, "发生错误")
		svcCtx.Log.Error(err)
		return
	}

	if user.ID == 0 || !svcCtx.PwdCheck(user.Password, req.Password) {
		ResponseClientErrDtl(c, CodeSigninErr, nil, "账号或密码错误")
		return
	}

	//token
	token, err := svcCtx.GenerateToken(user.ID)
	if err != nil {
		ResponseServerErr(c, "Token生成失败")
		svcCtx.Log.Error(err)
		return
	}

	ResponseOK(c, token, "登录成功")
}
