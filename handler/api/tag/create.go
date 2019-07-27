package tag

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
	log.Info("tag Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	m := model.Tag{
		Name:                 r.Name,
		Coefficient:          r.Coefficient,
		Parent:               r.Parent,
		CommensalismGroupIds: util.Uint64ArrayToInt64Array(r.GroupIds),
	}

	// Insert the group to the database.
	if err := m.Create(); err != nil {
		h.SendResponse(c, &errno.Errno{Code: errno.ErrDatabase.Code, Message: err.Error()}, err.Error())
		return
	}

	rsp := CreateResponse{
		Tag: &m,
	}
	model.CreateOperateRecord(c, fmt.Sprintf("创建标签,标签名: %s", r.Name))
	// Show the tag information.
	h.SendResponse(c, nil, rsp)
}
