package salary

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	h "hr-server/handler"
	"hr-server/util"
)

// GetPayroll : 获取工资单，职工只能查询自己的工资单
func GetPayroll(c *gin.Context) {
	log.Info("GetBaseSalary query function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	uid, ok := c.Get("userid")
	fmt.Println("uid, ok ", uid, ok)

	h.SendResponse(c, nil, CreateResponse{})
}
