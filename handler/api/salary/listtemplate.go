package salary

import (
	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/errno"

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
