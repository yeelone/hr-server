package salary

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/errno"
	"hrgdrc/pkg/template"
	"hrgdrc/util"
	"sort"
	"strconv"
)

// GetProfileMonthSalary : 获取指定员工指定月份的工资明细
func GetProfileMonthSalary(c *gin.Context) {
	log.Info("GetProfileMonthSalary query function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r ProfileSalaryRequest
	if err := c.Bind(&r); err != nil {
		fmt.Println(err)
		h.SendResponse(c, errno.ErrBind, err)
		return
	}
	pid, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	year := c.Param("year")
	month := c.Param("month")

	//如果用户是查询岗，则用户只能查询自己
	userid, ok := c.Get("userid")
	if !ok {
		h.SendResponse(c, errno.StatusUnauthorized, nil)
		return
	}

	user, err := model.GetUser(userid.(uint64))

	profile, err := model.GetProfileByIDCard(user.IDCard)
	if err != nil {
		h.SendResponse(c, errno.StatusUnauthorized, err)
		return
	}

	for _, role := range user.Roles {
		if role.Name == "查询岗" {
			if profile.ID != pid {
				h.SendResponse(c, errno.StatusUnauthorized, "您只有查询权限，只能访问你自己的档案")
				return
			}
		}
	}

	// 获取用户所有salary field
	fields, err := model.GetSalaryFieldByProfileAndMonth(year, month, pid)
	if err != nil {
		h.SendResponse(c, errno.ErrSalaryProfileDetail, err)
		return
	}

	if len(fields) < 1 {
		h.SendResponse(c, nil, ProfileSalaryResponse{})
		return
	}

	// 随便取出一个field ，得到 DepartmentGroupID and PostGroupID
	departmentID := fields[0].DepartmentGroupID
	postID := fields[0].PostGroupID

	departGroup, _ := model.GetGroup(departmentID, false)
	postGroup, _ := model.GetGroup(postID, false)
	//查看员工在哪个账套里
	//根据salary id 来获取账套
	salary, err := model.GetSalary(fields[0].SalaryID)
	if err != nil {
		h.SendResponse(c, errno.ErrSalaryProfileDetail, err)
		return
	}
	//获取账套，为的是取得账套内模板的order，方便前端显示时更有逻辑
	templateAccount, err := model.GetTemplateAccount(salary.TemplateAccount)
	if err != nil {
		h.SendResponse(c, errno.ErrGetTemplateAccount, err)
		return
	}

	fieldMap := make(map[string]float64)
	for _, f := range fields {
		fieldMap[f.Key] = f.Value
	}
	// 之所以以slice的形式了，是因为要使之变成有序列表
	templates := make([]map[string][]template.Field, len(templateAccount.Templates))
	orderMap := make(map[uint64]int)
	for i, id := range templateAccount.Order {
		orderMap[uint64(id)] = i
	}
	for _, t := range templateAccount.Templates {
		tp, err := template.ResolveTemplate(t.Name)

		m := make(map[string][]template.Field)

		if err != nil {
			h.SendResponse(c, errno.ErrTemplateInvalid, nil)
			return
		}

		fields := make([]template.Field, 0)

		orders := []int{}
		fieldTemp := make(map[int]template.Field)
		//按顺序排列
		for _, field := range tp.All {
			if _, ok := fieldMap[field.Key]; ok {
				field.Value = fieldMap[field.Key]
				fieldTemp[field.Order] = field
				orders = append(orders, field.Order)
			}
		}
		sort.Ints(orders)
		for _, index := range orders {
			fields = append(fields, fieldTemp[index])
		}

		m[t.Name] = fields

		templates[orderMap[t.ID]] = m
	}

	h.SendResponse(c, nil, ProfileSalaryResponse{
		TemplateList: templates,
		Profile:      profile,
		Department:   departGroup.Name,
		Post:         postGroup.Name,
	})
}
