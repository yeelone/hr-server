package message

import (
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
)

func UpdateStatus(c *gin.Context) {
	log.Info("UpdateStatus function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r UpdateStatusRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	// todo : 有空再改成批量更新
	for _, id := range r.MsgIds {
		m := model.Message{}
		m.ID = id
		m.SetStatus(r.Status)
	}

	h.SendResponse(c, nil,nil )
	return
}
