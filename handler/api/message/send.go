package message

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
)

func Send(c *gin.Context) {
	log.Info("send message function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	m := model.MessageText{
		SendId: r.SendId,
		Title: r.Title,
		Text: r.Text,
		MType: r.MType,
		Group: r.Group,
	}

	// Global 消息只有管理员可以发送
	if m.MType == "Global" {
    	user,err  := model.GetUser(m.SendId)
		if err != nil {
			h.SendResponse(c, errno.ErrUserNotFound, nil)
			return
		}

		if !user.IsSuper {
			h.SendResponse(c, errno.ErrSendGlobalMessage, nil)
			return
		}
	}

	if err := m.Create(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	// 如果是私信的话，要在message表中插入一行记录
	if m.MType == "Private" {
		if r.RecId == 0 {
			h.SendResponse(c, errors.New("请指定收件人"), nil)
			return
		}

		msg := model.Message{}
		msg.RecId = r.RecId
		msg.MType = m.MType
		msg.Status = 0
		msg.TextId = m.ID

		if err := msg.Create(); err != nil {
			model.DeleteMessageText(m.ID)
			h.SendResponse(c, errno.ErrDatabase, nil)
			return
		}
	}

	rsp := CreateResponse{
	}

	h.SendResponse(c, nil, rsp)
}

