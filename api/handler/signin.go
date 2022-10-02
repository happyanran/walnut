package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/happyanran/walnut/model"
)

type SigninReq struct {
	Username string `json:"username" validate:"required,min=1,max=20"`
	Password string `json:"password" validate:"required,min=1,max=60"`
}

func Signin(c *gin.Context) {
	var req SigninReq

	c.ShouldBindJSON(&req)

	if val := svcCtx.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	user := model.User{
		Username: req.Username,
		Password: svcCtx.PwdEnrypt(req.Password),
	}

	if err := user.UserFindSignin(); err != nil || user.ID == 0 {
		ResponseClientErrDtl(c, CodeSigninErr, nil, "账号或密码错误")
		return
	}

	//token
	token, err := svcCtx.GenerateToken(user.ID)
	if err != nil {
		ResponseServerErr(c, "Token生成失败")
		return
	}

	ResponseOK(c, token, "登录成功")
}
