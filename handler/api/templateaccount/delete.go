package templateaccount

import (
	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := model.DeleteTemplateAccount(uint64(id)); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}

	h.SendResponse(c, nil, nil)
}
