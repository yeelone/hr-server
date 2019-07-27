package templateaccount

import (
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

func Create(c *gin.Context) {
	log.Info("TemplateAccount Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	m := &model.TemplateAccount{}
	m.ID = r.ID
	m.Name = r.Name
	m.Order = r.Order
	if err := m.Save(); err != nil {
		h.SendResponse(c, errno.ErrCreateTemplateAccount, err.Error())
	}
	if err := model.AddTemplateAccountRelateGroups(m.ID, r.Groups); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}
	if err := model.AddTemplateAccountRelateTemplates(m.ID, r.Templates); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}
	rsp := CreateResponse{}
	h.SendResponse(c, nil, rsp)
}
