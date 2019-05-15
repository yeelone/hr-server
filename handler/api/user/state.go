package user

import (
	"fmt"
	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/errno"
	"hrgdrc/util"

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

	h.SendResponse(c, nil, nil)
}
