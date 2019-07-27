package salary

import (
	"fmt"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

func BaseConfig(c *gin.Context) {
	log.Info("TemplateConfig Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}
	m := model.SalaryConfig{
		Base:         r.BaseSalary,
		TaxThreshold: r.TaxThreshold,
	}

	if err := m.Create(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}
	record := model.Record{}
	record.Body = "基本工资发生变更, 变更为[" + fmt.Sprint(r.BaseSalary) + "],请仔细核对！"
	if err := record.Create(); err != nil {

	}
	h.SendResponse(c, nil, nil)
}
