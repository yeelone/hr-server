package profile

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/errno"
	"hrgdrc/util"
)

func Freeze(c *gin.Context) {
	log.Info("Lock function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r FreezeRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}
	// Save changed fields.
	if err := model.FreezeProfile(r.Profiles); err != nil {
		h.SendResponse(c, errno.ErrFreezeProfile, err.Error())
		return
	}
	//创建的同时需同时创建审核条目
	userid, _ := c.Get("userid")
	//查底所有用户资料，记录在audit中
	profiles,_ := model.GetProfiles(r.Profiles)
	body := ""

	for _,p := range profiles {
		body += "冻结用户:" + p.Name + ",证件:"+ p.IDCard + "\n"
	}

	ids := []int64{}
	for _,id := range r.Profiles{
		ids = append(ids, int64(id))
	}
	audit := &model.Audit{}
	audit.OperatorID = userid.(uint64)
	audit.Object = model.ProfileAuditObject
	audit.Action = model.AUDITUPDATEACTION
	audit.OrgObjectID = ids
	audit.State = model.AuditStateWaiting
	audit.Body = "描述:冻结职工档案;" +body
	audit.Remark = r.Remark

	if err := audit.Create(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}

	h.SendResponse(c, nil, nil)
}

func UnFreeze(c *gin.Context) {
	log.Info("UnFreeze function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r FreezeRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}
	// Save changed fields.
	if err := model.UnFreezeProfile(r.Profiles); err != nil {
		h.SendResponse(c, errno.ErrFreezeProfile, err.Error())
		return
	}
	//创建的同时需同时创建审核条目
	profiles,_ := model.GetProfiles(r.Profiles)
	body := ""

	for _,p := range profiles {
		body += "冻结用户:" + p.Name + ",证件:"+ p.IDCard + "\n"
	}

	userid, _ := c.Get("userid")
	ids := []int64{}
	for _,id := range r.Profiles{
		ids = append(ids, int64(id))
	}
	audit := &model.Audit{}
	audit.OperatorID = userid.(uint64)
	audit.Object = model.ProfileAuditObject
	audit.Action = model.AUDITUPDATEACTION
	audit.OrgObjectID = ids
	audit.State = model.AuditStateWaiting
	audit.Body = "描述:激活职工档案;" + body
	audit.Remark = r.Remark

	if err := audit.Create(); err != nil {
		fmt.Println(err)
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}

	h.SendResponse(c, nil, nil)
}