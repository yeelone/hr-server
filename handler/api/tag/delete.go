package tag

import (
	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/errno"
	"hrgdrc/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

// Delete :
func DeleteList(c *gin.Context) {
	log.Info("tag delete function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	if err := model.DeleteTags(r.IDS); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}

	h.SendResponse(c, nil, nil)
}

// Delete :
func Delete(c *gin.Context) {
	log.Info("tag Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	tagID, _ := strconv.Atoi(c.Param("id"))
	if err := model.DeleteTag(uint64(tagID)); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}

	h.SendResponse(c, nil, nil)
}
