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

func ResetPassword(c *gin.Context) {
	log.Info("User Active function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		fmt.Println(err)
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	if err := model.ResetUsersPassword(r.Users); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}

	h.SendResponse(c, nil, nil)
}

func ChangePassword(c *gin.Context) {
	log.Info("change user password function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		fmt.Println(err)
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	if err := model.ChangeUsersPassword(r.ID, r.Password); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}

	h.SendResponse(c, nil, nil)
}
