package salary

import (
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/util"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

func GetBaseSalary(c *gin.Context) {
	log.Info("GetBaseSalary query function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	s := model.SalaryConfig{}
	var err error
	if s, err = model.GetSalaryConfig(); err != nil {
		h.SendResponse(c, nil, CreateResponse{Base: 0.00, TaxThreshold: 0.00})
		return
	}

	h.SendResponse(c, nil, CreateResponse{Base: s.Base, TaxThreshold: s.TaxThreshold})
}
