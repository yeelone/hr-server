package templateaccount

import (
	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/errno"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	var r ListRequest
	if err := c.ShouldBind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	templateaccounts, count, err := model.ListTemplateAccounts(r.Offset, r.Limit)
	if err != nil {
		h.SendResponse(c, err, err.Error())
		return
	}

	h.SendResponse(c, nil, ListResponse{
		TotalCount: count,
		List:       templateaccounts,
	})
}
