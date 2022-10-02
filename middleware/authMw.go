package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/happyanran/walnut/common"
)

func AuthMw(svcCtx *common.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")

		if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "权限不足"})
			c.Abort()
			return
		}

		claims, err := svcCtx.ParseToken(tokenStr[7:])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "权限不足"})
			c.Abort()
			return
		}

		userID := claims.UserId

		//查库 用户状态改变

		c.Set("UserID", userID)

		c.Next()
	}
}
