package api

import (
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/happyanran/walnut/model"
)

type FileUploadReq struct {
	DirID   int  `form:"dirid" validate:"required,min=1"`
	IsCover bool `form:"iscover"`
}

func FileUpload(c *gin.Context) {
	var req FileUploadReq

	c.ShouldBind(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	dir := &model.Dir{
		ID: req.DirID,
	}

	if err := dir.DirFindByID(); err != nil {
		ResponseClientErrDtl(c, CodeIDNotExist, nil, "文件夹不存在")
		svcCtx.Log.Error(err)
		return
	}

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		ResponseClientErr(c, "文件上传失败")
		svcCtx.Log.Error(err)
		return
	}

	files := form.File["uploadfiles"]

	for _, file := range files {
		//写数据库
		mfile := model.File{
			DirID:   dir.ID,
			Name:    file.Filename,
			ExtType: filepath.Ext(file.Filename),
			Size:    file.Size,
		}

		//判断文件名是否冲突
		if req.IsCover {
			mfile.FileDeleteByDirName()
		} else {
			cnt, _ := mfile.FileNameCheckByDirID()
			if cnt != 0 {
				ResponseClientErrDtl(c, CodeNameExist, nil, "文件名冲突")
				return
			}
		}

		err = svcCtx.FileOp.FileUpload(
			filepath.Join(dir.Path, strconv.Itoa(dir.ID), mfile.Name),
			func(fullpath string) error {
				return c.SaveUploadedFile(file, fullpath)
			},
		)

		if err != nil {
			ResponseServerErr(c, "文件上传失败")
			svcCtx.Log.Error(err)
			return
		}

		if err := mfile.FileCreate(); err != nil {
			ResponseServerErr(c, "文件上传失败")
			svcCtx.Log.Error(err)
			return
		}
	}

	ResponseOK(c, nil, "文件上传成功")
}

type FileDelReq struct {
	ID int `json:"id" validate:"required,min=1"`
}

func FileDel(c *gin.Context) {
	var req FileDelReq

	c.ShouldBindJSON(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	file := model.File{
		ID: req.ID,
	}

	if err := file.FileDelete(); err != nil {
		ResponseServerErr(c, "文件删除失败")
		svcCtx.Log.Error(err)
		return
	}

	ResponseOK(c, nil, "文件删除成功")
}

type FileRenameReq struct {
	ID      int    `json:"id" validate:"required,min=1"`
	NewName string `json:"newname" validate:"required,min=1"`
}

func FileRename(c *gin.Context) {
	var req FileRenameReq

	c.ShouldBindJSON(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	file := model.File{
		ID: req.ID,
	}

	if err := file.FileFindByID(); err != nil {
		ResponseClientErrDtl(c, CodeIDNotExist, nil, "文件不存在")
		svcCtx.Log.Error(err)
		return
	}

	if file.Name == req.NewName {
		ResponseOK(c, nil, "文件重命名成功")
		return
	}

	oldName := file.Name
	file.Name = req.NewName
	file.ExtType = filepath.Ext(req.NewName)

	//判断文件名是否冲突
	cnt, _ := file.FileNameCheckByDirID()
	if cnt != 0 {
		ResponseClientErrDtl(c, CodeNameExist, nil, "文件名冲突")
		return
	}

	dir := model.Dir{
		ID: file.DirID,
	}

	if err := dir.DirFindByID(); err != nil {
		ResponseClientErrDtl(c, CodeIDNotExist, nil, "文件夹不存在")
		svcCtx.Log.Error(err)
		return
	}

	path := filepath.Join(dir.Path, strconv.Itoa(dir.ID))

	if err := svcCtx.FileOp.FileRename(path, oldName, path, file.Name); err != nil {
		ResponseServerErr(c, "文件重命名失败")
		svcCtx.Log.Error(err)
		return
	}

	if err := file.FileUpdate(); err != nil {
		ResponseServerErr(c, "文件重命名失败")
		svcCtx.Log.Error(err)
		return
	}

	ResponseOK(c, nil, "文件重命名成功")
}

type FileMoveReq struct {
	ID      int `json:"id" validate:"required,min=1"`
	ToDirID int `json:"todirid" validate:"required,min=1"`
}

func FileMove(c *gin.Context) {
	var req FileMoveReq

	c.ShouldBindJSON(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	file := model.File{
		ID: req.ID,
	}

	if err := file.FileFindByID(); err != nil {
		ResponseClientErrDtl(c, CodeIDNotExist, nil, "文件不存在")
		svcCtx.Log.Error(err)
		return
	}

	if file.DirID == req.ToDirID {
		ResponseOK(c, nil, "文件移动成功")
		return
	}

	FromDir := model.Dir{
		ID: file.DirID,
	}

	if err := FromDir.DirFindByID(); err != nil {
		ResponseClientErrDtl(c, CodeIDNotExist, nil, "文件夹不存在")
		svcCtx.Log.Error(err)
		return
	}

	fromPath := filepath.Join(FromDir.Path, strconv.Itoa(FromDir.ID))

	file.DirID = req.ToDirID

	//判断文件名是否冲突
	cnt, _ := file.FileNameCheckByDirID()
	if cnt != 0 {
		ResponseClientErrDtl(c, CodeNameExist, nil, "文件名冲突")
		return
	}

	toDir := model.Dir{
		ID: file.DirID,
	}

	if err := toDir.DirFindByID(); err != nil {
		ResponseClientErrDtl(c, CodeIDNotExist, nil, "文件夹不存在")
		svcCtx.Log.Error(err)
		return
	}

	//文件移动
	toPath := filepath.Join(toDir.Path, strconv.Itoa(toDir.ID))

	if err := svcCtx.FileOp.FileRename(fromPath, file.Name, toPath, file.Name); err != nil {
		ResponseServerErr(c, "文件移动失败")
		svcCtx.Log.Error(err)
		return
	}

	if err := file.FileUpdate(); err != nil {
		ResponseServerErr(c, "文件移动失败")
		svcCtx.Log.Error(err)
		return
	}

	ResponseOK(c, nil, "文件移动成功")
}

type FileGetAllReq struct {
	DirID int `json:"dirid" validate:"required,min=1"`
}

func FileGetAll(c *gin.Context) {
	var req FileGetAllReq

	c.ShouldBindJSON(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}
}
