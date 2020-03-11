package statistics

import (
	"fmt"
	"github.com/gin-gonic/gin"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
)

// 获取用户档案的增长司长 ，type 可以 分为几种情况：
// year : 获取近几年增长情况
// month : 获取最近12个月的增长情况
// day: 获取最近几天的增长情况
func ProfileIncrease(c *gin.Context){
	var r ProfileRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	fmt.Println(util.PrettyJson(r))

	resp := make([]model.ProfileIncrease, 0)
	var err error
	if r.GetDay {
		resp, err = model.GetIncreaseDay(r.Amount)

		if err != nil {
			h.SendResponse(c, errno.ErrDatabase, err.Error())
			return
		}
	}

	if r.GetMonth {
		resp, err = model.GetIncreaseMonth(r.Amount)

		if err != nil {
			h.SendResponse(c, errno.ErrDatabase, err.Error())
			return
		}
	}

	if r.GetYear {
		resp, err = model.GetIncreaseYear(r.Amount)

		if err != nil {
			h.SendResponse(c, errno.ErrDatabase, err.Error())
			return
		}
	}

	h.SendResponse(c, nil, ProfileResponse{Data: resp})
	return

}