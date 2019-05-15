package group

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

//Update
func Move(c *gin.Context) {
	log.Info("Move group function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	// Binding the group data.
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}
	// Save changed fields.
	if err := model.MoveGroup(r.ID, r.Parent); err != nil {
		fmt.Println(err.Error())
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	h.SendResponse(c, nil, nil)
}
