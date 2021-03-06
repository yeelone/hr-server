package profile

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

func Create(c *gin.Context) {
	log.Info("Profile Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	userid, ok := c.Get("userid")
	if !ok {
		h.SendResponse(c, errno.StatusUnauthorized, nil)
		return
	}
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	profile := r.Profile

	// Validate the data.
	if err := profile.Validate(); err != nil {
		h.SendResponse(c, errno.ErrValidation, nil)
		return
	}

	profile.AuditState = model.AuditStateWaiting

	// Insert the profile to the database.
	if err := profile.Create(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	////创建职工档案的时候，如果没有指定的组，只默认分配到系统默认创建的第一个组中.
	if err := model.AddProfileToDefaultGroup(profile.ID); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	//创建的同时需同时创建审核条目
	audit := &model.Audit{}
	audit.OperatorID = userid.(uint64)
	audit.Object = model.ProfileAuditObject
	audit.Action = model.AUDIT_CREATE_ACTION
	audit.OrgObjectID = []int64{int64(profile.ID)}
	audit.State = model.AuditStateWaiting
	audit.Body = "描述:创建职工档案;" +
		"档案ID:" + util.Uint2Str(profile.ID) + "; " +
		"员工姓名:" + profile.Name + ";" +
		"身份证号码:" + profile.IDCard + ";"
	audit.Remark = r.Remark

	if err := audit.Create(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	model.CreateOperateRecord(c, fmt.Sprintf("新建员工信息, 员工信息： %s", profile.Name))


	// 消息提示
	// 发送给复核岗
	role , err := model.GetRoleByName("复核岗")

	if err == nil {
		m := model.MessageText{
			SendId: userid.(uint64),
			Title: "有新的审核,请尽快处理",
			Text: "新员工信息",
			MType: "Public",
			Role:role.ID,
		}

		m.Create()
	}


	rsp := CreateResponse{
		ID:profile.ID,
		Name: profile.Name,
	}

	// Show the user information.
	h.SendResponse(c, nil, rsp)
}
