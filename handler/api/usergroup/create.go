package usergroup

import (
	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/errno"
	"hrgdrc/util"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

// Create :
func Create(c *gin.Context) {
	log.Info("user group Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	m := model.UserGroup{
		Name:   r.Name,
		Parent: r.Parent,
	}

	// Insert the group to the database.
	if err := m.Create(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}

	rsp := CreateResponse{
		Group: &m,
	}

	// Show the user information.
	h.SendResponse(c, nil, rsp)
}
