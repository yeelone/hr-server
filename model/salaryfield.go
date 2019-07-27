package model

import (
	"fmt"
	"hr-server/util"
	"strconv"
	"strings"
	"time"
)

type SalaryField struct {
	ID                uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	CreatedAt         time.Time `gorm:"column:createdAt" json:"-"`
	ProfileID         uint64
	SalaryID          uint64
	DepartmentGroupID uint64  `json:"department_group_id" ` //记录当时所在的部门
	PostGroupID       uint64  `json:"post_group_id"`        //记录当时所在的岗位
	Key               string  `json:"key"`                  //指定标识
	Name              string  `json:"name"`                 //字段名字
	Alias             string  `json:"alias"`                //字段别名，即显示名
	Value             float64 `json:"value"`                //值
	Content           string  `json:"content"`              //内容，比如key为name等数据时，content可以为张三，李四，而此时value为0
	ShouldTax         bool    `json:"should_tax"`           //是否纳税
	IsIncome          bool    `json:"is_income"`            //收入项
	IsDeduct          bool    `json:"is_deduct"`            //扣除项
	Year              string  `json:"year"`
	Month             string  `json:"month"`
	FitIntoYear       string  `json:"fit_into_year"`  //纳入指定月份
	FitIntoMonth      string  `json:"fit_into_month"` //纳入指定月份
}

var SALARYFIELDTABLENAME = "tb_salary_fields"

// TableName :
func (s *SalaryField) TableName() string {
	return SALARYFIELDTABLENAME
}

// Create creates a new salary field.
func (s *SalaryField) Create() error {
	return DB.Self.Create(&s).Error
}

// 批量插入
func BatchCreate(salary_id uint64, sf []SalaryField) error {
	keys := "salary_id,profile_id,department_group_id,post_group_id,key,name, alias,value,content,should_tax,is_income,is_deduct,year,month,fit_into_year,fit_into_month"
	sql := "insert into " + SALARYFIELDTABLENAME + "(" + keys + ") values"
	values := []string{}
	for _, f := range sf {
		s := `'` + util.Uint2Str(salary_id) + `','` +
			util.Uint2Str(f.ProfileID) + `','` +
			util.Uint2Str(f.DepartmentGroupID) + `','` +
			util.Uint2Str(f.PostGroupID) + `','` +
			f.Key + `','` +
			f.Name + `','` +
			f.Alias + `','` +
			fmt.Sprint(f.Value) + `','` +
			f.Content + `','` +
			fmt.Sprint(f.ShouldTax) + `','` +
			fmt.Sprint(f.IsIncome) + `','` +
			fmt.Sprint(f.IsDeduct) + `','` +
			f.Year + `','` +
			f.Month + `','` +
			f.FitIntoYear + `','` +
			f.FitIntoMonth + `'`
		values = append(values, `(`+s+`)`)
	}

	sql += strings.Join(values, ",")
	return DB.Self.Exec(sql).Error
}

func DeleteSalaryFieldsByMonthAndTemplate(salaryIDs []uint64) (err error) {
	tx := DB.Self.Begin()
	err = tx.Where("salary_id IN (?) ", salaryIDs).Delete(SalaryField{}).Error
	tx.Commit()
	return
}

// GetFieldByKeys ： //获取职工年底收入明细，这里需要注意，应该把 fit_into_month fit_into_year 考虑进来
// @param:
//       year : 指定年份
//       salaryMap : map[uint64][]string , salary id ,与所指定的field key
//  @return :
//   result : []SalaryField
// 通过账套 和模板取得 salary_id ,  再结合salary_id 和 字段名取得相应的fields
func GetFieldByKeys(year string, salarymap map[uint64][]string) (result []SalaryField, err error) {
	ids := []uint64{}
	whereSql := []string{} // sql 构造如下 (salary_id=? and key IN (?)) or (salary_id=? and key IN (?))
	for id, keys := range salarymap {
		ids = append(ids, id)

		// 要手工给key 加 双引号
		newKeys := []string{}
		for _, t := range keys {
			newKeys = append(newKeys, `'`+t+`'`)
		}
		s := "(salary_id=" + util.Uint2Str(id) + " AND key IN (" + strings.Join(newKeys, ",") + "))"
		whereSql = append(whereSql, s)
	}
	if err = DB.Self.Debug().Select("profile_id,department_group_id,post_group_id,salary_id,key,name, month,fit_into_year, fit_into_month,value").Where("(year='" + year + "' OR fit_into_year='" + year + "') AND " + strings.Join(whereSql, " or ")).Order("profile_id,key, month").Find(&result).Error; err != nil {
		fmt.Println("err ", err)
		return nil, err
	}

	return result, nil
}

