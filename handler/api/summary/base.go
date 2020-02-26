package summary

import (
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/util"
)

func Summary(c *gin.Context) {
	log.Info("query summary data : (Summary) function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})

	profileCount, _ := model.CountProfile()
	userCount, _ := model.CountUser()
	groupCount, _ := model.CountUserGroup()
	auditCount, _ := model.CountAudit()
	h.SendResponse(c, nil, CreateResponse{
		ProfileCount: profileCount,
		UserCount:    userCount,
		GroupCount:   groupCount,
		AuditCount:   auditCount,
	})
}
