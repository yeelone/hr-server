package user

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

func Update(c *gin.Context) {
	log.Info("Update function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	// Get the user id from the url parameter.
	userID, _ := strconv.Atoi(c.Param("id"))

	// Binding the user data.
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	// We update the record based on the user id.

	u := model.User{
		Email:    r.Email,
		IsSuper:  r.IsSuper,
		Username: r.Username,
		Nickname: r.Nickname,
		IDCard:   r.IDCard,
	}
	u.ID = uint64(userID)
	// Validate the data.
	if err := u.Validate(); err != nil {
		fmt.Println(err)
		h.SendResponse(c, errno.ErrValidation, err.Error())
		return
	}

	// Save changed fields.
	if err := u.Update(); err != nil {
		fmt.Println(err)
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}

	if r.Group > 0 {
		if err := model.AddUserGroupUsers(r.Group, []uint64{u.ID}); err != nil {
			h.SendResponse(c, errno.ErrDatabase, err.Error())
			return
		}
	} else {
		//创建用户的时候，如果没有指定的组，只默认分配到系统默认创建的第一个组中.
		if err := model.AddUserToDefaultGroup(u.ID); err != nil {
			h.SendResponse(c, errno.ErrDatabase, err.Error())
			return
		}
	}
	model.CreateOperateRecord(c, fmt.Sprintf("更新用户, 用户: %s ", u.Username))
	h.SendResponse(c, nil, nil)
}
