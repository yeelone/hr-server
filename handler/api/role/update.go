package role

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
	log.Info("Role Update function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	// Get the user id from the url parameter.
	rid, _ := strconv.Atoi(c.Param("id"))

	// Binding the user data.
	var m model.Role
	if err := c.Bind(&m); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	m.ID = uint64(rid)

	// Save changed fields.
	if err := m.Update(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	h.SendResponse(c, nil, nil)
}
