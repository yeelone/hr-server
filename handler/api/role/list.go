package role

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
	roles, count, err := model.ListRoles(r.Offset, r.Limit, r.WhereKey, r.WhereValue)
	if err != nil {
		h.SendResponse(c, err, nil)
		return
	}

	h.SendResponse(c, nil, ListResponse{
		TotalCount: count,
		RoleList:   roles,
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

	users, total, err := model.GetRoleRelatedUsers(i, r.Offset, r.Limit)
	if err != nil {
		h.SendResponse(c, err, nil)
		return
	}

	h.SendResponse(c, nil, ListResponse{
		UserList:   users,
		TotalCount: total,
	})
}
