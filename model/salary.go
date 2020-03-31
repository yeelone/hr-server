package model

import (
	"fmt"
	"hr-server/util"
	"strings"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

type Salary struct {
	ID              uint64         `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	CreatedAt       time.Time      `gorm:"column:createdAt" json:"-"`
	TemplateAccount uint64         `json:"template_account" gorm:"not null" binding:"required"`
	Template        string         `json:"template" gorm:"not null" binding:"required"`
	Year            string         `json:"year" gorm:"not null" binding:"required"`
	Month           string         `json:"month" gorm:"not null" binding:"required"`
	UserID          User           `json:"user_id"`
	Fields          []SalaryField  `json:"fields"`
	Locked          bool           `json:"locked"` //往月工资核算完毕将进行锁定，不可再修改
	Data            postgres.Jsonb //修改为加入fields,废除此项
}

var SALARYTABLENAME = "tb_salary"

// TableName :
func (w *Salary) TableName() string {
	return SALARYTABLENAME
}

// Create creates a new salary .
func (w *Salary) Create() error {

	//在create的时候，gorm 会进行一条一条的插入到tb_salary_fields表中，效率十分低下。所以在这里要优化成批量插入.
	// 批量插入之后 ，从原先的41s 缩短为9s .
	s := Salary{}
	s.Month = w.Month
	s.Year = w.Year
	s.Template = w.Template
	s.TemplateAccount = w.TemplateAccount
	s.UserID = w.UserID
	s.Locked = w.Locked
	if err := DB.Self.Create(&s).Error; err != nil {
		return err
	}

	if err := BatchCreate(s.ID, w.Fields); err != nil {
		return err
	}

	return nil
}

func ClearSalary(year, month string, templateAccount uint64) (err error) {
	//对于月度工资来说只能计算并保存一次，所以在每次创建当月工资前，先清空当月的表中数据 。避免用户多次计算导致的数据冲突
	//如果模板已经存在，则删除其中关联的fields，然后再重新添加关联
	sModel := []Salary{}
	if err = DB.Self.Where("template_account = ? AND year = ? AND month = ? ", templateAccount, year, month).Find(&sModel).Error; err != nil {
		return err
	}

	if len(sModel) > 0 {
		ids := make([]uint64, len(sModel))
		for i, s := range sModel {
			ids[i] = s.ID
		}
		if err := DeleteSalaryFieldsByMonthAndTemplate(ids); err != nil {
			return err
		}

		return DB.Self.Where("id IN (?)", ids).Delete(&sModel).Error
	} else {
		fmt.Println("clear error", err)
	}

	return nil
}

func GetRelatedTemplateValue(year, month, template string, templateAccountID uint64, fields []string) (result []SalaryField) {
	sModel := &Salary{}
	if err := DB.Self.Model(&sModel).Where("template_account = ? AND template = ? AND  year = ? AND month = ?", templateAccountID, template, year, month).First(&sModel).Error; err != nil {
		log.Info("GetRelatedTemplateValue function called.", lager.Data{"message": "无法找到相关的工资模板:" + template, "error": err.Error()})
	}

	newFieldStr := []string{}
	for _, field := range fields {
		newFieldStr = append(newFieldStr, "'"+field+"'")
	}

	where := ""
	if len(newFieldStr) > 0 {
		where = "key IN (" + strings.Join(newFieldStr, ",") + ")  and "
	}

	where = where + " salary_id = " + util.Uint2Str(sModel.ID)
	fieldList := &SalaryField{}
	if err := DB.Self.Model(fieldList).Where(where).Find(&result).Error; err != nil {
		fmt.Println("error is :" + err.Error())
	}
	return result
}

func GetSalaryByAccountAndTemplate(year string, account uint64, templates []string) (result []Salary, err error) {
	if err = DB.Self.Select("id,template").Where("year = ? AND template_account = ? AND template IN  (?) ", year, account, templates).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func GetSalaryByAccount(year string, account uint64) (result []Salary, err error) {
	if err = DB.Self.Select("id,template,month").Where("year = ? AND template_account = ?", year, account).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

//select department_group_id,month ,count(profile_id),sum(value) from tb_salary_fields where  key='实发金额' group by department_group_id,month ;
func GetSalaryByDepartment(year string, account uint64, templates []string) (result []Salary, err error) {
	if err = DB.Self.Select("id,template").Where("year = ? AND template_account = ? AND template IN  (?) ", year, account, templates).Find(&result).Error; err != nil {
		fmt.Println("err ", err)
		return nil, err
	}

	return result, nil
}

func GetSalaryByTemplateAccount(account uint64, template, year, month string) (result []Salary, err error) {
	if err = DB.Self.Select("id,template,year,month").Where("template_account = ? AND template = ? AND month=? AND year=?  ", account, template, year, month).Find(&result).Error; err != nil {
		fmt.Println("err ", err)
		return nil, err
	}

	return result, nil
}

func GetSalary(id uint64) (result Salary, err error) {
	if err = DB.Self.Where("id = ? ", id).First(&result).Error; err != nil {
		fmt.Println("err ", err)
		return result, err
	}

	return result, nil
}
