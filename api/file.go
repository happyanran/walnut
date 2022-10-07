package api

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/happyanran/walnut/model"
)

type FileUploadReq struct {
	DirID int `form:"dirid" validate:"required,min=2"`
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
		ResponseClientErrDtl(c, CodeDirNotExist, nil, "文件夹不存在")
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
		fmt.Println(file.Filename)
		fmt.Println(file.Size)
		fmt.Println(file.Header["Content-Type"])

		//写数据库 file名字改成主键

		err = svcCtx.FileOp.FileUpload(
			filepath.Join(dir.Path, strconv.Itoa(dir.ID), file.Filename),
			func(fullpath string) error {
				return c.SaveUploadedFile(file, fullpath)
			},
		)

		if err != nil {
			ResponseServerErr(c, "文件上传失败")
			svcCtx.Log.Error(err)
			return
		}
	}

	ResponseOK(c, nil, "文件上传成功")
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
