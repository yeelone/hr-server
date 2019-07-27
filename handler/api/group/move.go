package group

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

//Update
func Move(c *gin.Context) {
	log.Info("Move group function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	// Binding the group data.
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}
	// Save changed fields.
	if err := model.MoveGroup(r.ID, r.Parent); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	group, _ := model.GetGroup(r.ID, false)
	model.CreateOperateRecord(c, fmt.Sprintf("机构调动, 机构名: %s ", group.Name))
	h.SendResponse(c, nil, nil)
}
