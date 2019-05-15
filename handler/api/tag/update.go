package tag

import (
	"fmt"
	"strconv"

	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/errno"
	"hrgdrc/util"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

//Update
func Update(c *gin.Context) {
	log.Info("Update function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	// Get the user id from the url parameter.
	tagID, _ := strconv.Atoi(c.Param("id"))

	var m model.Tag
	if err := c.Bind(&m); err != nil {
		h.SendResponse(c, errno.ErrBind,  err.Error())
		return
	}

	m.ID = uint64(tagID)
	oldTag, err := model.GetTag(m.ID , false ) // 取出旧数据来记录

	if err := m.Update(); err != nil {
		h.SendResponse(c, errno.ErrDatabase,  err.Error())
		return
	}

	parentTag , err := model.GetTagParent(m.ID)
	if err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}
	record := model.Record{}
	record.Body = "描述:系数变更; 系数名：" + parentTag.Name + ";"
	record.Body += "系数变化 ，原系数名：" + oldTag.Name + "; 原系数：" + fmt.Sprint(oldTag.Coefficient) + ";"
	record.Body += "系数变化 ，新系数名：" + m.Name + "; 新系数：" +  fmt.Sprint(m.Coefficient) + ";"

	if err := record.Create(); err != nil {
		fmt.Println(err)
	}

	h.SendResponse(c, nil, nil)
}