// GetFieldByKeyAndMonth ： //获取指定两个月份之间的key value
// @param:
//       key :
//       fromMonth : 格式比如 9
//		 toMonth :
//  @return :
//   result : []SalaryField
// 通过账套 和模板取得 salary_id ,  再结合salary_id 和 字段名取得相应的fields
func GetFieldByKeyAndMonth(key, year, fromMonth, toMonth string) (result []SalaryField, err error) {

	month1, _ := strconv.Atoi(fromMonth)
	month2, _ := strconv.Atoi(toMonth)
	specialMonth := make([]string, 0)
	for {
		if month1+1 < month2+2 {
			m := strconv.Itoa(month1)
			if len(m) < 2 {
				m = "0" + m
			}
			specialMonth = append(specialMonth, m)
			month1 += 1
		} else {
			break
		}
	}
	if err = DB.Self.Debug().Select("profile_id,key,value,year,month").Where("key=? AND year=? AND month IN (?)", key, year, specialMonth).Order("profile_id,key, month").Find(&result).Error; err != nil {
		fmt.Println("err ", err)
		return nil, err
	}

	return result, nil
}

// GetFieldLastMonthValueByKey ： //获取指定月份的key value
// @param:
//       key :
//       year : 格式比如 9
//		 month :
//  @return :
//   result : []SalaryField
// 通过账套 和模板取得 salary_id ,  再结合salary_id 和 字段名取得相应的fields
func GetFieldLastMonthValueByKey(key, year, month string) (result []SalaryField, err error) {
	if len(month) == 1 {
		month = "0" + month
	}

	if err = DB.Self.Debug().Select("profile_id,key,value,year,month").Where("key=? AND year=? AND month=?", key, year, month).Order("profile_id,key, month").Find(&result).Error; err != nil {
		fmt.Println("err ", err)
		return nil, err
	}

	return result, nil
}

// GetDepartmentTotalIncome ：
// @param:
//       year : 指定年份
//       salaryMap : map[uint64][]string , salary id ,与所指定的field key
//  @return :
//   result : []SalaryField
// 通过账套 和模板取得 salary_id ,  再结合salary_id 和 字段名取得相应的fields
func GetDepartmentTotalIncome(year string, salarymap map[uint64][]string) (result []Statistics, err error) {
	ids := []uint64{}
	whereSql := []string{} // sql 构造如下 (salary_id=? and key IN (?)) or (salary_id=? and key IN (?))
	for id, keys := range salarymap {
		ids = append(ids, id)

		// 要手工给key 加 双引号
		newKeys := []string{}
		for _, t := range keys {
			newKeys = append(newKeys, `'`+t+`'`)
		}
		s := "(salary_id=" + util.Uint2Str(id) + " AND key IN (" + strings.Join(newKeys, ",") + "))"
		whereSql = append(whereSql, s)
	}

	if err = DB.Self.Table(SALARYFIELDTABLENAME).Debug().Select("department_group_id as department,fit_into_year as year ,fit_into_month as month,count(profile_id) as number,sum(value) as total").Where("fit_into_year='" + year + "' AND " + strings.Join(whereSql, " or ")).Group("department_group_id,fit_into_year,fit_into_month").Order("fit_into_month").Scan(&result).Error; err != nil {
		fmt.Println("err ", err)
		return nil, err
	}

	return result, nil
}

func GetSalaryFieldsByMonthAndTemplate(salaryIDs []uint64) (fields []SalaryField, err error) {
	err = DB.Self.Where("salary_id IN (?) ", salaryIDs).Find(fields).Error
	return fields, err
}

func GetSalaryFieldByDepartment(groups []uint64, beginDate, endDate string) (fields []SalaryField, err error) {
	// sModel := &Salary{}
	// if err := DB.Self.Model(&sModel).Where("department_group_id IN （?） AND  year = ? AND month = ?", groups, year, month).First(&sModel).Error; err != nil {
	// 	fmt.Println("无法找到相关的工资模板")
	// }

	// newFieldStr := []string{}
	// for _, field := range fields {
	// 	newFieldStr = append(newFieldStr, "'"+field+"'")
	// }

	// where := ""
	// if len(newFieldStr) > 0 {
	// 	where = "key IN (" + strings.Join(newFieldStr, ",") + ")  and "
	// }

	// where = where + " salary_id = " + util.Uint2Str(sModel.ID)
	// fieldList := &SalaryField{}
	// if err := DB.Self.Model(fieldList).Where(where).Find(&result).Error; err != nil {
	// 	fmt.Println("error is :" + err.Error())
	// }
	return fields, err
}

func GetSalaryFieldByProfileAndMonth(year string, month string, profileID uint64) (result []SalaryField, err error) {
	if err = DB.Self.Where("year = ? AND month = ? AND profile_id=? ", year, month, profileID).Find(&result).Error; err != nil {
		fmt.Println("err ", err)
		return nil, err
	}

	return result, nil
}

func GetFieldsBySalaryAndProfilesAndYear(year string, salary []uint64, profiles []uint64) (result []SalaryField, err error) {
	if err = DB.Self.Debug().Where("year = ? AND salary_id IN (?) AND profile_id IN  (?) ", year, salary, profiles).Find(&result).Error; err != nil {
		fmt.Println("err ", err)
		return nil, err
	}

	return result, nil
}
