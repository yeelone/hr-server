package templateaccount

import (
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
	"strconv"
)

func GetAccountTemplate(c *gin.Context) {
	log.Info("Get Account Template  function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	id, _ := strconv.Atoi(c.Param("id"))
	t, err := model.GetAccountTemplates(uint64(id))
	if err != nil {
		h.SendResponse(c, errno.ErrGetTemplate, nil)
		return
	}

	h.SendResponse(c, nil, TemplateResponse{Templates: t})
}
