package tag

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
	log.Info("tag Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	m := model.Tag{
		Name:        r.Name,
		Coefficient: r.Coefficient,
		Parent:      r.Parent,
	}

	// Insert the group to the database.
	if err := m.Create(); err != nil {
		h.SendResponse(c, &errno.Errno{Code: errno.ErrDatabase.Code, Message: err.Error()}, err.Error())
		return
	}

	rsp := CreateResponse{
		Tag: &m,
	}

	// Show the tag information.
	h.SendResponse(c, nil, rsp)
}
