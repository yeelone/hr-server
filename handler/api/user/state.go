package user

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

func Freeze(c *gin.Context) {
	log.Info("User Freeze function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	if err := model.FreezeUsers(r.Users); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}
	model.CreateOperateRecord(c, fmt.Sprintf("冻结用户, 用户列表: [ %s ]", util.ArrayToString(r.Users, ",")))
	h.SendResponse(c, nil, nil)
}

func Active(c *gin.Context) {
	log.Info("User Active function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		fmt.Println(err)
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	if err := model.ActiveUsers(r.Users); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}
	model.CreateOperateRecord(c, fmt.Sprintf("激活用户, 用户列表: [ %s ]", util.ArrayToString(r.Users, ",")))
	h.SendResponse(c, nil, nil)
}
