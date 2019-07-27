package templateaccount

import (
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	template, err := model.GetTemplateAccount(id)
	if err != nil {
		h.SendResponse(c, errno.ErrGetTemplateAccount, err.Error())
		return
	}

	h.SendResponse(c, nil, CreateResponse{
		TemplateAccount: template,
	})
}
