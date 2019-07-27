package tag

import (
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

func List(c *gin.Context) {
	var r ListRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	tags, count, err := model.GetAllTags(r.Offset, r.Limit, r.WhereKey, r.WhereValue)
	if err != nil {
		h.SendResponse(c, err, err.Error())
		return
	}

	h.SendResponse(c, nil, ListResponse{
		TotalCount: count,
		TagList:    tags,
	})
}

//ProfileList:
func ProfileList(c *gin.Context) {
	log.Info("List function called.")
	id := c.Param("id")
	i, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		h.SendResponse(c, err, err.Error())
		return
	}

	var r ListRequest
	if err := c.ShouldBind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	profiles, total, err := model.GetTagRelatedProfiles(i, r.Offset, r.Limit)
	if err != nil {
		h.SendResponse(c, err, err.Error())
		return
	}

	h.SendResponse(c, nil, ListResponse{
		ProfileList: profiles,
		TotalCount:  total,
	})
}
