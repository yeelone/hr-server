package profile

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

func Update(c *gin.Context) {
	log.Info("Update function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	// Get the user id from the url parameter.
	profileID, _ := strconv.Atoi(c.Param("id"))

	// Binding the user data.
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}
	profile := r.Profile
	// We update the record based on the user id.
	profile.ID = uint64(profileID)
	oldProfile,_  := model.GetProfile(profile.ID)
	//取出原档案数据 ，跟新的进行对比，判断更新的字段
	change := util.FindUpdatedField(oldProfile, profile)
	// Validate the data.
	if err := profile.Validate(); err != nil {
		h.SendResponse(c, errno.ErrValidation, nil)
		return
	}

	// Save changed fields.
	if err := profile.Update(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	uid, _ := c.Get("userid")
	audit := &model.Audit{}
	audit.OperatorID = uid.(uint64)
	audit.Object = model.ProfileAuditObject
	audit.Action = model.AUDITUPDATEACTION
	audit.OrgObjectID = []int64{int64(profile.ID)}
	audit.State = model.AuditStateWaiting
	body := ""
	for k, v := range change {
		if k != "audit_state" && k != "groups" { // 这两个字段不需要进行对比
			body  = "更新了:" + model.ProfileI18nMap[k] + ";"
			body += "更新内容:从[" + fmt.Sprint(v["from"]) + "]到[" + fmt.Sprint(v["to"]) + "];"
		}
	}
	audit.Body = "描述:更新职工档案;" +
		"档案ID:" + util.Uint2Str(profile.ID) + "; " +
		"职工姓名:" + profile.Name + ";" +
		"身份证号码:" + profile.IDCard + ";" + body
	audit.Remark = r.Remark

	if err := audit.Create(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	h.SendResponse(c, nil, nil)
}
