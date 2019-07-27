package tag

import (
	"fmt"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RelatedUserList(c *gin.Context) {
	tagID, _ := strconv.Atoi(c.Param("id"))
	tag, err := model.GetTag(uint64(tagID), true)
	if err != nil {
		fmt.Println("errors", err)
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}

	h.SendResponse(c, nil, CreateResponse{
		Tag: tag,
	})
}
