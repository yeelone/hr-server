package templateaccount

import (
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/pkg/template"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListTemplateWithFields(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	templateAccount, err := model.GetTemplateAccount(uint64(id))
	if err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}

	resp := ListResponse{}

	for _, t := range templateAccount.Templates {
		templatesResp := TemplateResponse{}
		fields := template.GetFields(t.Name)
		templatesResp.Name = t.Name
		templatesResp.Fields = fields
		resp.Templates = append(resp.Templates, templatesResp)
	}

	h.SendResponse(c, nil, resp)
}
