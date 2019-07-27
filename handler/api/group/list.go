package group

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
	if err := c.ShouldBind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}
	groups, count, err := model.ListGroup(r.Offset, r.Limit, r.WhereKey, r.WhereValue)
	if err != nil {
		h.SendResponse(c, err, nil)
		return
	}

	h.SendResponse(c, nil, ListResponse{
		TotalCount: count,
		GroupList:  groups,
	})
}

func UserList(c *gin.Context) {
	log.Info("List function called.")
	id := c.Param("id")
	i, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		h.SendResponse(c, err, nil)
		return
	}

	var r ListRequest
	if err := c.ShouldBind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	users, total, err := model.GetGroupRelatedUsers(i, r.Offset, r.Limit)
	if err != nil {
		h.SendResponse(c, err, nil)
		return
	}

	h.SendResponse(c, nil, ListResponse{
		UserList:   users,
		TotalCount: total,
	})
}

//ProfileList:
func ProfileList(c *gin.Context) {
	log.Info("Profile List function called.")
	id := c.Param("id")
	i, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		h.SendResponse(c, err, nil)
		return
	}

	var r ListRequest
	if err := c.ShouldBind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	profiles, total, err := model.GetGroupRelatedProfiles(i, r.Offset, r.Limit, r.Freezed)
	if err != nil {
		h.SendResponse(c, err, nil)
		return
	}

	h.SendResponse(c, nil, ListResponse{
		ProfileList: profiles,
		TotalCount:  total,
	})
}
