package group

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
	"os"
	"strings"
)

func RelateProfiles(c *gin.Context) {
	log.Info("group relateprofiles function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	if err := model.AddGroupProfiles(r.ID, r.Profiles); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}
	group, _ := model.GetGroup(r.ID, false)
	profiles, _ := model.GetProfiles(r.Profiles)
	record := model.Record{}
	record.Body = "描述:组关联人员; 组名：" + group.Name + ";新添加人员："
	for _, p := range profiles {
		record.Body += p.Name + ",身份证：" + p.IDCard + ";"
	}

	if err := record.Create(); err != nil {
		fmt.Println(err)
	}
	h.SendResponse(c, nil, nil)
}

func RemoveRelateProfiles(c *gin.Context) {
	log.Info("group relateprofiles function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}
	if err := model.RemoveGroupProfiles(r.ID, r.Profiles); err != nil {
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	group, _ := model.GetGroup(r.ID, false)
	profiles, _ := model.GetProfiles(r.Profiles)
	record := model.Record{}
	record.Body = "描述:移除组与人员的关联; 组名：" + group.Name + ";移除人员包括："
	for _, p := range profiles {
		record.Body += p.Name + ",身份证：" + p.IDCard + ";"
	}

	if err := record.Create(); err != nil {
		fmt.Println(err)
	}

	h.SendResponse(c, nil, nil)
}

func RelateTags(c *gin.Context) {
	log.Info("RelateTags function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r RelateTagsRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	if err := model.ClearThenAddGroupTags(r.Group, r.Tags); err != nil {
		e := errno.ErrDatabase
		e.Message = e.Message + ";" + err.Error()
		h.SendResponse(c, e, nil)
		return
	}

	//记录
	topTags, _, _ := model.GetAllTags(0, 10000, "parent", "0")
	topTagMap := make(map[uint64]string)
	for _, tag := range topTags {
		topTagMap[tag.ID] = tag.Name
	}

	group, _ := model.GetGroup(r.Group, false)
	tags, _ := model.GetTagsByIDList(r.Tags)
	record := model.Record{}
	record.Body = "描述:群组与标签的关联发化变更; 组名：" + group.Name + ";现所关联系数："
	for _, t := range tags {
		record.Body += topTagMap[t.Parent] + ",系数：" + fmt.Sprint(t.Coefficient) + ";"
	}

	if err := record.Create(); err != nil {
		fmt.Println(err)
	}
	if err := createCSV(group, r.Rules); err != nil {
		fmt.Println("create csv error", err)
		h.SendResponse(c, errno.ErrCreateFile, nil)
		return
	}
	h.SendResponse(c, nil, nil)
}

func createCSV(group *model.Group, rules []string) (err error) {

	filename := "conf/group_tag_rule.csv"
	var f *os.File
	if !util.Exists(filename) {
		f, err = os.Create(filename) //创建文件
		if err != nil {
			panic(err)
		}
		defer f.Close()
	}

	targetGroupName := ""
	if len(rules) > 0 { //如果rules为空，有可能是用户删除了所有的规则
		targetGroupName = strings.Split(rules[0], ",")[0]
	} else {
		parent, _ := model.GetGroup(group.Parent, false)
		targetGroupName = parent.Name + "." + group.Name
	}

	//csv的格式类似：职称.员级, 部门.新亨支行,独生子女费.10
	// 当 职称.员级, 部门.新亨支行 存在时，更新 "独生子女费.10"
	// 不存在时则创建
	newLines := []string{}
	lines, err := util.ReadLines(filename)
	for _, line := range lines {
		s := strings.Split(line, ",")
		if s[0] != targetGroupName {
			newLines = append(newLines, line)
		}
	}

	for _, r := range rules {
		newLines = append(newLines, r)
	}

	util.WriteLines(newLines, filename)

	return err
}
