package usergroup

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
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}
	groups, count, err := model.ListUserGroup(r.Offset, r.Limit, r.WhereKey, r.WhereValue)
	if err != nil {
		h.SendResponse(c, err, err.Error())
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
		h.SendResponse(c, err, err.Error())
		return
	}

	var r ListRequest
	if err := c.ShouldBind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	users, total, err := model.GetUserGroupRelatedUsers(i, r.Offset, r.Limit)
	if err != nil {
		h.SendResponse(c, err, err.Error())
		return
	}

	h.SendResponse(c, nil, ListResponse{
		UserList:   users,
		TotalCount: total,
	})
}
