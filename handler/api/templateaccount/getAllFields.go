package templateaccount

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
	"strconv"
)

// GetAccountFields
// 根据账套ID 从已核算过的工资表里取出所有模板信息，以及相关联的字段信息，为用户作统计用
func GetAccountFields(c *gin.Context) {
	log.Info("GetAccountFields  function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	id, _ := strconv.Atoi(c.Param("id"))
	year := c.Param("year")
	// 第一步
	account, err := model.GetTemplateAccount(uint64(id))
	if err != nil {
		h.SendResponse(c, errno.ErrGetTemplateAccount, nil)
		return
	}

	// 取出账套里的随便一个组，再从组里随机挑选一个profile出来。这是为了简化搜索
	group := model.Group{}

	if len(account.Groups) > 0 {
		group = account.Groups[0]
	}
	//todo 为了随机挑选一个人出来却查询所有的人，这里要优化。目前时间比较紧
	group, err = model.GetGroupWithProfile(group.ID, true)
	if err != nil {
		fmt.Println("err", err)
	}
	profiles := []uint64{}
	if len(group.Profiles) > 0 {
		profiles = append(profiles, group.Profiles[0].ID)
	}

	salaries, err := model.GetSalaryByAccount(year, account.ID)
	if err != nil {
		h.SendResponse(c, errno.ErrGetTemplate, nil)
		return
	}
	// map[string][string][]string   => month => template name => related fields
	templates := make(map[string]map[string][]string)
	salaryMap := make(map[uint64]string)
	salaryIds := []uint64{}
	templateOrder := []string{}
	for _, salary := range salaries {
		salaryMap[salary.ID] = salary.Template
		salaryIds = append(salaryIds, salary.ID)
		templateOrder = append(templateOrder, salary.Template)
	}

	fields, err := model.GetFieldsBySalaryAndProfilesAndYear(year, salaryIds, profiles)

	for _, field := range fields {
		if _, ok := salaryMap[field.SalaryID]; ok {
			name := salaryMap[field.SalaryID]
			if len(field.Name) == 0 || field.Name == "姓名" || field.Name == "身份证号码" {
				continue
			}
			if _, ok := templates[field.Month]; !ok {
				templates[field.Month] = make(map[string][]string)
			}

			if _, ok := templates[field.Month][name]; !ok {
				templates[field.Month][name] = make([]string, 0)
			}
			templates[field.Month][name] = append(templates[field.Month][name], field.Name)
			//记录模板的顺序
			if _, ok := templates[field.Month]["__order__"]; !ok {
				templates[field.Month]["__order__"] = make([]string, 0)
				templates[field.Month]["__order__"] = templateOrder
			}

		}
	}

	h.SendResponse(c, nil, TemplateFieldsResponse{
		Fields: templates,
	})
}
