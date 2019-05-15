package group

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
	groupID, _ := strconv.Atoi(c.Param("id"))

	// Binding the user data.
	var p model.Group
	if err := c.Bind(&p); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}
	// We update the record based on the category id.
	p.ID = uint64(groupID)
	// Save changed fields.
	if err := p.Update(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}
	group, _ := model.GetGroup(p.ID, false)
	record := model.Record{}
	record.Body = "描述:组发生了变更; 组名：" + group.Name + ";变更后的系数:" + fmt.Sprint(group.Coefficient) + ";"

	if err := record.Create(); err != nil {
		fmt.Println(err)
	}
	h.SendResponse(c, nil, nil)
}
