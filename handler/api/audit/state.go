package audit

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

func UpdateState(c *gin.Context) {
	log.Info("Update audit state function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	// Get the user id from the url parameter.
	userid, ok := c.Get("userid")
	if !ok {
		h.SendResponse(c, errno.StatusUnauthorized, nil)
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	var r CreateRequest

	if err := c.Bind(&r); err != nil {
		fmt.Println("err", err)
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}
	audit := &model.Audit{}
	audit.ID = uint64(id)
	audit.State = r.State
	audit.Reply = r.Reply
	audit.AuditorID = userid.(uint64)

	// Save changed fields.
	if err := audit.Update(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	model.CreateOperateRecord(c, fmt.Sprintf("更新审核,ID :  %s ", audit.ID))
	h.SendResponse(c, nil, nil)
}
