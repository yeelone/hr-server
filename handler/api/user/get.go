package user

import (
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"

	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	username := c.Param("username")
	// Get the user by the `username` from the database.
	user, err := model.GetUserByName(username)
	if err != nil {
		h.SendResponse(c, errno.ErrUserNotFound, err.Error())
		return
	}

	h.SendResponse(c, nil, user)
}
