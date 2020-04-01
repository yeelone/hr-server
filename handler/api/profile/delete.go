package profile

import (
	"fmt"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

//此API废除
func Delete(c *gin.Context) {
	userid, ok := c.Get("userid")
	if !ok {
		h.SendResponse(c, errno.StatusUnauthorized, nil)
		return
	}

	profileID, _ := strconv.Atoi(c.Param("id"))
	profile, err := model.GetProfile(uint64(profileID))
	if err != nil {
		h.SendResponse(c, nil, nil)
		return
	}

	var r DeleteRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	profile.UpdateState(model.AuditStateWaiting)

	uid, _ := c.Get("userid")
	audit := &model.Audit{}
	audit.OperatorID = uid.(uint64)
	audit.Object = model.ProfileAuditObject
	audit.Action = model.AUDIT_DELETE_ACTION
	//audit.OrgObjectID = profile.ID
	audit.State = model.AuditStateWaiting
	audit.Body = "描述:删除职工档案;" +
		"档案ID:" + util.Uint2Str(profile.ID) + "; " +
		"职工姓名:" + profile.Name + ";" +
		"身份证号码:" + profile.IDCard + ";"
	audit.Remark = r.Remark

	if err := audit.Create(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	// 消息提示
	role , err := model.GetRoleByName("复核岗")
	if err == nil {
		m := model.MessageText{
			SendId: userid.(uint64),
			Title: "有新的审核,请尽快处理",
			Text: "删除员工信息",
			MType: "Public",
			Role:role.ID,
		}

		m.Create()
	}

	model.CreateOperateRecord(c, fmt.Sprintf("删除员工信息, 员工信息： %s", profile.Name))
	h.SendResponse(c, nil, nil)
}
