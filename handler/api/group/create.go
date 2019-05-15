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

// Create :
func Create(c *gin.Context) {
	log.Info("group Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		fmt.Println(err)
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	m := model.Group{
		Name:        r.Name,
		Code:        r.Code,
		Parent:      r.Parent,
		Coefficient: r.Coefficient,
	}

	// Insert the group to the database.
	if err := m.Create(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	rsp := CreateResponse{
		Group: &m,
	}

	// Show the user information.
	h.SendResponse(c, nil, rsp)
}
