package tag

import (
	"fmt"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
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
	tag, _ := model.GetTag(uint64(tagID), false)

	if err := model.DeleteTag(tag.ID); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}
	model.CreateOperateRecord(c, fmt.Sprintf("删除标签,标签名: %s", tag.Name))
	h.SendResponse(c, nil, nil)
}
