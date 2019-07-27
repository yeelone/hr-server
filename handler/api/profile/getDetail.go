package profile

import (
	"fmt"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	uids := []uint64{id}
	profiles, err := model.GetProfileWithGroupAndTag(uids)
	fmt.Println(util.PrettyJson(profiles))
	if err != nil {
		h.SendResponse(c, errno.ErrUserNotFound, err.Error())
		return
	}
	h.SendResponse(c, nil, CreateResponse{Profile: profiles[0]})
}
