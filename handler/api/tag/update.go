package tag

import (
	"fmt"
	"strconv"

	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

//Update
func Update(c *gin.Context) {
	log.Info("Update function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	// Get the user id from the url parameter.
	tagID, _ := strconv.Atoi(c.Param("id"))

	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		fmt.Println("err", err)
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	m := model.Tag{
		Name:                 r.Name,
		Coefficient:          r.Coefficient,
		Parent:               r.Parent,
		CommensalismGroupIds: util.Uint64ArrayToInt64Array(r.GroupIds),
	}

	m.ID = uint64(tagID)
	oldTag, err := model.GetTag(m.ID, false) // 取出旧数据来记录

	if err := m.Update(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}

	parentTag, err := model.GetTagParent(m.ID)
	if err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}
	record := model.Record{}
	record.Object = "tag"
	record.Body = "描述:系数变更; 系数名:" + parentTag.Name + ";"
	record.Body += "系数变化,原系数名:" + oldTag.Name + "; 原系数:" + fmt.Sprint(oldTag.Coefficient) + ";"
	record.Body += "系数变化,新系数名:" + m.Name + "; 新系数:" + fmt.Sprint(m.Coefficient) + ";"

	if err := record.Create(); err != nil {
		fmt.Println(err)
	}
	model.CreateOperateRecord(c, fmt.Sprintf("更新标签,标签名: %s", oldTag.Name))
	h.SendResponse(c, nil, nil)
}
