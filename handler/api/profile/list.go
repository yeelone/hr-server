package profile

import (
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

func List(c *gin.Context) {
	log.Info("List function called.")
	var r ListRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}
	infos, count, err := model.ListProfile(r.Key, r.Value, r.Offset, r.Limit, r.Freezed)
	if err != nil {
		h.SendResponse(c, err, nil)
		return
	}

	// 给 id 加密

	for _, p := range infos {
		str := strconv.FormatUint(p.ID, 10)
		p.UUID = util.HashIDEncode(str)
	}

	h.SendResponse(c, nil, ListResponse{
		TotalCount:  count,
		ProfileList: infos,
	})
}
