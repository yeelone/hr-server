package salary

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
)

// TemplateConfig :
func TemplateOrder(c *gin.Context) {
	log.Info("TemplateOrder Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r TemplateOrderRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}
	if err := model.UpdateTemplateOrder(r.Orders); err != nil {
		fmt.Println("create error", err)
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}

	h.SendResponse(c, nil, nil)
}
