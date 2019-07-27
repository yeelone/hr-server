package tag

import (
	"fmt"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

func RelateProfiles(c *gin.Context) {
	log.Info("tag relate profiles function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	if err := model.AddTagProfiles(r.ID, r.Profiles); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}
	profiles, _ := model.GetProfiles(r.Profiles)
	pTag, _ := model.GetTagParent(r.ID)
	tag, _ := model.GetTag(r.ID, false)
	record := model.Record{}
	record.Body = "描述:职工档案关联系数; 系数名：[" + pTag.Name + "-" + tag.Name + "];系数：" + fmt.Sprint(tag.Coefficient) + ";"
	for _, p := range profiles {
		record.Body += "系数变化,涉及职工：" + p.Name + ",身份证：" + p.IDCard + ";"
	}

	if err := record.Create(); err != nil {
		fmt.Println(err)
	}

	h.SendResponse(c, nil, nil)
}

func RemoveRelateProfiles(c *gin.Context) {
	log.Info("group relateprofiles function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	if err := model.RemoveTagProfiles(r.ID, r.Profiles); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}

	profiles, _ := model.GetProfiles(r.Profiles)
	tag, _ := model.GetTagParent(r.ID)
	record := model.Record{}
	record.Body = "描述:取消职工与系数的关联; 系数名：" + tag.Name + ";系数：" + fmt.Sprint(tag.Coefficient) + ";"
	for _, p := range profiles {
		record.Body += "系数变化,涉及职工：" + p.Name + ",身份证：" + p.IDCard + ";"
	}

	if err := record.Create(); err != nil {
		fmt.Println(err)
	}

	h.SendResponse(c, nil, nil)
}
