package group

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

	model.CreateOperateRecord(c, fmt.Sprintf("新建机构, 机构名: %s ", r.Name))
	// Show the user information.
	h.SendResponse(c, nil, rsp)
}
