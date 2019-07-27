package role

import (
	"fmt"
	"strconv"

	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

// Delete :
func Delete(c *gin.Context) {
	log.Info("role Delete function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	rid, _ := strconv.Atoi(c.Param("id"))
	role, _ := model.GetRole(uint64(rid), false)
	if err := model.DeleteRole(role.ID); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	model.CreateOperateRecord(c, fmt.Sprintf("删除Role, Role：%s", role.Name))
	h.SendResponse(c, nil, nil)
}
