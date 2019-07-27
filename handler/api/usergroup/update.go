package usergroup

import (
	"strconv"

	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

//Update
func Update(c *gin.Context) {
	log.Info("UserGroup Update function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	// Get the user id from the url parameter.
	groupID, _ := strconv.Atoi(c.Param("id"))

	// Binding the user data.
	var p model.UserGroup
	if err := c.Bind(&p); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	// We update the record based on the category id.
	p.ID = uint64(groupID)

	// Save changed fields.
	if err := p.Update(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}

	h.SendResponse(c, nil, nil)
}
