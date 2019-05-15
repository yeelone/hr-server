package statistics

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/errno"
	"hrgdrc/util"
	"strconv"
)

var departmentNameMap map[uint64]string
func DepartmentIncomeQuery(c *gin.Context) {
	log.Info("DepartmentIncomeQuery function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r DetailRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind,  err.Error())
		return
	}
	currentYear = r.Year
	m := make(map[string][]string)
	for _, t := range r.Templates {
		m[t.Template] = t.Fields
	}
	ts := []string{}
	for k := range m {
		ts = append(ts, k)
	}
	salaries, err := model.GetSalaryByAccountAndTemplate(r.Year, r.Account, ts)
	if err != nil {
		fmt.Println("GetSalaryByAccountAndTemplate error :", err.Error())
		h.SendResponse(c, errno.ErrBind,  err.Error())
		return
	}

	if len(salaries) < 1 {
		//这里失败有一种场景，比如说2018年1月发的加班工资实际应该归纳于2017年12月，但是因为核算是在2018-01， 所以在tb_salary里并没有留有2017的记录，所以在这里会有查不到的问题。
		// 这种情况，我们还要试着去找下一年的数据 .
		year, _ := strconv.Atoi(r.Year)
		year = year + 1
		salaries, err = model.GetSalaryByAccountAndTemplate(strconv.Itoa(year), r.Account, ts)
		if err != nil {
			h.SendResponse(c, errno.ErrBind,  err.Error())
			return
		}
	}

	sMap := make(map[uint64][]string)
	departmentNameMap = make(map[uint64]string)
	for _, item := range salaries {
		sMap[item.ID] = m[item.Template]
		departmentNameMap[item.ID] = item.Template
	}

	result, err := model.GetDepartmentTotalIncome(r.Year, sMap)
	if err != nil {
		h.SendResponse(c, errno.ErrDatabase,  err.Error())
		return
	}

	filename, err := writeDepartmentIntoExcel(result)
	if err != nil {
		h.SendResponse(c, errno.ErrWriteExcel,  err.Error())
		return
	}
	h.SendResponse(c, nil, CreateResponse{File: filename})
	return

}

func writeDepartmentIntoExcel(data []model.Statistics) (filename string, err error) {
	//获取所有用户档案 ，因为最终写入excel的是名字而不是ID
	allGroups, err := model.GetGroupWithAllChildren("部门")
	if err != nil {
		fmt.Println("GetAllProfileWidthGroup error :", err.Error())
	}
	departmentIndexMap := make(map[uint64]int)

	for i, g := range allGroups {
		departmentIndexMap[g.ID] = i
	}
	sheet1 := "Sheet1"
	monthOrder := make(map[string]int)
	monthOrder["01"] = 1
	monthOrder["02"] = 2
	monthOrder["03"] = 3
	monthOrder["04"] = 4
	monthOrder["05"] = 5
	monthOrder["06"] = 6
	monthOrder["07"] = 7
	monthOrder["08"] = 8
	monthOrder["09"] = 9
	monthOrder["10"] = 10
	monthOrder["11"] = 11
	monthOrder["12"] = 12

	xlsx, err := excelize.OpenFile("./export/template/department_detail_template.xlsx")
	//index := xlsx.NewSheet(sheet1)
	if err != nil {
		fmt.Println("OpenFile", filename, err)
		return "", err
	}

	startRow := 5
	startCol := 3
	departmentRowMap := make(map[string]int) //记录部门所在的行
	for i, d := range data {
		if d.Year != currentYear {
			continue
		}
		row := startRow + i
		if _, ok := departmentRowMap[d.Department]; !ok {
			departmentRowMap[d.Department] = row
			xlsx.SetCellValue(sheet1, "A"+strconv.Itoa(row), i+1)
			department, _ := strconv.ParseUint(d.Department, 10, 64)
			xlsx.SetCellValue(sheet1, "B"+strconv.Itoa(row), allGroups[departmentIndexMap[department]].Name)
			xlsx.SetCellValue(sheet1, "C"+strconv.Itoa(row), allGroups[departmentIndexMap[department]].Code)
		}

		col := (monthOrder[d.Month]-1)*2 + startCol
		xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(col+1)+strconv.Itoa(departmentRowMap[d.Department]), d.Number)
		xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(col+2)+strconv.Itoa(departmentRowMap[d.Department]), d.Total)
	}
	// Save xlsx file by the given path.
	filename = "department.xlsx"
	err = xlsx.SaveAs("./export/" + filename)
	if err != nil {
		fmt.Println(err)
	}
	return filename, nil
}
