package middleware

import (
	"fmt"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"hr-server/handler"
	"hr-server/pkg/errno"
	"hr-server/pkg/token"
	"net/http"
)

func Authority(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, _ := token.ParseRequest(c)

		sub := body.Role          // the user that wants to access a resource.
		obj := c.Request.URL.Path // the resource that is going to be accessed.
		act := c.Request.Method   // the operation that the user performs on the resource.
		fmt.Println("sub,obj,act", sub, obj, act)
		c.Set("userid", body.ID)

		result, _ := e.Enforce(sub, obj, act)
		if result == true {
			// permit alice to read data1
			c.Next()
		} else {
			// deny the request, show an error
			c.JSON(http.StatusUnauthorized, handler.Response{
				Code:    errno.StatusUnauthorized.Code,
				Message: "deny\n",
				Data:    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
