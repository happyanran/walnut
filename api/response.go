package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeOK        = iota //http 200
	CodeClientErr        //http 400
	CodeServerErr        //http 500
	CodeReqValErr
	CodeSigninErr
	CodeNotAdmin
	CodeDataExist
	CodeIDNotExist
	CodeNameExist
	CodeDirMoveToChild
)

func Response(c *gin.Context, httpStatus int, code int, data interface{}, msg string) {
	c.JSON(httpStatus, gin.H{"code": code, "data": data, "msg": msg})
}

func ResponseOK(c *gin.Context, data interface{}, msg string) {
	Response(c, http.StatusOK, CodeOK, data, msg)
}

func ResponseClientErr(c *gin.Context, msg string) {
	Response(c, http.StatusBadRequest, CodeClientErr, nil, msg)
}

func ResponseClientErrDtl(c *gin.Context, code int, data interface{}, msg string) {
	Response(c, http.StatusBadRequest, code, data, msg)
}

func ResponseServerErr(c *gin.Context, msg string) {
	Response(c, http.StatusInternalServerError, CodeServerErr, nil, msg)
}

func ResponseServerErrDtl(c *gin.Context, code int, data interface{}, msg string) {
	Response(c, http.StatusInternalServerError, code, data, msg)
}
