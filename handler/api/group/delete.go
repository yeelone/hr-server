package group

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
	log.Info("group Delete function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	groupID, _ := strconv.Atoi(c.Param("id"))
	group, _ := model.GetGroup(uint64(groupID), false)
	if err := model.DeleteGroup(group.ID); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	model.CreateOperateRecord(c, fmt.Sprintf("删除机构, 机构名: %s ", group.Name))

	h.SendResponse(c, nil, nil)
}
