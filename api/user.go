package api

import (
	"github.com/gin-gonic/gin"
	"github.com/happyanran/walnut/model"
)

func UserGetAllCnt(c *gin.Context) {
	if c.GetInt("UserID") != 1 {
		ResponseClientErrDtl(c, CodeNotAdmin, nil, "无权访问")
		return
	}

	var u model.User
	var cnt int64

	if err := u.UserCount(&cnt); err != nil {
		ResponseServerErr(c, "发生错误")
		svcCtx.Log.Error(err)
		return
	}

	ResponseOK(c, cnt, "成功")
}

func UserGetAll(c *gin.Context) {
	if c.GetInt("UserID") != 1 {
		ResponseClientErrDtl(c, CodeNotAdmin, nil, "无权访问")
		return
	}

	var req PageReq

	c.ShouldBindJSON(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	var u model.User
	var users []model.User

	if err := u.UserGetAll(&users, req.Size, (req.Page-1)*req.Size); err != nil {
		ResponseServerErr(c, "发生错误")
		svcCtx.Log.Error(err)
		return
	}

	ResponseOK(c, users, "成功")
}

func UserAdd(c *gin.Context) {
	if c.GetInt("UserID") != 1 {
		ResponseClientErrDtl(c, CodeNotAdmin, nil, "无权访问")
		return
	}

	var req UserReq

	c.ShouldBindJSON(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	var u = &model.User{
		Username: req.Username,
		Password: svcCtx.Utilw.PwdEnrypt(req.Password),
	}

	if err := u.UserCreate(); err != nil {
		ResponseServerErr(c, "用户创建失败")
		svcCtx.Log.Error(err)
		return
	}

	ResponseOK(c, nil, "用户创建成功")
}

func UserDel(c *gin.Context) {
	if c.GetInt("UserID") != 1 {
		ResponseClientErrDtl(c, CodeNotAdmin, nil, "无权访问")
		return
	}

	var req UserDelReq

	c.ShouldBindJSON(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	var u = &model.User{
		ID: req.UserID,
	}

	if err := u.UserDelete(); err != nil {
		ResponseServerErr(c, "用户删除失败")
		svcCtx.Log.Error(err)
		return
	}

	ResponseOK(c, nil, "用户删除成功")
}

func UserChange(c *gin.Context) {
	var req UserChgReq

	c.ShouldBindJSON(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	if c.GetInt("UserID") != 1 && c.GetInt("UserID") != req.UserID {
		ResponseClientErrDtl(c, CodeNotAdmin, nil, "无权访问")
		return
	}

	var u = &model.User{
		ID:       req.UserID,
		Password: svcCtx.Utilw.PwdEnrypt(req.Password),
	}

	if err := u.UserUpdate(); err != nil {
		ResponseServerErr(c, "密码更新失败")
		svcCtx.Log.Error(err)
		return
	}

	ResponseOK(c, nil, "密码更新成功")
}
