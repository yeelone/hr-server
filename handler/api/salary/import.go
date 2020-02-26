package salary

import (
	"fmt"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

func Import(c *gin.Context) {
	log.Info("Import Salary from excel to create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	file, err := c.FormFile("file")
	if err != nil {
		h.SendResponse(c, errno.ErrTemplateInvalid, nil)
		return
	}

	filename, subffix := util.ExtractFileName(file.Filename)

	newfilename := "upload/salary/" + filename + "-" + time.Now().Format("20060102150405") + subffix

	if !util.Exists("upload/salary/") {
		os.MkdirAll("upload/salary/", os.ModePerm) //创建文件
	}

	if err := c.SaveUploadedFile(file, newfilename); err != nil {
		h.SendResponse(c, errno.ErrTemplateInvalid, nil)
		return
	}
	//上传成功之后，要对数据进行处理，比如在上传的表中，[sheet1]模板记录张三被扣罚300，[sheet2]模板记录张三其它收入1000，将这两条分散的数据归集到一起，再返回给客户端，
	//供操作员确认上传数据是否正确.
	//salaryDataMap := readDataFromExcel(newfilename)
	rsp := UploadResponse{}
	rsp.UploadFile = newfilename
	//rsp.DataPreview = salaryDataMap
	// Show the user information.
	h.SendResponse(c, nil, rsp)
}

func handleUploadedExcel(filename string) (fields string, err error) {
	//分析第一行
	xlsx, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println("OpenFile", err)
		return "", err
	}

	rows, _ := xlsx.GetRows("Sheet1")
	cols := make([]string, len(rows[0]))
	for index, colCell := range rows[0] {
		cols[index] = model.ProfileI18nMap[colCell]
	}

	return strings.Join(cols, ","), err
}

// readDataFromExcel : 读取上传的excel文件内容，保存到 map[string]map[string]interface{}
//  map[string]interface{} 结构如下:
//

func readDataFromExcel(filename string) (dataRows []map[string]interface{}) {
	xlsx, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
		return dataRows
	}
	dataRows = make([]map[string]interface{}, 0)
	dataMap := make(map[string]map[string]interface{}) // profile , sheet , column , value
	for _, sheet := range xlsx.GetSheetMap() {
		rows, _ := xlsx.GetRows(sheet)
		for _, row := range rows[2:] {
			idCard := row[0]
			if _, ok := dataMap[idCard]; !ok {
				dataMap[idCard] = make(map[string]interface{})
			}
			dataMap[idCard]["__name__"] = row[1]

			for j, col := range row[2:] {
				colName := rows[0][j+2]
				if len(col) > 1 {
					v, _ := strconv.ParseFloat(col, 64)
					dataMap[idCard][colName] = util.Decimal(v)
				}
			}
		}
	}

	for idcard, d := range dataMap {
		data := make(map[string]interface{})
		for k, v := range d {
			data[k] = v
		}
		data["__id_card__"] = idcard
		dataRows = append(dataRows, data)
	}
	return dataRows
}
