package profile

import (
	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	profile, err := model.GetProfile(id)
	if err != nil {
		h.SendResponse(c, errno.ErrUserNotFound, err.Error())
		return
	}
	h.SendResponse(c, nil, CreateResponse{Profile:profile})
}
