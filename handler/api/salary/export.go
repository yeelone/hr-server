package salary

import (
	"fmt"
	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/errno"
	"hrgdrc/util"
	"sort"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var xlsx *excelize.File
var sheetName = "Sheet1"

//导出到excel
func Export(c *gin.Context) {

	var r ExportRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}
	accountID := r.AccountID
	year := r.Year
	month := r.Month
	//用户上传模板，将之保存于 export/template 中
	if len(r.Template) > 0 {
		templateFile := "export/template/工资导出模板.xlsx"
		if util.Exists(templateFile) {
			err := util.MoveFile(templateFile, "export/template/old/工资导出模板.xlsx")
			if err != nil {
				fmt.Println("cannot move file to new directory" + err.Error())
			}
		}
		if util.Exists(r.Template) {
			err := util.MoveFile(r.Template, templateFile)
			if err != nil {
				fmt.Println("cannot move file to new directory" + err.Error())
			}
		}
	}

	account, err := model.GetTemplateAccount(accountID)
	if err != nil {
		fmt.Println("无法根据账套找到模板", err)
	}

	fieldMap := make(map[uint64]map[string]interface{}) // profile id => name => value or content
	templateFileName := "export/template/工资导出模板.xlsx"

	isSheetExists := false

	if util.Exists(templateFileName) {
		xlsx, err = excelize.OpenFile(templateFileName)
		if xlsx.GetSheetIndex(account.Name) > 0 {
			isSheetExists = true
		}
		sheetName = account.Name
	}

	if !isSheetExists {
		xlsx = excelize.NewFile()
		xlsx.NewSheet(sheetName)
		sheetName = "Sheet1"
	}

	colsMap := make(map[string]struct{}) // 去除重复项
	colsOrder := []string{}
	for _, t := range account.Templates {
		fields := model.GetRelatedTemplateValue(year, month, t.Name, accountID, []string{})
		cols, _ := readTemplateStruct(t.Name)
		for _, col := range cols { //
			if _, ok := colsMap[col]; !ok {
				colsMap[col] = struct{}{}
				colsOrder = append(colsOrder, col)
			}
		}
		readFields(fields, fieldMap)
		writeExcel(fieldMap, cols, t.Name, year, month)
	}
	if isSheetExists {
		writeSummaryExcel(fieldMap, colsOrder)
	}
	// Save xlsx file by the given path.
	filename := account.Name + year + "-" + month + ".xlsx"
	err = xlsx.SaveAs("./export/" + filename)
	if err != nil {
		fmt.Println(err)
		h.SendResponse(c, errno.ErrCreateFile, err.Error())
		return
	}

	h.SendResponse(c, nil, ExportResponse{
		File: filename,
	})

}

func readTemplateStruct(name string) (cols []string, err error) {
	var runtimeViper = viper.New()
	keys := make(map[int]string)
	runtimeViper.AddConfigPath("conf/templates") // 如果没有指定配置文件，则解析默认的配置文件
	runtimeViper.SetConfigName(name)

	runtimeViper.SetConfigType("yaml")                  // 设置配置文件格式为YAML
	if err := runtimeViper.ReadInConfig(); err != nil { // viper解析配置文件
		return nil, err
	}

	//先把键跟顺序取出来
	for _, key := range runtimeViper.AllKeys() {
		s := strings.Split(key, ".")
		keys[runtimeViper.GetInt(s[0]+".order")] = s[0]
	}

	//根据key排序
	var orders []int
	for k := range keys {
		orders = append(orders, k)
	}
	sort.Ints(orders)
	cols = append(cols, "ID")
	cols = append(cols, "行别")
	for _, i := range orders {
		if runtimeViper.GetBool(keys[i] + ".visible") {
			name := util.Strip(runtimeViper.GetString(keys[i] + ".name"))
			cols = append(cols, name)
		}
	}

	return cols, nil
}

// 读取field内容或值
func readFields(fields []model.SalaryField, fieldMap map[uint64]map[string]interface{}) {
	// todo : 这是特别定制的,转化前要删掉
	groups, _ := model.GetGroupWithAllChildren("部门")
	codeMap := make(map[uint64]int)
	for _, g := range groups {
		codeMap[g.ID] = g.Code
	}

	for _, field := range fields {
		if _, ok := fieldMap[field.ProfileID]; !ok {
			fieldMap[field.ProfileID] = make(map[string]interface{})
		}
		if field.Value >= 0 {
			fieldMap[field.ProfileID][field.Name] = field.Value
		} else {
			fieldMap[field.ProfileID][field.Name] = field.Content
		}
		fieldMap[field.ProfileID]["行别"] = codeMap[field.DepartmentGroupID]
		fieldMap[field.ProfileID]["ID"] = field.ProfileID
	}
}

func writeExcel(fieldMap map[uint64]map[string]interface{}, cols []string, name string, year string, month string) {
	// Create a new sheet.
	index := xlsx.NewSheet(name)
	// Set value of a cell.
	if len(cols) < 1 {
		return
	}
	xlsx.MergeCell(name, "A1", util.ConvertToNumberingScheme(len(cols))+"1")
	xlsx.SetCellValue(name, "A1", year+"-"+month+" "+name)
	col := 1
	for _, colName := range cols {
		pos := util.ConvertToNumberingScheme(col) + strconv.Itoa(2)
		xlsx.SetCellValue(name, pos, colName)
		col++
	}

	row := 3
	for _, item := range fieldMap {
		col := 1
		for _, colName := range cols {
			pos := util.ConvertToNumberingScheme(col) + strconv.Itoa(row)
			xlsx.SetCellValue(name, pos, item[colName])
			col++
		}

		row++
	}
	// Set active sheet of the workbook.
	xlsx.SetActiveSheet(index)
}

func writeSummaryExcel(fieldMap map[uint64]map[string]interface{}, cols []string) {
	// 注意，用户上传的表中，项目名为  field.Alias ,而不是field.Name
	sheet := sheetName

	nameIndex := make(map[string]int)
	rows := xlsx.GetRows(sheet)
	for i, colCell := range rows[0] {
		nameIndex[util.Strip(colCell)] = i + 1
	}
	// Set value of a cell.
	if len(cols) < 1 {
		return
	}
	row := 2
	for _, item := range fieldMap {
		col := 1
		for _, colName := range cols {
			pos := util.ConvertToNumberingScheme(nameIndex[colName]) + strconv.Itoa(row)
			if nameIndex[colName] == 0 {
				continue
			}
			if _, ok := item[colName]; ok {
				xlsx.SetCellValue(sheet, pos, item[colName])
			} else {
				fmt.Println("name", colName)
			}
			col++
		}
		row++
	}
	// Set active sheet of the workbook.
	xlsx.SetActiveSheet(xlsx.GetSheetIndex(sheet))
}
