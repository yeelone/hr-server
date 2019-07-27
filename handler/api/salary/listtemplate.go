package salary

import (
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"

	"github.com/gin-gonic/gin"
)

// TemplateConfig :
func ListTemplate(c *gin.Context) {
	list, err := model.ListTemplates()

	if err != nil {
		h.SendResponse(c, errno.ErrListTemplate, nil)
	}

	h.SendResponse(c, nil, ListTemplateResponse{
		List: list,
	})
}
