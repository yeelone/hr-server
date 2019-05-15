package role

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
	log.Info("role Delete function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	rid, _ := strconv.Atoi(c.Param("id"))
	if err := model.DeleteRole(uint64(rid)); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	h.SendResponse(c, nil, nil)
}
