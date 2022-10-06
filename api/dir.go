package api

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/happyanran/walnut/model"
)

type DirAddReq struct {
	PID  int    `json:"pid" validate:"required,min=1"`
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

	pdir := &model.Dir{
		ID: req.PID,
	}

	if err := pdir.DirFindByID(); err != nil {
		ResponseClientErrDtl(c, CodeDirNotExist, nil, "文件夹不存在")
		svcCtx.Log.Error(err)
		return
	}

	path := filepath.Join(pdir.Path, strconv.Itoa(pdir.ID))

	dir := &model.Dir{
		PID:  req.PID,
		Path: path,
		Name: req.Name,
		Note: req.Note,
	}

	if err := dir.DirCreate(); err != nil {
		ResponseServerErr(c, "文件夹创建失败")
		svcCtx.Log.Error(err)
		return
	}

	if err := svcCtx.FileOp.DirCreate(path, strconv.Itoa(dir.ID)); err != nil {
		ResponseServerErr(c, "文件夹创建失败")
		svcCtx.Log.Error(err)
		return
	}

	ResponseOK(c, nil, "文件夹创建成功")
}

type DirDelReq struct {
	ID int `json:"id" validate:"required,min=2"`
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

	if err := dir.DirFindByID(); err != nil {
		ResponseClientErrDtl(c, CodeDirNotExist, nil, "文件夹不存在")
		svcCtx.Log.Error(err)
		return
	}

	if err := dir.DirDelNested(); err != nil {
		ResponseServerErr(c, "文件夹删除失败")
		svcCtx.Log.Error(err)
		return
	}

	//文件物理删除

	ResponseOK(c, nil, "文件夹删除成功")
}

type DirRenameReq struct {
	ID      int    `json:"id" validate:"required,min=2"`
	NewName string `json:"newname" validate:"required,min=1"`
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
		ResponseClientErrDtl(c, CodeDirNotExist, nil, "文件夹不存在")
		svcCtx.Log.Error(err)
		return
	}

	dir.Name = req.NewName

	if err := dir.DirUpdate(); err != nil {
		ResponseServerErr(c, "文件夹重命名失败")
		svcCtx.Log.Error(err)
		return
	}

	ResponseOK(c, nil, "文件夹重命名成功")
}

type DirMoveReq struct {
	ID   int `json:"id" validate:"required,min=2"`
	ToID int `json:"toid" validate:"required,min=1,nefield=ID"`
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
		ResponseClientErrDtl(c, CodeDirNotExist, nil, "文件夹不存在")
		svcCtx.Log.Error(err)
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
		ResponseClientErrDtl(c, CodeDirNotExist, nil, "文件夹不存在")
		svcCtx.Log.Error(err)
		return
	}

	//判断是否是自己的子目录 a/b/c -> a/b/c/
	childDirPath := strings.TrimRight(filepath.Join(dir.Path, strconv.Itoa(dir.ID), "*"), "*")
	if strings.Contains(filepath.Join(toDir.Path, "*"), childDirPath) {
		ResponseClientErrDtl(c, CodeDirMoveToChild, nil, "文件夹不能移动到它的子目录中")
		return
	}

	//文件夹移动
	newPath := filepath.Join(toDir.Path, strconv.Itoa(toDir.ID))

	if err := svcCtx.FileOp.FileRename(dir.Path, strconv.Itoa(dir.ID), newPath, strconv.Itoa(dir.ID)); err != nil {
		ResponseServerErr(c, "文件夹移动失败")
		svcCtx.Log.Error(err)
		return
	}

	//更新子目录path
	if err := dir.DirUpdateChild(filepath.Join(newPath, strconv.Itoa(dir.ID))); err != nil {
		ResponseServerErr(c, "文件夹移动失败")
		svcCtx.Log.Error(err)
		return
	}

	//更新当前目录
	dir.PID = req.ToID
	dir.Path = newPath

	if err := dir.DirUpdate(); err != nil {
		ResponseServerErr(c, "文件夹移动失败")
		svcCtx.Log.Error(err)
		return
	}

	ResponseOK(c, nil, "文件夹移动成功")
}