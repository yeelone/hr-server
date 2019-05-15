package group

import (
	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/errno"
	"hrgdrc/util"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	group, err := model.GetGroup(id, false)
	if err != nil {
		e := errno.ErrGroupNotFound
		e.Message = e.Message + ";" + err.Error()
		h.SendResponse(c, errno.ErrGroupNotFound, nil)
		return
	}

	//获取组时,同时解析群与标签的规则关系
	//todo : 是否需要拆开成独立的API
	rules, _ := getRulesFromCSV(group)
	h.SendResponse(c, nil, CreateResponse{Group: group,Rules:rules})
}

func getRulesFromCSV(group *model.Group) (rules []string , err error) {
	parent,_ := model.GetGroup(group.Parent, false)
	name :=  parent.Name + "." + group.Name
	filename := "conf/group_tag_rule.csv"
	var f *os.File
	if !util.Exists(filename) {
		f, err = os.Create(filename) //创建文件
		if err != nil {
			panic(err)
		}
		defer f.Close()
	}

	//csv的格式类似：职称.员级, 部门.新亨支行,独生子女费.10
	// 当 职称.员级, 部门.新亨支行 存在时，更新 "独生子女费.10"
	lines, err := util.ReadLines(filename)
	for _, line := range lines {
		s := strings.Split(line, ",")
		if s[0] == name {
			rules = append(rules, line)
		}
	}

	return rules, err
}
