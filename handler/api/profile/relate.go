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

var tempProfile = model.Profile{}

func RelateTags(c *gin.Context) {
	log.Info("RelateTags function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r RelateTagsRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	//将父标签取出来
	topTags, _, _ := model.GetAllTags(0, 10000, "parent", "0")
	topTagMap := make(map[uint64]string)
	for _, tag := range topTags {
		topTagMap[tag.ID] = tag.Name
	}

	removeTagStr, _ := writeProfileRecord(topTagMap, r.Profile)

	if err := model.ClearThenAddProfileTags(r.Profile, r.Tags); err != nil {
		fmt.Println(err)
		h.SendResponse(c, errno.ErrDatabase, err)
		return
	}
	addTagStr, _ := writeProfileRecord(topTagMap, r.Profile)
	record := model.Record{}
	record.Object = "profile"
	record.Body = "描述:标签关联变动; 姓名:" + tempProfile.Name + ";身份证号码:" + tempProfile.IDCard + ";"

	if len(removeTagStr) > 0 {
		record.Body += "系数变化:删除了以下系数;" + removeTagStr
	}

	if len(addTagStr) > 0 {
		record.Body += "系数变化:新增了以下系数;" + addTagStr
	}
	if err := record.Create(); err != nil {
		fmt.Println(err)
	}
	model.CreateOperateRecord(c, fmt.Sprintf("员工关联标签, 员工信息：[ %s ]", tempProfile.Name))

	h.SendResponse(c, nil, nil)
}

func writeProfileRecord(tagmap map[uint64]string, pid uint64) (body string, err error) {
	profile, err := model.GetProfileWithTags(pid)
	tempProfile.Name = profile.Name
	tempProfile.IDCard = profile.IDCard
	if err != nil {
		return "", nil
	}

	recordStr := ""
	for _, tag := range profile.Tags {
		recordStr += "系数名:" + tagmap[tag.Parent] + "; 系数:" + fmt.Sprint(tag.Coefficient) + ";"
	}

	return recordStr, err
}
