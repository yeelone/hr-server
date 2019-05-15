package tag

import (
	"fmt"
	"github.com/gin-gonic/gin"
	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/errno"
	"strconv"
)

func GetChild(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	fmt.Println("id", id )
	tags, err := model.GetSubTag(id)
	if err != nil {
		h.SendResponse(c, errno.ErrTagNoFount, err.Error())
		return
	}

	h.SendResponse(c, nil, ListResponse{TagList: tags})
}
