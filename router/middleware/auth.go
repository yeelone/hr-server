package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"hrgdrc/handler"
	"hrgdrc/pkg/errno"
	"hrgdrc/pkg/token"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Parse the json web token.
		ctx, err := token.ParseRequest(c)
		if err != nil {
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}
		log.Infof("user have passed authenticated and trying to access resource . ID: %d | username: %s | resource: %s  ", ctx.ID, ctx.Username, c.Request.URL)
		c.Next()
	}
}
