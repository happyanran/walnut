package api

import (
	"github.com/gin-gonic/gin"
	"github.com/happyanran/walnut/model"
)

type SigninReq struct {
	UserName string `json:"userName" validate:"required,min=1,max=20"`
	Password string `json:"password" validate:"required,min=1,max=60"`
}

type SigninResp struct {
	NickName string `json:"nickName"`
	Token    string `json:"token"`
}

func Signin(c *gin.Context) {
	var req SigninReq

	c.ShouldBindJSON(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	user := model.User{
		UserName: req.UserName,
	}

	if err := user.UserFindByName(); err != nil {
		ResponseServerErr(c)
		return
	}

	if user.ID == 0 || !svcCtx.Utilw.PwdCheck(user.Password, req.Password) {
		ResponseClientErrDtl(c, CodeSigninErr, nil, "账号或密码错误")
		return
	}

	//token
	token, err := svcCtx.Jwtw.GenerateToken(user.ID)
	if err != nil {
		ResponseServerErr(c)
		return
	}

	ResponseOK(c, SigninResp{NickName: user.NickName, Token: token}, "登录成功")
}
