package profile

import (
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	uids := []uint64{id}
	profiles, err := model.GetProfileWithGroupAndTag(uids)
	if err != nil {
		h.SendResponse(c, errno.ErrUserNotFound, err.Error())
		return
	}

	if len(profiles) > 0 {
		h.SendResponse(c, nil, CreateResponse{Profile: profiles[0]})
		return
	}

	h.SendResponse(c, nil, CreateResponse{})
}
