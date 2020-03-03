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

//
//func writeIntoExcel(fields []model.SalaryField) (filename string, err error) {
//
//	//获取所有用户档案 ，因为最终写入excel的是名字而不是ID
//	allProfiles, err := model.GetAllProfile()
//	allGroups, err := model.GetGroupWithAllChildren("部门")
//	if err != nil {
//		fmt.Println("GetAllProfileWidthGroup error :", err.Error())
//	}
//	profileIndexMap := make(map[uint64]int) // 记录 profile id 对应的 all profiles index
//	departmentIndexMap := make(map[uint64]int)
//
//	for i, p := range allProfiles {
//		profileIndexMap[p.ID] = i
//	}
//
//	for i, g := range allGroups {
//		departmentIndexMap[g.ID] = i
//	}
//	//对field 要进行排序归类等相关处理
//	rows := make(map[uint64]map[string]map[string]map[string]float64) // 用户 -- > 模板名 --> 字段 --> 月份 -> 值
//	departmentMap := make(map[uint64]map[string]uint64)               // 用户 -- > 月份 --> 部门
//	monthOrder := []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}
//	for _, field := range fields {
//		month := field.Month
//		// 当fit_into_month 存在的时候，优先认fit_into_month
//		//
//		if field.FitIntoYear == currentYear {
//			if len(field.FitIntoMonth) > 0 {
//				month = field.FitIntoMonth
//			}
//		}
//		// fitintoyear 有配置且不为当前年份，表现该笔数据是归属于其它年份的。忽略
//		if field.FitIntoYear != currentYear && len(field.FitIntoYear) > 0 {
//			field.Name = field.Name + "(归属于" + field.FitIntoYear + "-" + field.FitIntoMonth + ")"
//		}
//
//		if _, ok := rows[field.ProfileID]; !ok {
//			rows[field.ProfileID] = make(map[string]map[string]map[string]float64)
//		}
//
//		if _, ok := departmentMap[field.ProfileID]; !ok {
//			departmentMap[field.ProfileID] = make(map[string]uint64)
//		}
//		//记录用户在某个月属于哪个部门
//		departmentMap[field.ProfileID][month] = field.DepartmentGroupID
//
//		if _, ok := rows[field.ProfileID][idNameMap[field.SalaryID]]; !ok {
//			rows[field.ProfileID][idNameMap[field.SalaryID]] = make(map[string]map[string]float64)
//		}
//
//		if _, ok := rows[field.ProfileID][idNameMap[field.SalaryID]][field.Name]; !ok {
//			rows[field.ProfileID][idNameMap[field.SalaryID]][field.Name] = make(map[string]float64)
//		}
//
//		//这里要区分纳入的问题，比如说这个月是2018年7月，但加班补贴在2018年8月发。8月份发到手的这笔加班补贴实际应为7月的收入项目.
//		rows[field.ProfileID][idNameMap[field.SalaryID]][field.Name][month] = field.Value
//	}
//	fmt.Println(util.PrettyJson(rows))
//	xlsx := excelize.NewFile()
//	sheet1 := "员工明细"
//	index := xlsx.NewSheet(sheet1)
//	row := 4
//	startCol := 3
//
//	xlsx.SetCellValue(sheet1, "A3", "姓名")
//	xlsx.SetCellValue(sheet1, "B3", "身份证")
//	xlsx.SetCellValue(sheet1, "C3", "状态")
//	for profile, salary := range rows {
//		xlsx.MergeCell(sheet1, "A"+strconv.Itoa(row), "A"+strconv.Itoa(row+2))
//		xlsx.SetCellValue(sheet1, "A"+strconv.Itoa(row), allProfiles[profileIndexMap[profile]].Name)
//		xlsx.MergeCell(sheet1, "B"+strconv.Itoa(row), "B"+strconv.Itoa(row+2))
//		xlsx.SetCellValue(sheet1, "B"+strconv.Itoa(row), allProfiles[profileIndexMap[profile]].IDCard)
//		xlsx.MergeCell(sheet1, "C"+strconv.Itoa(row), "C"+strconv.Itoa(row+2))
//		xlsx.SetCellValue(sheet1, "C"+strconv.Itoa(row), allProfiles[profileIndexMap[profile]].Status)
//		templateCount := 0 // 每个字段会横向占用12个cell，使用这个变量来计算
//		fieldCount := 0
//		for template, fields := range salary {
//			for field, months := range fields {
//				startCell := util.ConvertToNumberingScheme(fieldCount*13 + 1 + startCol)
//				endCell := util.ConvertToNumberingScheme(fieldCount*13 + startCol + 13)
//				xlsx.MergeCell(sheet1, startCell+"2", endCell+"2")
//				xlsx.SetCellValue(sheet1, startCell+"2", template+"."+field)
//				j := fieldCount*12 + startCol
//				fmt.Println("field", field)
//				for month, value := range months {
//					fmt.Println("month ,value ", month, value)
//					cell := fieldCount*13 + startCol + 1
//					valueRow := row + 1
//					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+1)+"3", "1月")
//					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+2)+"3", "2月")
//					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+3)+"3", "3月")
//					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+4)+"3", "4月")
//					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+5)+"3", "5月")
//					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+6)+"3", "6月")
//					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+7)+"3", "7月")
//					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+8)+"3", "8月")
//					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+9)+"3", "9月")
//					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+10)+"3", "10月")
//					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+11)+"3", "11月")
//					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+12)+"3", "12月")
//					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell)+strconv.Itoa(valueRow-1), "部门")
//					xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell)+strconv.Itoa(valueRow), "金额")
//					//xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell)+strconv.Itoa(valueRow + 1), "")
//					switch month {
//					case "01":
//						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+1)+strconv.Itoa(valueRow), value)
//					case "02":
//						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+2)+strconv.Itoa(valueRow), value)
//					case "03":
//						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+3)+strconv.Itoa(valueRow), value)
//					case "04":
//						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+4)+strconv.Itoa(valueRow), value)
//					case "05":
//						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+5)+strconv.Itoa(valueRow), value)
//					case "06":
//						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+6)+strconv.Itoa(valueRow), value)
//					case "07":
//						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+7)+strconv.Itoa(valueRow), value)
//					case "08":
//						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+8)+strconv.Itoa(valueRow), value)
//					case "09":
//						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+9)+strconv.Itoa(valueRow), value)
//					case "10":
//						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+10)+strconv.Itoa(valueRow), value)
//					case "11":
//						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+11)+strconv.Itoa(valueRow), value)
//					case "12":
//						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+12)+strconv.Itoa(valueRow), value)
//					}
//					j++
//				}
//				//对部门进行分阶段，然后在excel中写入到已合并的表格中
//				//j := 1
//				//departLen := 1
//
//				//分析出部门所占的Cell数目
//				//departs := make([]uint64 , 12)
//				departArr := make([]uint64, 0)
//				lenArr := make([]int, 0)
//
//				length := 1
//				last := departmentMap[profile]["01"]
//				for _, d := range monthOrder[1:] {
//					if departmentMap[profile][d] != last { // 这个月跟上个月不同，表示调整了部门
//						lenArr = append(lenArr, length) // 记录上个部门所占有的月数
//						departArr = append(departArr, last)
//						last = departmentMap[profile][d]
//						length = 1
//					} else {
//						length++
//					}
//
//					if d == "12" { //当是最后一个月的时候
//						lenArr = append(lenArr, length)
//						departArr = append(departArr, last)
//					}
//				}
//
//				from := 1
//				for i, d := range departArr {
//					if d == 0 {
//						from += lenArr[i] + 1
//						continue
//					}
//					start := util.ConvertToNumberingScheme(fieldCount*13 + startCol + from)
//					if lenArr[i] > 1 {
//						end := util.ConvertToNumberingScheme(fieldCount*13 + startCol + from + lenArr[i] - 1)
//						xlsx.MergeCell(sheet1, start+strconv.Itoa(row), end+strconv.Itoa(row))
//					}
//
//					xlsx.SetCellValue(sheet1, start+strconv.Itoa(row), allGroups[departmentIndexMap[d]].Name)
//					from = from + lenArr[i]
//				}
//
//				fieldCount++
//			}
//			fieldCount++
//		}
//		templateCount++
//		row += 3
//	}
//
//	//style, _ := xlsx.NewStyle(`{"custom_number_format": "[$-380A]dddd\\,\\ dd\" de \"mmmm\" de \"yyyy;@"}`)
//	//xlsx.SetCellStyle(sheet1, "A1", util.ConvertToNumberingScheme(fieldc) + strconv.Itoa(row+1), style)
//
//	xlsx.SetActiveSheet(index)
//	// Save xlsx file by the given path.
//	filename = "detail.xlsx"
//	err = xlsx.SaveAs("./export/" + filename)
//	if err != nil {
//		fmt.Println(err)
//	}
//	return filename, nil
//}

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
	//monthOrder := []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}
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

	fmt.Println("templateOrder", util.PrettyJson(templateOrder))
	fmt.Println("fieldOrder", util.PrettyJson(fieldOrder))

	for profile, salary := range rows {
		fmt.Println("profile", profile)
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

				fmt.Println("template,field", startCell, template, field)

				for month, value := range months {

					cell := fieldCount*13 + startCol + 1
					valueRow := row + 1
					if field == "独生子女费" && value > 30 {
						fmt.Println("month ,value ", profile, field, month, value, cell)
					}
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
						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+1)+strconv.Itoa(valueRow), value)
					case "02":
						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+2)+strconv.Itoa(valueRow), value)
					case "03":
						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+3)+strconv.Itoa(valueRow), value)
					case "04":
						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+4)+strconv.Itoa(valueRow), value)
					case "05":
						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+5)+strconv.Itoa(valueRow), value)
					case "06":
						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+6)+strconv.Itoa(valueRow), value)
					case "07":
						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+7)+strconv.Itoa(valueRow), value)
					case "08":
						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+8)+strconv.Itoa(valueRow), value)
					case "09":
						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+9)+strconv.Itoa(valueRow), value)
					case "10":
						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+10)+strconv.Itoa(valueRow), value)
					case "11":
						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+11)+strconv.Itoa(valueRow), value)
					case "12":
						xlsx.SetCellValue(sheet1, util.ConvertToNumberingScheme(cell+12)+strconv.Itoa(valueRow), value)
					}
					j++
				}
				//对部门进行分阶段，然后在excel中写入到已合并的表格中
				//j := 1
				//departLen := 1

				//分析出部门所占的Cell数目
				//departs := make([]uint64 , 12)
				//departArr := make([]uint64, 0)
				//lenArr := make([]int, 0)
				//
				//length := 1
				//last := departmentMap[profile]["01"]
				//for _, d := range monthOrder[1:] {
				//	if departmentMap[profile][d] != last { // 这个月跟上个月不同，表示调整了部门
				//		lenArr = append(lenArr, length) // 记录上个部门所占有的月数
				//		departArr = append(departArr, last)
				//		last = departmentMap[profile][d]
				//		length = 1
				//	} else {
				//		length++
				//	}
				//
				//	if d == "12" { //当是最后一个月的时候
				//		lenArr = append(lenArr, length)
				//		departArr = append(departArr, last)
				//	}
				//}
				//
				//from := 1
				//for i, d := range departArr {
				//	if d == 0 {
				//		from += lenArr[i] + 1
				//		continue
				//	}
				//	start := util.ConvertToNumberingScheme(fieldCount*13 + startCol + from)
				//	if lenArr[i] > 1 {
				//		end := util.ConvertToNumberingScheme(fieldCount*13 + startCol + from + lenArr[i] - 1)
				//		xlsx.MergeCell(sheet1, start+strconv.Itoa(row), end+strconv.Itoa(row))
				//	}
				//
				//	xlsx.SetCellValue(sheet1, start+strconv.Itoa(row), allGroups[departmentIndexMap[d]].Name)
				//	from = from + lenArr[i]
				//}

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
