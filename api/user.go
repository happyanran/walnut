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

type UserGetAllReq struct {
	Page int `json:"page" validate:"required,min=1"`
	Size int `json:"size" validate:"required,min=1"`
}

func UserGetAll(c *gin.Context) {
	if c.GetInt("UserID") != 1 {
		ResponseClientErrDtl(c, CodeNotAdmin, nil, "无权访问")
		return
	}

	var req UserGetAllReq

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

type UserAddReq struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=5,max=60"`
}

func UserAdd(c *gin.Context) {
	if c.GetInt("UserID") != 1 {
		ResponseClientErrDtl(c, CodeNotAdmin, nil, "无权访问")
		return
	}

	var req UserAddReq

	c.ShouldBindJSON(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	user := &model.User{
		Username: req.Username,
		Password: svcCtx.Utilw.PwdEnrypt(req.Password),
	}

	if err := user.UserCreate(); err != nil {
		ResponseServerErr(c, "用户创建失败")
		svcCtx.Log.Error(err)
		return
	}

	ResponseOK(c, nil, "用户创建成功")
}

type UserDelReq struct {
	UserID int `json:"userid" validate:"required,min=2"`
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

	user := &model.User{
		ID: req.UserID,
	}

	if err := user.UserDelete(); err != nil {
		ResponseServerErr(c, "用户删除失败")
		svcCtx.Log.Error(err)
		return
	}

	ResponseOK(c, nil, "用户删除成功")
}

type UserChangeReq struct {
	UserID   int    `json:"userid" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=5,max=60"`
}

func UserChange(c *gin.Context) {
	var req UserChangeReq

	c.ShouldBindJSON(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	if c.GetInt("UserID") != 1 && c.GetInt("UserID") != req.UserID {
		ResponseClientErrDtl(c, CodeNotAdmin, nil, "无权访问")
		return
	}

	user := &model.User{
		ID:       req.UserID,
		Password: svcCtx.Utilw.PwdEnrypt(req.Password),
	}

	if err := user.UserUpdate(); err != nil {
		ResponseServerErr(c, "密码更新失败")
		svcCtx.Log.Error(err)
		return
	}

	ResponseOK(c, nil, "密码更新成功")
}
