package usergroup

import (
	"strconv"

	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

// Delete :
func Delete(c *gin.Context) {
	log.Info("user group Delete function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	groupID, _ := strconv.Atoi(c.Param("id"))
	if err := model.DeleteUserGroup(uint64(groupID)); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}

	h.SendResponse(c, nil, nil)
}
