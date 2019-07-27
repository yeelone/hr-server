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

// @Summary Add new user to the database
// @Description Add a new user
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body user.CreateRequest true "Create a new user"
// @Success 200 {object} user.CreateResponse "{"code":0,"message":"OK","data":{"username":"kong"}}"
// @Router /user [post]
func Create(c *gin.Context) {
	log.Info("User Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	u := model.User{
		Email:    r.Email,
		IsSuper:  r.IsSuper,
		Picture:  r.Picture,
		Username: r.Username,
		IDCard:   r.IDCard,
		Nickname: r.Nickname,
		Password: r.Password,
	}

	// Validate the data.
	if err := u.Validate(); err != nil {
		h.SendResponse(c, errno.ErrValidation, err.Error())
		return
	}

	// Encrypt the user password.
	if err := u.Encrypt(); err != nil {
		h.SendResponse(c, errno.ErrEncrypt, err.Error())
		return
	}
	// Insert the user to the database.
	if err := u.Create(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}

	model.CreateOperateRecord(c, fmt.Sprintf("创建用户:  %s ", u.Username))
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

	rsp := CreateResponse{
		Username: r.Username,
	}

	// Show the user information.
	h.SendResponse(c, nil, rsp)
}
