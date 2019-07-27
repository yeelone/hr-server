package record

import (
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

func List(c *gin.Context) {
	log.Info("Audit List function called.")

	var r ListRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}
	records, count, err := model.GetRecordListByMonth(r.QueryDate, r.Offset, r.Limit)
	if err != nil {
		h.SendResponse(c, err, err.Error())
		return
	}

	h.SendResponse(c, nil, ListResponse{
		TotalCount: count,
		List:       records,
	})
}
