package statistics

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
	"strconv"
)

var idNameMap map[uint64]string

var currentYear = ""

func DetailQuery(c *gin.Context) {
	log.Info("Detail function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r DetailRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
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
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	if len(salaries) < 1 {
		// 这里失败有一种场景，比如说2018年1月发的加班工资实际应该归纳于2017年12月，但是因为核算是在2018-01， 所以在tb_salary里并没有留有2017的记录，所以在这里会有查不到的问题。
		// 这种情况，我们还要试着去找下一年的数据 .
		year, _ := strconv.Atoi(r.Year)
		year = year + 1
		salaries, err = model.GetSalaryByAccountAndTemplate(strconv.Itoa(year), r.Account, ts)
		if err != nil {
			h.SendResponse(c, errno.ErrBind, err.Error())
			return
		}
	}
	sMap := make(map[uint64][]string)   //
	idNameMap = make(map[uint64]string) //记录用户ID 与名字的映射
	for _, item := range salaries {
		sMap[item.ID] = m[item.Template]
		idNameMap[item.ID] = item.Template
	}
	fields, err := model.GetFieldByKeys(r.Year, sMap)
	if err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	filename, err := writeIntoExcel(fields)
	if err != nil {
		h.SendResponse(c, errno.ErrWriteExcel, err.Error())
		return
	}
	model.CreateOperateRecord(c, fmt.Sprintf("查询员工收入情况"))

	h.SendResponse(c, nil, CreateResponse{File: filename})
	return

}

const (
	centerStyle = `"alignment":{"horizontal":"center","vertical":"center","ident":1}`
	borderStyle = `"border":[
						{"type":"left","color":"607D8B","style":1},
						{"type":"top","color":"607D8B","style":1},
						{"type":"bottom","color":"607D8B","style":1},
						{"type":"right","color":"607D8B","style":1}
					]`
	colorStyle = `"fill":{"type":"pattern","color":["#607D8B"],"pattern":13}}`
	valueColorStyle = `"fill":{"type":"pattern","color":["#B2DFDB"],"pattern":1}}`
	)

func writeIntoExcel(fields []model.SalaryField) (filename string, err error) {

	//获取所有用户档案 ，因为最终写入excel的是名字而不是ID
	allProfiles, err := model.GetAllProfile()
	allGroups, err := model.GetGroupWithAllChildren("部门")

	if err != nil {
		fmt.Println("GetAllProfileWidthGroup error :", err.Error())
	}
	profileIndexMap := make(map[uint64]int) // 记录 profile id 对应的 all profiles index
	departmentIndexMap := make(map[uint64]int)

	for i, p := range allProfiles {
		profileIndexMap[p.ID] = i
	}

	for i, g := range allGroups {
		departmentIndexMap[g.ID] = i
	}
	//对field 要进行排序归类等相关处理
	rows := make(map[uint64]map[string]map[string]map[string]float64) // 用户 -- > 模板名 --> 字段 --> 月份 -> 值
	departmentMap := make(map[uint64]map[string]uint64)               // 用户 -- > 月份 --> 部门
	monthOrder := []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}
	//
	for _, field := range fields {
		month := field.Month
		// 当fit_into_month 存在的时候，优先认fit_into_month
		//
		if field.FitIntoYear == currentYear {
			if len(field.FitIntoMonth) > 0 {
				month = field.FitIntoMonth
			}
		}
		// fitintoyear 有配置且不为当前年份，表现该笔数据是归属于其它年份的。忽略
		if field.FitIntoYear != currentYear && len(field.FitIntoYear) > 0 {
			field.Name = field.Name + "(归属于" + field.FitIntoYear + "-" + field.FitIntoMonth + ")"
		}

		if _, ok := rows[field.ProfileID]; !ok {
			rows[field.ProfileID] = make(map[string]map[string]map[string]float64)
		}

		if _, ok := departmentMap[field.ProfileID]; !ok {
			departmentMap[field.ProfileID] = make(map[string]uint64)
		}
		//记录用户在某个月属于哪个部门
		departmentMap[field.ProfileID][month] = field.DepartmentGroupID

		if _, ok := rows[field.ProfileID][idNameMap[field.SalaryID]]; !ok {
			rows[field.ProfileID][idNameMap[field.SalaryID]] = make(map[string]map[string]float64)
		}

		if _, ok := rows[field.ProfileID][idNameMap[field.SalaryID]][field.Name]; !ok {
			rows[field.ProfileID][idNameMap[field.SalaryID]][field.Name] = make(map[string]float64)
		}

		//这里要区分纳入的问题，比如说这个月是2018年7月，但加班补贴在2018年8月发。8月份发到手的这笔加班补贴实际应为7月的收入项目.
		rows[field.ProfileID][idNameMap[field.SalaryID]][field.Name][month] = field.Value
	}
	xlsx := excelize.NewFile()
	sheet1 := "员工明细"
	index := xlsx.NewSheet(sheet1)
	row := 4
	startCol := 3

	xlsx.SetCellValue(sheet1, "A3", "姓名")
	xlsx.SetCellValue(sheet1, "B3", "身份证")
	xlsx.SetCellValue(sheet1, "C3", "状态")

	style, _ := xlsx.NewStyle(`{`+centerStyle+`,`+borderStyle + `}` )
	err = xlsx.SetCellStyle(sheet1, "A1", "AZ1000", style)

	xlsx.SetColWidth(sheet1, "B", "B", 25)

	//fieldCellMap := make(map[string]string) //因为rows是无序化的map结构，所以需要用一个变量来记录起始cell位置
	templateOrder := []string{}
	fieldOrder := make(map[string][]string) // template => field[]
	for _, salary := range rows {
		for template, fields := range salary {
			templateOrder = append(templateOrder, template)
			for field := range fields {
				if _, ok := fieldOrder[template]; !ok {
					fieldOrder[template] = make([]string, 0)
				}
				fieldOrder[template] = append(fieldOrder[template], field)
			}
		}
		break //只需要遍历一行
	}

	for profile, salary := range rows {
		xlsx.MergeCell(sheet1, "A"+strconv.Itoa(row), "A"+strconv.Itoa(row+2))
		xlsx.SetCellValue(sheet1, "A"+strconv.Itoa(row), allProfiles[profileIndexMap[profile]].Name)
		xlsx.MergeCell(sheet1, "B"+strconv.Itoa(row), "B"+strconv.Itoa(row+2))
		xlsx.SetCellValue(sheet1, "B"+strconv.Itoa(row), allProfiles[profileIndexMap[profile]].IDCard)
		xlsx.MergeCell(sheet1, "C"+strconv.Itoa(row), "C"+strconv.Itoa(row+2))
		xlsx.SetCellValue(sheet1, "C"+strconv.Itoa(row), allProfiles[profileIndexMap[profile]].ID)
		templateCount := 0 // 每个字段会横向占用12个cell，使用这个变量来计算
		fieldCount := 0
		for _, template := range templateOrder {
			for _, field := range fieldOrder[template] {
				months := salary[template][field]
				//if _, ok := fieldCellMap[template+"."+field]; !ok {
				//	fieldCellMap[template+"."+field] = util.ConvertToNumberingScheme(fieldCount*13 + 1 + startCol)
				//}
				startCell := util.ConvertToNumberingScheme(fieldCount*13 + 1 + startCol)

				endCell := util.ConvertToNumberingScheme(fieldCount*13 + startCol + 13)
				xlsx.MergeCell(sheet1, startCell+"2", endCell+"2")
				xlsx.SetCellValue(sheet1, startCell+"2", template+"."+field)
				j := fieldCount*12 + startCol

				for month, value := range months {

					cell := fieldCount*13 + startCol + 1
					valueRow := row + 1
					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+1)+"3", "1月")
					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+2)+"3", "2月")
					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+3)+"3", "3月")
					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+4)+"3", "4月")
					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+5)+"3", "5月")
					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+6)+"3", "6月")
					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+7)+"3", "7月")
					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+8)+"3", "8月")
					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+9)+"3", "9月")
					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+10)+"3", "10月")
					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+11)+"3", "11月")
					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+12)+"3", "12月")
					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell)+strconv.Itoa(valueRow-1), "部门")
					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell)+strconv.Itoa(valueRow), "金额")


							//xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell)+strconv.Itoa(valueRow + 1), "")
					switch month {
					case "01":
						start := util.ConvertToNumberingScheme(cell+1)+strconv.Itoa(valueRow)
						xlsx.SetCellValue(sheet1, start, value)
						style, _ := xlsx.NewStyle(`{`+centerStyle+`,`+borderStyle+` ,`+ valueColorStyle )
						err = xlsx.SetCellStyle(sheet1, start, start, style)
					case "02":
						start := util.ConvertToNumberingScheme(cell+2)+strconv.Itoa(valueRow)
						xlsx.SetCellValue(sheet1, start, value)
						style, _ := xlsx.NewStyle(`{`+centerStyle+`,`+borderStyle+` ,`+ valueColorStyle )
						err = xlsx.SetCellStyle(sheet1, start, start, style)
					case "03":
						start := util.ConvertToNumberingScheme(cell+3)+strconv.Itoa(valueRow)
						xlsx.SetCellValue(sheet1, start, value)
						style, _ := xlsx.NewStyle(`{`+centerStyle+`,`+borderStyle+` ,`+ valueColorStyle )
						err = xlsx.SetCellStyle(sheet1, start, start, style)
					case "04":
						start := util.ConvertToNumberingScheme(cell+4)+strconv.Itoa(valueRow)
						xlsx.SetCellValue(sheet1, start, value)
						style, _ := xlsx.NewStyle(`{`+centerStyle+`,`+borderStyle+` ,`+ valueColorStyle )
						err = xlsx.SetCellStyle(sheet1, start, start, style)
					case "05":
						start := util.ConvertToNumberingScheme(cell+5)+strconv.Itoa(valueRow)
						xlsx.SetCellValue(sheet1, start, value)
						style, _ := xlsx.NewStyle(`{`+centerStyle+`,`+borderStyle+` ,`+ valueColorStyle )
						err = xlsx.SetCellStyle(sheet1, start, start, style)
					case "06":
						start := util.ConvertToNumberingScheme(cell+6)+strconv.Itoa(valueRow)
						xlsx.SetCellValue(sheet1, start, value)
						style, _ := xlsx.NewStyle(`{`+centerStyle+`,`+borderStyle+` ,`+ valueColorStyle )
						err = xlsx.SetCellStyle(sheet1, start, start, style)
					case "07":
						start := util.ConvertToNumberingScheme(cell+7)+strconv.Itoa(valueRow)
						xlsx.SetCellValue(sheet1, start, value)
						style, _ := xlsx.NewStyle(`{`+centerStyle+`,`+borderStyle+` ,`+ valueColorStyle )
						err = xlsx.SetCellStyle(sheet1, start, start, style)
					case "08":
						start := util.ConvertToNumberingScheme(cell+8)+strconv.Itoa(valueRow)
						xlsx.SetCellValue(sheet1, start, value)
						style, _ := xlsx.NewStyle(`{`+centerStyle+`,`+borderStyle+` ,`+ valueColorStyle )
						err = xlsx.SetCellStyle(sheet1, start, start, style)
					case "09":
						start := util.ConvertToNumberingScheme(cell+9)+strconv.Itoa(valueRow)
						xlsx.SetCellValue(sheet1, start, value)
						style, _ := xlsx.NewStyle(`{`+centerStyle+`,`+borderStyle+` ,`+ valueColorStyle )
						err = xlsx.SetCellStyle(sheet1, start, start, style)
					case "10":
						start := util.ConvertToNumberingScheme(cell+10)+strconv.Itoa(valueRow)
						xlsx.SetCellValue(sheet1, start, value)
						style, _ := xlsx.NewStyle(`{`+centerStyle+`,`+borderStyle+` ,`+ valueColorStyle )
						err = xlsx.SetCellStyle(sheet1, start, start, style)
					case "11":
						start := util.ConvertToNumberingScheme(cell+11)+strconv.Itoa(valueRow)
						xlsx.SetCellValue(sheet1, start, value)
						style, _ := xlsx.NewStyle(`{`+centerStyle+`,`+borderStyle+` ,`+ valueColorStyle )
						err = xlsx.SetCellStyle(sheet1, start, start, style)
					case "12":
						start := util.ConvertToNumberingScheme(cell+12)+strconv.Itoa(valueRow)
						xlsx.SetCellValue(sheet1, start, value)
						style, _ := xlsx.NewStyle(`{`+centerStyle+`,`+borderStyle+` ,`+ valueColorStyle )
						err = xlsx.SetCellStyle(sheet1, start, start, style)
					}
					j++
				}
				//对部门进行分阶段，然后在excel中写入到已合并的表格中
				// j := 1
				//departLen := 1

				// 分析出部门所占的Cell数目
				//departs := make([]uint64 , 12)
				departArr := make([]uint64, 0)
				lenArr := make([]int , 0)

				for _, d := range monthOrder {
					departArr = append(departArr, departmentMap[profile][d])
				}

				// departArr 的格式 ： [69 69 70 70 70 70 70 70 70 70 70 69]
				// 现在要得出长度： [2,9,1]
				last := departArr[0]
				lenArr = append(lenArr,1)

				newDepartArr := make([]uint64, 0)
				newDepartArr = append(newDepartArr, departArr[0])
				for i, d := range departArr[1:]{
					if d != last {  //
						lenArr = append(lenArr,1)
						newDepartArr = append(newDepartArr, d)
					}else {
						lenArr[len(lenArr) - 1] += 1
					}

					last = departArr[i+1]
				}
				from := 2
				for i, d := range newDepartArr {
					start := util.ConvertToNumberingScheme(fieldCount*13 + startCol + from)
					end := start
					if lenArr[i] > 1 {
						end = util.ConvertToNumberingScheme(fieldCount*13 + startCol + from + lenArr[i] -1  )
						xlsx.MergeCell(sheet1, start+strconv.Itoa(row), end+strconv.Itoa(row))
					}

					xlsx.SetCellValue(sheet1, start+strconv.Itoa(row), allGroups[departmentIndexMap[d]].Name)
					style, _ := xlsx.NewStyle(`{`+centerStyle+`,`+borderStyle+` ,`+ colorStyle )
					err = xlsx.SetCellStyle(sheet1, start+strconv.Itoa(row), end +strconv.Itoa(row), style)
					from = from + lenArr[i]
				}

				fieldCount++
			}
			//fieldCount++
		}
		templateCount++
		row += 3
	}

	//style, _ := xlsx.NewStyle(`{"custom_number_format": "[$-380A]dddd\\,\\ dd\" de \"mmmm\" de \"yyyy;@"}`)
	//xlsx.SetCellStyle(sheet1, "A1", util.ConvertToNumberingScheme(fieldc) + strconv.Itoa(row+1), style)

	xlsx.SetActiveSheet(index)
	// Save xlsx file by the given path.
	filename = "detail.xlsx"
	err = xlsx.SaveAs("./export/" + filename)
	if err != nil {
		fmt.Println(err)
	}
	return filename, nil
}
