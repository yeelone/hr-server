package group

import (
	"strconv"

	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/errno"
	"hrgdrc/util"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

// Delete :
func Delete(c *gin.Context) {
	log.Info("group Delete function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	groupID, _ := strconv.Atoi(c.Param("id"))
	if err := model.DeleteGroup(uint64(groupID)); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	h.SendResponse(c, nil, nil)
}
