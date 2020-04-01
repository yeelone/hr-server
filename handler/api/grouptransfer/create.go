package grouptransfer

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

// Create :
func Create(c *gin.Context) {
	log.Info("group transfer Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	m := model.GroupTransfer{
		Profile:     r.ProfileID,
		OldGroup:    r.OldGroupID,
		NewGroup:    r.NewGroupID,
		Description: r.Description,
	}
	// Insert the group to the database.
	if err := m.Create(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}
	profile, _ := model.GetProfile(r.ProfileID)
	oldG, err := model.GetGroup(r.OldGroupID, false)
	if err != nil {
		// 员工归属的组有可能出现为空的情况
		oldG = &model.Group{}
		oldG.Name = ""
	}
	newG, _ := model.GetGroup(r.NewGroupID, false)

	//创建的同时需同时创建审核条目
	userid, _ := c.Get("userid")
	audit := &model.Audit{}
	audit.OperatorID = userid.(uint64)
	audit.Object = model.ProfileAuditObject
	audit.Action = model.AUDIT_MOVE_ACTION
	audit.OrgObjectID = []int64{int64(profile.ID)}
	audit.DestObjectID = []int64{int64(m.ID)}
	audit.State = model.AuditStateWaiting
	audit.Body = "描述:职工部门调动;" +
		"档案ID:" + util.Uint2Str(profile.ID) + "; " +
		"员工姓名:" + profile.Name + ";" +
		"身份证号码:" + profile.IDCard + ";" +
		"从:" + oldG.Name + ";" +
		"到:" + newG.Name + ";"

	audit.Remark = r.Remark

	if err := audit.Create(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	profile.UpdateState(model.AuditStateWaiting)
	model.CreateOperateRecord(c, fmt.Sprintf("人员调动, 姓名: %s ", profile.Name))

	// 消息提示
	role , err := model.GetRoleByName("复核岗")
	if err == nil {
		m := model.MessageText{
			SendId: userid.(uint64),
			Title: "有新的审核,请尽快处理",
			Text: "职工部门调动",
			MType: "Public",
			Role:role.ID,
		}

		m.Create()
	}

	// Show the tag information.
	h.SendResponse(c, nil, nil)
}
