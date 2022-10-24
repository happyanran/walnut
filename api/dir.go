package api

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/happyanran/walnut/model"
)

type DirAddReq struct {
	PID  int    `json:"PID" validate:"required,min=1"`
	Name string `json:"name" validate:"required,min=1"`
	Note string `json:"note"`
}

func DirAdd(c *gin.Context) {
	var req DirAddReq

	c.ShouldBindJSON(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	PDir := &model.Dir{
		ID: req.PID,
	}

	if err := PDir.DirFindByID(); err != nil {
		ResponseClientErrDtl(c, CodeIDNotExist, nil, "文件夹不存在")
		return
	}

	dir := &model.Dir{
		PID:  req.PID,
		Path: PDir.Path + strconv.Itoa(PDir.ID) + ",",
		Name: req.Name,
		Note: req.Note,
	}

	if err := dir.DirCreate(); err != nil {
		ResponseClientErrDtl(c, CodeNameExist, nil, "文件夹名冲突")
		return
	}

	ResponseOK(c, nil, "文件夹创建成功")
}

type DirDelReq struct {
	ID int `json:"ID" validate:"required,min=2"`
}

func DirDel(c *gin.Context) {
	var req DirDelReq

	c.ShouldBindJSON(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	dir := &model.Dir{
		ID: req.ID,
	}

	if err := dir.DirDelete(); err != nil {
		ResponseServerErr(c)
		return
	}

	ResponseOK(c, nil, "文件夹删除成功")
}

type DirRenameReq struct {
	ID      int    `json:"ID" validate:"required,min=2"`
	NewName string `json:"newName" validate:"required,min=1"`
}

func DirRename(c *gin.Context) {
	var req DirRenameReq

	c.ShouldBindJSON(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	dir := &model.Dir{
		ID: req.ID,
	}

	if err := dir.DirFindByID(); err != nil {
		ResponseClientErrDtl(c, CodeIDNotExist, nil, "文件夹不存在")
		return
	}

	dir.Name = req.NewName

	if err := dir.DirUpdate(); err != nil {
		ResponseClientErrDtl(c, CodeNameExist, nil, "文件夹名冲突")
		return
	}

	ResponseOK(c, nil, "文件夹重命名成功")
}

type DirMoveReq struct {
	ID   int `json:"ID" validate:"required,min=2"`
	ToID int `json:"toID" validate:"required,min=1,nefield=ID"`
}

func DirMove(c *gin.Context) {
	var req DirMoveReq

	c.ShouldBindJSON(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	dir := &model.Dir{
		ID: req.ID,
	}

	if err := dir.DirFindByID(); err != nil {
		ResponseClientErrDtl(c, CodeIDNotExist, nil, "文件夹不存在")
		return
	}

	if dir.PID == req.ToID {
		ResponseOK(c, nil, "文件夹移动成功")
		return
	}

	toDir := &model.Dir{
		ID: req.ToID,
	}

	if err := toDir.DirFindByID(); err != nil {
		ResponseClientErrDtl(c, CodeIDNotExist, nil, "文件夹不存在")
		return
	}

	//判断是否是自己的子目录
	if strings.Contains(toDir.Path, ","+strconv.Itoa(dir.ID)+",") {
		ResponseClientErrDtl(c, CodeDirMoveToChild, nil, "文件夹不能移动到它的子目录中")
		return
	}

	oldPath := dir.Path + strconv.Itoa(dir.ID) + ","
	//更新自己
	dir.PID = toDir.ID
	dir.Path = toDir.Path + strconv.Itoa(toDir.ID) + ","

	if err := dir.DirUpdate(); err != nil {
		ResponseClientErrDtl(c, CodeNameExist, nil, "文件夹名冲突")
		return
	}

	//更新子目录path
	if err := dir.DirMoveChilds(oldPath); err != nil {
		ResponseServerErr(c)
		return
	}

	ResponseOK(c, nil, "文件夹移动成功")
}

type DirGetReq struct {
	ID int `form:"ID" validate:"required,min=1"`
}

func DirGet(c *gin.Context) {
	var req DirGetReq

	c.ShouldBind(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	dir := &model.Dir{
		ID: req.ID,
	}

	if err := dir.DirFindByID(); err != nil {
		ResponseClientErrDtl(c, CodeIDNotExist, nil, "文件夹不存在")
		return
	}

	ResponseOK(c, dir, "获取文件夹成功")
}

func DirGetChilds(c *gin.Context) {
	var req DirGetReq

	c.ShouldBind(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	dir := &model.Dir{
		ID: req.ID,
	}

	var childDirs []model.Dir

	if err := dir.FindChildDirs(&childDirs); err != nil {
		ResponseServerErr(c)
		return
	}

	ResponseOK(c, childDirs, "获取子文件夹成功")
}
