package tag

import (
	"github.com/gin-gonic/gin"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"strconv"
)

func GetChild(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	tags, err := model.GetSubTag(id)
	if err != nil {
		h.SendResponse(c, errno.ErrTagNoFount, err.Error())
		return
	}

	h.SendResponse(c, nil, ListResponse{TagList: tags})
}
