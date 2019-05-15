package salary

import (
	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/errno"
	"hrgdrc/pkg/template"
	"hrgdrc/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

func GetTemplate(c *gin.Context) {
	log.Info("GetTemplate query function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	id, _ := strconv.Atoi(c.Param("id"))
	t, err := model.GetTemplate(uint64(id))
	if err != nil {
		h.SendResponse(c, errno.ErrGetTemplate, nil)
		return
	}

	fields := template.GetFields(t.Name)

	h.SendResponse(c, nil, TemplateResponse{Template: t, Fields: fields})
}

//GetAuditTemplate : 当模板处于待审核的状态时，审核者通过GetAuditTemplate来查看变动的模板详细.
func GetAuditTemplate(c *gin.Context) {
	log.Info("GetTemplate query function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	id, _ := strconv.Atoi(c.Param("id"))
	t, err := model.GetTemplate(uint64(id))
	if err != nil {
		h.SendResponse(c, errno.ErrGetTemplate, nil)
		return
	}

	fields := template.GetFields(t.Name + "-" + util.Uint2Str(t.ID))

	h.SendResponse(c, nil, TemplateResponse{Template: t, Fields: fields})
}
