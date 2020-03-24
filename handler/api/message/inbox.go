package message

import (
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
	"strconv"
)

// 收件箱未读统计
func InboxCount(c *gin.Context) {
	log.Info("Inbox message function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	recId, _ := strconv.Atoi(c.Param("id"))
	private ,public , global := model.CheckUserMessage(uint64(recId), 0 )

	response := InboxCountResponse{}
	response.Private = private
	response.Global = global
	response.Public = public

	h.SendResponse(c, nil,response)
	return
}

func Inbox(c *gin.Context) {
	log.Info("Inbox message function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	recId, _ := strconv.Atoi(c.Param("id"))
	var r ListRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	//userid, _ := c.Get("userid")

	//if userid == nil {
	//	h.SendResponse(c, errno.StatusUnauthorized, errors.New("无权查看"))
	//	return
	//}
	//
	//if uint64(recId) != userid.(uint64) {
	//	h.SendResponse(c, errno.StatusUnauthorized, errors.New("无权查看"))
	//	return
	//}

	list,total, err := model.GetMessages(r.Offset, r.Limit, uint64(recId),"status", "0")

	if err != nil {
		h.SendResponse(c, errno.ErrDatabase,err )
		return
	}


	respItems := make([]CreateResponse,0)

	for _, item := range list {
		user, _ := model.GetUser(item.SendId)

		resp := CreateResponse{}
		resp.Id = item.MessageId
		resp.SenderId = user.ID
		resp.SenderName = user.Username
		resp.MType = item.MType
		resp.Title = item.Title
		resp.Text = item.Text
		resp.Date = item.PostDate.String()

		if item.Group != 0 {
			g, _ := model.GetUserGroup(item.Group, false)
			resp.GroupName = g.Name
		}

		if item.Role != 0 {
			r,_ := model.GetRole(item.Role,false)
			resp.RoleName = r.Name
		}

		respItems = append(respItems,resp)
	}

	h.SendResponse(c, nil,ListResponse{
		Total:total,
		List:respItems,
	})
	return

}