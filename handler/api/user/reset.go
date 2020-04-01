package user

import (
	"fmt"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/auth"
	"hr-server/pkg/errno"
	"hr-server/util"

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
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	// 先验证旧密码是否正确
	u, err := model.GetUser(r.ID)
	if err != nil {
		log.Warnf("ChangePassword function called. ID: %d | username: %s | info: %s  ", r.ID, "", "user trying to login")
		h.SendResponse(c, errno.ErrUserNotFound, err.Error())
		return
	}

	// Compare the login password with the user password.
	if err := auth.Compare(u.Password , r.OldPassword ); err != nil {
		log.Warnf("Login function called. ID: %d | username: %s | info: %s  ", u.ID, u.Username, "password incorrect")
		h.SendResponse(c, errno.ErrPasswordIncorrect, err.Error())
		return
	}

	if err := model.ChangeUsersPassword(r.ID, r.Password); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}

	h.SendResponse(c, nil, nil)
}
