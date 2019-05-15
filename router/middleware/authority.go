package middleware

import (
	"fmt"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"hrgdrc/handler"
	"hrgdrc/pkg/errno"
	"hrgdrc/pkg/token"
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

		//if strings.Contains(obj, "/sd/") {
		//	c.Next()
		//	return
		//}
		//
		//if strings.Contains(obj, "/login") {
		//	c.Next()
		//	return
		//}
		//
		//if strings.Contains(obj, "/static") {
		//	c.Next()
		//	return
		//}
		//
		//if strings.Contains(obj, "/favicon.ico") {
		//	c.Next()
		//	return
		//}
		//
		//if strings.Contains(obj, "/download") {
		//	c.Next()
		//	return
		//}
		//if strings.Contains(obj, "/upload") {
		//	c.Next()
		//	return
		//}
		if e.Enforce(sub, obj, act) == true {
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
