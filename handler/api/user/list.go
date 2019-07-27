package user

import (
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary List the users in the database
// @Description List users
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body user.ListRequest true "List users"
// @Success 200 {object} user.SwaggerListResponse "{"code":0,"message":"OK","data":{"totalCount":1,"userList":[{"id":0,"username":"admin","random":"user 'admin' get random string 'EnqntiSig'","password":"$2a$10$veGcArz47VGj7l9xN7g2iuT9TF21jLI1YGXarGzvARNdnt4inC9PG","createdAt":"2018-05-28 00:25:33","updatedAt":"2018-05-28 00:25:33"}]}}"
// @Router /user [get]
func List(c *gin.Context) {
	log.Info("List function called.")
	var r ListRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}
	infos, count, err := model.ListUser(r.Username, r.Offset, r.Limit)
	if err != nil {
		h.SendResponse(c, err, err.Error())
		return
	}

	h.SendResponse(c, nil, ListResponse{
		TotalCount: count,
		UserList:   infos,
	})
}
