package role

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
	log.Info("role Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	m := model.Role{
		Name: r.Name,
	}

	// Insert the group to the database.
	if err := m.Create(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}
	model.CreateOperateRecord(c, fmt.Sprintf("新建Role, Role：%s", r.Name))

	rsp := CreateResponse{
		Role: &m,
	}

	// Show the user information.
	h.SendResponse(c, nil, rsp)
}
