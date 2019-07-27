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
	audit.Action = model.AUDITDELETEACTION
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
	model.CreateOperateRecord(c, fmt.Sprintf("删除员工信息, 员工信息： %s", profile.Name))
	h.SendResponse(c, nil, nil)
}
