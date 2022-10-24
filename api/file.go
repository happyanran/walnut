package api

import (
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/happyanran/walnut/model"
)

type FileUploadReq struct {
	DirID   int  `form:"dirID" validate:"required,min=1"`
	IsCover bool `form:"isCover"`
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
		return
	}

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		ResponseClientErr(c, "文件上传失败")
		svcCtx.Log.Error(err)
		return
	}

	files := form.File["uploadFiles"]

	//var wg sync.WaitGroup

	for _, f := range files {
		file := model.File{
			DirID:   req.DirID,
			Name:    f.Filename,
			ExtType: strings.ToUpper(strings.Replace(filepath.Ext(f.Filename), ".", "", -1)),
		}

		//判断文件名是否冲突
		file.FileFindByName()

		if file.ID > 0 {
			if !req.IsCover {
				ResponseClientErrDtl(c, CodeNameExist, nil, "文件名冲突")
				return
			}
		} else {
			nanoName := svcCtx.Utilw.GenNanoName()

			file.OriginalName = "Ori" + "-" + nanoName + "." + file.ExtType
			file.SmallImgName = "Sma" + "-" + nanoName + "." + file.ExtType
			file.LargeImgName = "Lar" + "-" + nanoName + "." + file.ExtType
		}

		file.OriginalSize = f.Size

		err = svcCtx.FileOp.FileUpload(
			file.ExtType,
			file.OriginalName,
			func(fullpath string) error {
				return c.SaveUploadedFile(f, fullpath)
			},
		)

		if err != nil {
			ResponseServerErr(c)
			svcCtx.Log.Error(err)
			return
		}

		if file.ExtType == "JPG" || file.ExtType == "JPEG" {
			//wg.Add(2)
			go svcCtx.FileOp.ImgResize(file.ExtType, file.OriginalName, file.SmallImgName, 200, 0)
			go svcCtx.FileOp.ImgResize(file.ExtType, file.OriginalName, file.LargeImgName, 800, 0)
		}

		if file.ID > 0 {
			if err := file.FileUpdate(); err != nil {
				ResponseServerErr(c)
				return
			}
		} else {
			if err := file.FileCreate(); err != nil {
				ResponseServerErr(c)
				return
			}
		}
	}

	ResponseOK(c, nil, "文件上传成功")

	//wg.Wait()
}

type FileDelReq struct {
	ID int `json:"ID" validate:"required,min=1"`
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
		ResponseServerErr(c)
		return
	}

	ResponseOK(c, nil, "文件删除成功")
}

type FileRenameReq struct {
	ID      int    `json:"ID" validate:"required,min=1"`
	NewName string `json:"newName" validate:"required,min=1"`
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
		return
	}

	if file.ExtType != strings.ToUpper(strings.Replace(filepath.Ext(req.NewName), ".", "", -1)) {
		ResponseClientErrDtl(c, CodeFileExtErr, nil, "不允许修改文件类型")
		return
	}

	file.Name = req.NewName

	if err := file.FileUpdate(); err != nil {
		ResponseClientErrDtl(c, CodeNameExist, nil, "文件名冲突")
		return
	}

	ResponseOK(c, nil, "文件重命名成功")
}

type FileMoveReq struct {
	ID      int `json:"ID" validate:"required,min=1"`
	ToDirID int `json:"toDirID" validate:"required,min=1"`
}

func FileMove(c *gin.Context) {
	var req FileMoveReq

	c.ShouldBindJSON(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	toDir := model.Dir{
		ID: req.ToDirID,
	}

	if err := toDir.DirFindByID(); err != nil {
		ResponseClientErrDtl(c, CodeIDNotExist, nil, "文件夹不存在")
		return
	}

	file := model.File{
		ID: req.ID,
	}

	if err := file.FileFindByID(); err != nil {
		ResponseClientErrDtl(c, CodeIDNotExist, nil, "文件不存在")
		return
	}

	file.DirID = req.ToDirID

	if err := file.FileUpdate(); err != nil {
		ResponseClientErrDtl(c, CodeIDNotExist, nil, "文件名冲突")
		return
	}

	ResponseOK(c, nil, "文件移动成功")
}

type FileGetByDirReq struct {
	DirID int `form:"dirID" validate:"required,min=1"`
}

func FileGetByDir(c *gin.Context) {
	var req FileGetByDirReq

	c.ShouldBind(&req)

	if val := svcCtx.ZhVal.Struct(req); val != nil {
		ResponseClientErrDtl(c, CodeReqValErr, val, "请求参数错误")
		return
	}

	file := model.File{
		DirID: req.DirID,
	}

	var files []model.File

	if err := file.FileFindByDirID(&files); err != nil {
		return
	}

	ResponseOK(c, files, "获取文件成功")
}
