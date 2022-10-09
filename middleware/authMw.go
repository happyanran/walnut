package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/happyanran/walnut/common"
	"github.com/happyanran/walnut/model"
)

func AuthMw(svcCtx *common.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")

		if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "权限不足"})
			c.Abort()
			return
		}

		claims, err := svcCtx.Jwtw.ParseToken(tokenStr[7:])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "权限不足"})
			c.Abort()
			return
		}

		//查库 用户状态改变
		user := &model.User{
			ID: claims.UserId,
		}

		if err := user.UserFindByID(); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "权限不足"})
			c.Abort()
			return
		}

		c.Set("UserID", user.ID)
		c.Set("Username", user.Username)

		c.Next()
	}
}
