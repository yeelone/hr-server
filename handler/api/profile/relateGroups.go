package profile

import (
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
)

// 一次性关联多个组
// 这个功能一般只有新增加Profile的时候用到，其它时候要变更员工的组，需要进行“调动”管理
func RelateGroups(c *gin.Context) {
	log.Info("RelateGroups function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r RelateGroupsRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	for _, id := range r.Groups {
		pids := make([]uint64,1)
		pids[0] = r.Profile
		model.AddGroupProfiles(id,pids)
	}

	h.SendResponse(c, nil, nil)
}
